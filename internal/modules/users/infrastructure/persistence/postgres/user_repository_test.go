package postgres

import (
	"os"
	"testing"
	"time"

	"fin-flow-api/internal/infrastructure/config"
	"fin-flow-api/internal/infrastructure/db"
	"fin-flow-api/internal/modules/users/domain"

	"github.com/google/uuid"
)

func setupTestDB(t *testing.T) (*Repository, func()) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewDB(&cfg.Database)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	repo := NewRepository(database.Pool)

	cleanup := func() {
		database.Close()
	}

	return repo, cleanup
}

func TestRepository_Create(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	user := domain.NewUser(
		uuid.New().String(),
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
}

func TestRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	user := domain.NewUser(
		userID,
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.GetByID(userID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.ID != userID {
		t.Errorf("expected ID %s, got %s", userID, found.ID)
	}

	if found.FirstName != "John" {
		t.Errorf("expected FirstName John, got %s", found.FirstName)
	}

	if found.Email != "john@example.com" {
		t.Errorf("expected Email john@example.com, got %s", found.Email)
	}
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := repo.GetByID("nonexistent")
	if err == nil {
		t.Error("GetByID should fail when user not found")
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	email := "jane@example.com"
	user := domain.NewUser(
		uuid.New().String(),
		"Jane",
		"Smith",
		email,
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.GetByEmail(email)
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}

	if found.Email != email {
		t.Errorf("expected Email %s, got %s", email, found.Email)
	}
}

func TestRepository_GetByEmail_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := repo.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("GetByEmail should fail when user not found")
	}
}

func TestRepository_GetByAuthID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	authID := "auth-123"
	userID := uuid.New().String()
	user := domain.NewUserWithAuthID(
		userID,
		authID,
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.GetByAuthID(authID)
	if err != nil {
		t.Fatalf("GetByAuthID failed: %v", err)
	}

	if found.AuthID != authID {
		t.Errorf("expected AuthID %s, got %s", authID, found.AuthID)
	}

	if found.ID != userID {
		t.Errorf("expected ID %s, got %s", userID, found.ID)
	}
}

func TestRepository_GetByAuthID_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := repo.GetByAuthID("nonexistent-auth-id")
	if err == nil {
		t.Error("GetByAuthID should fail when user not found")
	}
}

func TestRepository_Update(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	user := domain.NewUser(
		userID,
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	originalModifiedAt := user.ModifiedAt
	time.Sleep(10 * time.Millisecond)

	user.FirstName = "Jane"
	user.LastName = "Smith"
	user.Email = "jane@example.com"
	user.Entity.UpdateModified("test-user-updated")

	err = repo.Update(user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.GetByID(userID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if updated.FirstName != "Jane" {
		t.Errorf("expected FirstName Jane, got %s", updated.FirstName)
	}

	if updated.LastName != "Smith" {
		t.Errorf("expected LastName Smith, got %s", updated.LastName)
	}

	if updated.Email != "jane@example.com" {
		t.Errorf("expected Email jane@example.com, got %s", updated.Email)
	}

	if updated.ModifiedBy != "test-user-updated" {
		t.Errorf("expected ModifiedBy test-user-updated, got %s", updated.ModifiedBy)
	}

	if !updated.ModifiedAt.After(originalModifiedAt) {
		t.Error("ModifiedAt should be updated")
	}
}

func TestRepository_Update_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	user := domain.NewUser(
		"nonexistent",
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Update(user)
	if err == nil {
		t.Error("Update should fail when user not found")
	}
}

func TestRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	user := domain.NewUser(
		userID,
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Delete(userID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(userID)
	if err == nil {
		t.Error("User should be deleted")
	}
}

func TestRepository_Delete_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	err := repo.Delete("nonexistent")
	if err == nil {
		t.Error("Delete should fail when user not found")
	}
}

func TestRepository_List(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	user1 := domain.NewUser(
		uuid.New().String(),
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	user2 := domain.NewUser(
		uuid.New().String(),
		"Jane",
		"Smith",
		"jane@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Create(user2)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	users, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(users) < 2 {
		t.Errorf("expected at least 2 users, got %d", len(users))
	}
}

func TestRepository_List_Empty(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	users, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestRepository_Create_WithAuthID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	authID := "auth-123"
	user := domain.NewUserWithAuthID(
		uuid.New().String(),
		authID,
		"John",
		"Doe",
		"john@example.com",
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.GetByAuthID(authID)
	if err != nil {
		t.Fatalf("GetByAuthID failed: %v", err)
	}

	if found.AuthID != authID {
		t.Errorf("expected AuthID %s, got %s", authID, found.AuthID)
	}
}

func TestRepository_UniqueEmail(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	email := "unique@example.com"
	user1 := domain.NewUser(
		uuid.New().String(),
		"John",
		"Doe",
		email,
		"hashedpassword",
		"test-user",
	)

	err := repo.Create(user1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	user2 := domain.NewUser(
		uuid.New().String(),
		"Jane",
		"Smith",
		email,
		"hashedpassword",
		"test-user",
	)

	err = repo.Create(user2)
	if err == nil {
		t.Error("Create should fail when email already exists")
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

