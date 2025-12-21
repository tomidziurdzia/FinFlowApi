package services

import (
	"context"
	"errors"
	"fin-flow-api/internal/modules/categories/application/contracts/commands"
	"fin-flow-api/internal/modules/categories/domain"
	"fin-flow-api/internal/shared/middleware"
	"testing"
)

type mockCategoryRepository struct {
	categories map[string]*domain.Category
	createErr  error
	getByIDErr error
	updateErr  error
	deleteErr  error
	listErr    error
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[string]*domain.Category),
	}
}

func (m *mockCategoryRepository) Create(category *domain.Category) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.categories[category.ID] = category
	return nil
}

func (m *mockCategoryRepository) GetByID(id string, userID string) (*domain.Category, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	category, exists := m.categories[id]
	if !exists {
		return nil, errors.New("category not found")
	}
	if category.UserID != userID {
		return nil, errors.New("unauthorized access to category")
	}
	return category, nil
}

func (m *mockCategoryRepository) List(userID string) ([]*domain.Category, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var result []*domain.Category
	for _, category := range m.categories {
		if category.UserID == userID {
			result = append(result, category)
		}
	}
	return result, nil
}

func (m *mockCategoryRepository) Update(category *domain.Category) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	existing, exists := m.categories[category.ID]
	if !exists {
		return errors.New("category not found")
	}
	if existing.UserID != category.UserID {
		return errors.New("unauthorized access to category")
	}
	m.categories[category.ID] = category
	return nil
}

func (m *mockCategoryRepository) Delete(id string, userID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	category, exists := m.categories[id]
	if !exists {
		return errors.New("category not found")
	}
	if category.UserID != userID {
		return errors.New("unauthorized access to category")
	}
	delete(m.categories, id)
	return nil
}

type mockContext struct {
	context.Context
	userID string
	hasID  bool
}

func (m *mockContext) Value(key interface{}) interface{} {
	if key == middleware.UserIDKey {
		if m.hasID {
			return m.userID
		}
		return nil
	}
	return nil
}

func TestNewCategoryService(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	if service == nil {
		t.Fatal("NewCategoryService returned nil")
	}
	if service.repository != repo {
		t.Error("repository not set correctly")
	}
	if service.systemUser != "system" {
		t.Error("systemUser not set correctly")
	}
}

func TestCategoryService_Create(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.CategoryRequest{
		Name: "Groceries",
		Type: 0,
	}

	err := service.Create(ctx, req)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if len(repo.categories) != 1 {
		t.Errorf("expected 1 category, got %d", len(repo.categories))
	}
}

func TestCategoryService_Create_InvalidType(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.CategoryRequest{
		Name: "Groceries",
		Type: 99,
	}

	err := service.Create(ctx, req)
	if err == nil {
		t.Error("Create should fail with invalid type")
	}
	if err != domain.ErrInvalidCategoryType {
		t.Errorf("expected ErrInvalidCategoryType, got %v", err)
	}
}

func TestCategoryService_Create_NotAuthenticated(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	ctx := &mockContext{hasID: false}

	req := commands.CategoryRequest{
		Name: "Groceries",
		Type: 0,
	}

	err := service.Create(ctx, req)
	if err == nil {
		t.Error("Create should fail when user not authenticated")
	}
	if err.Error() != "user not authenticated" {
		t.Errorf("expected 'user not authenticated', got %v", err)
	}
}

func TestCategoryService_GetByID(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	category := domain.NewCategory("cat1", "user1", "Groceries", domain.CategoryTypeExpense, "system")
	repo.categories["cat1"] = category

	ctx := &mockContext{userID: "user1", hasID: true}

	result, err := service.GetByID(ctx, "cat1")
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if result.ID != "cat1" {
		t.Errorf("expected ID 'cat1', got %s", result.ID)
	}
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	_, err := service.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Error("GetByID should fail when category not found")
	}
	if err.Error() != "category not found" {
		t.Errorf("expected 'category not found', got %v", err)
	}
}

