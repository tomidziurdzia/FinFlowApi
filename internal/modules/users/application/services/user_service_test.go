package services

import (
	"errors"
	"fin-flow-api/internal/modules/users/application/contracts/commands"
	"fin-flow-api/internal/modules/users/application/contracts/queries"
	"fin-flow-api/internal/modules/users/domain"
	"testing"
)

type mockRepository struct {
	users      map[string]*domain.User
	createErr  error
	getByIDErr error
	updateErr  error
	deleteErr  error
	listErr    error
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *mockRepository) Create(user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockRepository) GetByID(id string) (*domain.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockRepository) GetByEmail(email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockRepository) GetByAuthID(authID string) (*domain.User, error) {
	for _, user := range m.users {
		if user.AuthID == authID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockRepository) Update(user *domain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockRepository) Delete(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *mockRepository) List() ([]*domain.User, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

type mockHashService struct {
	hashFunc   func(password string) (string, error)
	verifyFunc func(password, hash string) bool
}

func newMockHashService() *mockHashService {
	return &mockHashService{
		hashFunc: func(password string) (string, error) {
			return "hashed-" + password, nil
		},
		verifyFunc: func(password, hash string) bool {
			return hash == "hashed-"+password
		},
	}
}

func (m *mockHashService) Hash(password string) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(password)
	}
	return "hashed", nil
}

func (m *mockHashService) Verify(password, hash string) bool {
	if m.verifyFunc != nil {
		return m.verifyFunc(password, hash)
	}
	return true
}

func TestNewUserService(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	systemUser := "system"

	service := NewUserService(repo, hashService, systemUser)
	if service == nil {
		t.Fatal("NewUserService returned nil")
	}
}

func TestUserService_Create(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	req := commands.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	err := service.Create(req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if len(repo.users) != 1 {
		t.Errorf("expected 1 user, got %d", len(repo.users))
	}

	var createdUser *domain.User
	for _, user := range repo.users {
		createdUser = user
		break
	}

	if createdUser.FirstName != req.FirstName {
		t.Errorf("expected FirstName %s, got %s", req.FirstName, createdUser.FirstName)
	}

	if createdUser.Email != req.Email {
		t.Errorf("expected Email %s, got %s", req.Email, createdUser.Email)
	}

	if createdUser.Password == req.Password {
		t.Error("password should be hashed, not stored as plain text")
	}
}

func TestUserService_Create_HashError(t *testing.T) {
	repo := newMockRepository()
	hashService := &mockHashService{
		hashFunc: func(password string) (string, error) {
			return "", errors.New("hash error")
		},
	}
	service := NewUserService(repo, hashService, "system")

	req := commands.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	err := service.Create(req)
	if err == nil {
		t.Error("Create should fail when hash fails")
	}
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	user := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed", "system")
	repo.users["user-1"] = user

	req := queries.GetUserRequest{ID: "user-1"}
	response, err := service.GetByID(req)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if response.ID != "user-1" {
		t.Errorf("expected ID user-1, got %s", response.ID)
	}

	if response.FirstName != "John" {
		t.Errorf("expected FirstName John, got %s", response.FirstName)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	req := queries.GetUserRequest{ID: "nonexistent"}
	_, err := service.GetByID(req)
	if err == nil {
		t.Error("GetByID should fail when user not found")
	}
}

func TestUserService_Update(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "admin")

	user := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed", "system")
	repo.users["user-1"] = user

	req := commands.UpdateUserRequest{
		ID:        "user-1",
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	err := service.Update(req)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updatedUser := repo.users["user-1"]
	if updatedUser.FirstName != "Jane" {
		t.Errorf("expected FirstName Jane, got %s", updatedUser.FirstName)
	}

	if updatedUser.ModifiedBy != "admin" {
		t.Errorf("expected ModifiedBy admin, got %s", updatedUser.ModifiedBy)
	}

	if updatedUser.ModifiedAt.Before(user.ModifiedAt) {
		t.Error("ModifiedAt should be updated")
	}
}

func TestUserService_Update_NotFound(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	req := commands.UpdateUserRequest{
		ID:        "nonexistent",
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	err := service.Update(req)
	if err == nil {
		t.Error("Update should fail when user not found")
	}
}

func TestUserService_Delete(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	user := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed", "system")
	repo.users["user-1"] = user

	req := commands.DeleteUserRequest{ID: "user-1"}
	err := service.Delete(req)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, exists := repo.users["user-1"]; exists {
		t.Error("user should be deleted")
	}
}

func TestUserService_Delete_NotFound(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	req := commands.DeleteUserRequest{ID: "nonexistent"}
	err := service.Delete(req)
	if err == nil {
		t.Error("Delete should fail when user not found")
	}
}

func TestUserService_List(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	user1 := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed", "system")
	user2 := domain.NewUser("user-2", "Jane", "Smith", "jane@example.com", "hashed", "system")
	repo.users["user-1"] = user1
	repo.users["user-2"] = user2

	req := queries.ListUsersRequest{}
	responses, err := service.List(req)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(responses) != 2 {
		t.Errorf("expected 2 users, got %d", len(responses))
	}
}

func TestUserService_List_Empty(t *testing.T) {
	repo := newMockRepository()
	hashService := newMockHashService()
	service := NewUserService(repo, hashService, "system")

	req := queries.ListUsersRequest{}
	responses, err := service.List(req)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(responses) != 0 {
		t.Errorf("expected 0 users, got %d", len(responses))
	}
}