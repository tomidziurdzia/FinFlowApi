package http

import (
	"errors"
	"fin-flow-api/internal/modules/users/domain"
)

type mockUserRepository struct {
	users         map[string]*domain.User
	createFunc    func(user *domain.User) error
	getByIDFunc   func(id string) (*domain.User, error)
	getByEmailFunc func(email string) (*domain.User, error)
	getByAuthIDFunc func(authID string) (*domain.User, error)
	updateFunc    func(user *domain.User) error
	deleteFunc    func(id string) error
	listFunc      func() ([]*domain.User, error)
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(user)
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(id string) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(email)
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) GetByAuthID(authID string) (*domain.User, error) {
	if m.getByAuthIDFunc != nil {
		return m.getByAuthIDFunc(authID)
	}
	for _, user := range m.users {
		if user.AuthID == authID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) Update(user *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(user)
	}
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) List() ([]*domain.User, error) {
	if m.listFunc != nil {
		return m.listFunc()
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
	return "hashed-" + password, nil
}

func (m *mockHashService) Verify(password, hash string) bool {
	if m.verifyFunc != nil {
		return m.verifyFunc(password, hash)
	}
	return hash == "hashed-"+password
}

type mockJWTService struct {
	generateTokenFunc func(userID string) (string, error)
	validateTokenFunc func(tokenString string) (string, error)
}

func newMockJWTService() *mockJWTService {
	return &mockJWTService{
		generateTokenFunc: func(userID string) (string, error) {
			return "mock-token", nil
		},
		validateTokenFunc: func(tokenString string) (string, error) {
			return "user-123", nil
		},
	}
}

func (m *mockJWTService) GenerateToken(userID string) (string, error) {
	if m.generateTokenFunc != nil {
		return m.generateTokenFunc(userID)
	}
	return "mock-token", nil
}

func (m *mockJWTService) ValidateToken(tokenString string) (string, error) {
	if m.validateTokenFunc != nil {
		return m.validateTokenFunc(tokenString)
	}
	return "user-123", nil
}