func TestCategoryService_GetByID_Unauthorized(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	category := domain.NewCategory("cat1", "user1", "Groceries", domain.CategoryTypeExpense, "system")
	repo.categories["cat1"] = category

	ctx := &mockContext{userID: "user2", hasID: true}

	_, err := service.GetByID(ctx, "cat1")
	if err == nil {
		t.Error("GetByID should fail when user doesn't own category")
	}
	if err.Error() != "unauthorized access to category" {
		t.Errorf("expected 'unauthorized access to category', got %v", err)
	}
}

func TestCategoryService_Update(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	category := domain.NewCategory("cat1", "user1", "Groceries", domain.CategoryTypeExpense, "system")
	repo.categories["cat1"] = category

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.CategoryRequest{
		Name: "Food",
		Type: 0,
	}

	err := service.Update(ctx, "cat1", req)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	updated := repo.categories["cat1"]
	if updated.Name != "Food" {
		t.Errorf("expected name 'Food', got %s", updated.Name)
	}
}

func TestCategoryService_Update_NotFound(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.CategoryRequest{
		Name: "Food",
		Type: 0,
	}

	err := service.Update(ctx, "nonexistent", req)
	if err == nil {
		t.Error("Update should fail when category not found")
	}
	if err.Error() != "category not found" {
		t.Errorf("expected 'category not found', got %v", err)
	}
}

func TestCategoryService_Update_Unauthorized(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	category := domain.NewCategory("cat1", "user1", "Groceries", domain.CategoryTypeExpense, "system")
	repo.categories["cat1"] = category

	ctx := &mockContext{userID: "user2", hasID: true}

	req := commands.CategoryRequest{
		Name: "Food",
		Type: 0,
	}

	err := service.Update(ctx, "cat1", req)
	if err == nil {
		t.Error("Update should fail when user doesn't own category")
	}
	if err.Error() != "unauthorized access to category" {
		t.Errorf("expected 'unauthorized access to category', got %v", err)
	}
}

func TestCategoryService_Delete(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	category := domain.NewCategory("cat1", "user1", "Groceries", domain.CategoryTypeExpense, "system")
	repo.categories["cat1"] = category

	ctx := &mockContext{userID: "user1", hasID: true}

	err := service.Delete(ctx, "cat1")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	if _, exists := repo.categories["cat1"]; exists {
		t.Error("category should be deleted")
	}
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	err := service.Delete(ctx, "nonexistent")
	if err == nil {
		t.Error("Delete should fail when category not found")
	}
	if err.Error() != "category not found" {
		t.Errorf("expected 'category not found', got %v", err)
	}
}

func TestCategoryService_Delete_Unauthorized(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	category := domain.NewCategory("cat1", "user1", "Groceries", domain.CategoryTypeExpense, "system")
	repo.categories["cat1"] = category

	ctx := &mockContext{userID: "user2", hasID: true}

	err := service.Delete(ctx, "cat1")
	if err == nil {
		t.Error("Delete should fail when user doesn't own category")
	}
	if err.Error() != "unauthorized access to category" {
		t.Errorf("expected 'unauthorized access to category', got %v", err)
	}
}

func TestCategoryService_List(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	category1 := domain.NewCategory("cat1", "user1", "Groceries", domain.CategoryTypeExpense, "system")
	category2 := domain.NewCategory("cat2", "user1", "Salary", domain.CategoryTypeIncome, "system")
	category3 := domain.NewCategory("cat3", "user2", "Rent", domain.CategoryTypeExpense, "system")
	repo.categories["cat1"] = category1
	repo.categories["cat2"] = category2
	repo.categories["cat3"] = category3

	ctx := &mockContext{userID: "user1", hasID: true}

	result, err := service.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 categories, got %d", len(result))
	}
}

func TestCategoryService_List_Empty(t *testing.T) {
	repo := newMockCategoryRepository()
	service := NewCategoryService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	result, err := service.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 categories, got %d", len(result))
	}
}