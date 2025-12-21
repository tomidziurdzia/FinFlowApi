package postgres

import (
	"os"
	"testing"

	"fin-flow-api/internal/infrastructure/config"
	"fin-flow-api/internal/infrastructure/db"
	"fin-flow-api/internal/modules/categories/domain"

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

	userID := uuid.New().String()
	category := domain.NewCategory(
		uuid.New().String(),
		userID,
		"Test Category",
		domain.CategoryTypeExpense,
		"test-user",
	)

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
}

func TestRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	categoryID := uuid.New().String()
	category := domain.NewCategory(
		categoryID,
		userID,
		"Test Category",
		domain.CategoryTypeExpense,
		"test-user",
	)

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.GetByID(categoryID, userID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.ID != categoryID {
		t.Errorf("expected ID %s, got %s", categoryID, found.ID)
	}

	if found.Name != "Test Category" {
		t.Errorf("expected Name Test Category, got %s", found.Name)
	}

	if found.Type != domain.CategoryTypeExpense {
		t.Errorf("expected Type CategoryTypeExpense, got %v", found.Type)
	}
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := repo.GetByID("nonexistent", "user-1")
	if err == nil {
		t.Error("GetByID should fail when category not found")
	}
}

func TestRepository_GetByID_WrongUser(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID1 := uuid.New().String()
	userID2 := uuid.New().String()
	categoryID := uuid.New().String()
	category := domain.NewCategory(
		categoryID,
		userID1,
		"Test Category",
		domain.CategoryTypeExpense,
		"test-user",
	)

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err = repo.GetByID(categoryID, userID2)
	if err == nil {
		t.Error("GetByID should fail when category belongs to different user")
	}
}

func TestRepository_List(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()

	category1 := domain.NewCategory(
		uuid.New().String(),
		userID,
		"Category 1",
		domain.CategoryTypeExpense,
		"test-user",
	)

	category2 := domain.NewCategory(
		uuid.New().String(),
		userID,
		"Category 2",
		domain.CategoryTypeIncome,
		"test-user",
	)

	err := repo.Create(category1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Create(category2)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	categories, err := repo.List(userID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(categories) < 2 {
		t.Errorf("expected at least 2 categories, got %d", len(categories))
	}
}

func TestRepository_List_Empty(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	categories, err := repo.List(userID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(categories) != 0 {
		t.Errorf("expected 0 categories, got %d", len(categories))
	}
}

func TestRepository_Update(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	categoryID := uuid.New().String()
	category := domain.NewCategory(
		categoryID,
		userID,
		"Original Name",
		domain.CategoryTypeExpense,
		"test-user",
	)

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	category.Name = "Updated Name"
	category.Type = domain.CategoryTypeIncome
	category.Entity.UpdateModified("test-user-updated")

	err = repo.Update(category)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := repo.GetByID(categoryID, userID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if updated.Name != "Updated Name" {
		t.Errorf("expected Name Updated Name, got %s", updated.Name)
	}

	if updated.Type != domain.CategoryTypeIncome {
		t.Errorf("expected Type CategoryTypeIncome, got %v", updated.Type)
	}
}

func TestRepository_Update_WrongUser(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID1 := uuid.New().String()
	userID2 := uuid.New().String()
	categoryID := uuid.New().String()
	category := domain.NewCategory(
		categoryID,
		userID1,
		"Original Name",
		domain.CategoryTypeExpense,
		"test-user",
	)

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	category.UserID = userID2
	category.Name = "Updated Name"
	category.Entity.UpdateModified("test-user-updated")

	err = repo.Update(category)
	if err == nil {
		t.Error("Update should fail when category belongs to different user")
	}
}

func TestRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	categoryID := uuid.New().String()
	category := domain.NewCategory(
		categoryID,
		userID,
		"Test Category",
		domain.CategoryTypeExpense,
		"test-user",
	)

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Delete(categoryID, userID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(categoryID, userID)
	if err == nil {
		t.Error("Category should be deleted")
	}
}

func TestRepository_Delete_WrongUser(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID1 := uuid.New().String()
	userID2 := uuid.New().String()
	categoryID := uuid.New().String()
	category := domain.NewCategory(
		categoryID,
		userID1,
		"Test Category",
		domain.CategoryTypeExpense,
		"test-user",
	)

	err := repo.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Delete(categoryID, userID2)
	if err == nil {
		t.Error("Delete should fail when category belongs to different user")
	}
}

func TestRepository_Delete_NotFound(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()
	err := repo.Delete("nonexistent", userID)
	if err == nil {
		t.Error("Delete should fail when category not found")
	}
}

func TestRepository_CategoryTypes(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	userID := uuid.New().String()

	tests := []struct {
		name         string
		categoryType domain.CategoryType
	}{
		{"Expense", domain.CategoryTypeExpense},
		{"Income", domain.CategoryTypeIncome},
		{"Investment", domain.CategoryTypeInvestment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := domain.NewCategory(
				uuid.New().String(),
				userID,
				"Test "+tt.name,
				tt.categoryType,
				"test-user",
			)

			err := repo.Create(category)
			if err != nil {
				t.Fatalf("Create failed: %v", err)
			}

			found, err := repo.GetByID(category.ID, userID)
			if err != nil {
				t.Fatalf("GetByID failed: %v", err)
			}

			if found.Type != tt.categoryType {
				t.Errorf("expected Type %v, got %v", tt.categoryType, found.Type)
			}
		})
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		os.Exit(0)
	}
	
	code := m.Run()
	os.Exit(code)
}