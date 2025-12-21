package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fin-flow-api/internal/modules/categories/application/contracts/commands"
	"fin-flow-api/internal/modules/categories/application/contracts/queries"
	"fin-flow-api/internal/shared/middleware"
)

type mockCategoryService struct {
	createErr  error
	getByIDErr error
	updateErr  error
	deleteErr  error
	listErr    error
	category   *queries.CategoryResponse
	categories []*queries.CategoryResponse
}

func newMockCategoryService() *mockCategoryService {
	return &mockCategoryService{
		categories: []*queries.CategoryResponse{},
	}
}

func (m *mockCategoryService) Create(ctx context.Context, req commands.CategoryRequest) error {
	return m.createErr
}

func (m *mockCategoryService) GetByID(ctx context.Context, id string) (*queries.CategoryResponse, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.category, nil
}

func (m *mockCategoryService) Update(ctx context.Context, id string, req commands.CategoryRequest) error {
	return m.updateErr
}

func (m *mockCategoryService) Delete(ctx context.Context, id string) error {
	return m.deleteErr
}

func (m *mockCategoryService) List(ctx context.Context) ([]*queries.CategoryResponse, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.categories, nil
}

func createContextWithUserID(userID string) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, middleware.UserIDKey, userID)
}

func TestCreateCategory_Success(t *testing.T) {
	service := newMockCategoryService()
	handler := &Handler{categoryService: service}

	typeValue := 0
	body := CategoryRequest{
		Name: "Groceries",
		Type: &typeValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createContextWithUserID("user1"))

	rr := httptest.NewRecorder()
	handler.CreateCategory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	if response["message"] != "Category created successfully" {
		t.Errorf("expected success message, got %s", response["message"])
	}
}

func TestCreateCategory_InvalidMethod(t *testing.T) {
	service := newMockCategoryService()
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("GET", "/categories", nil)
	rr := httptest.NewRecorder()
	handler.CreateCategory(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestCreateCategory_InvalidBody(t *testing.T) {
	service := newMockCategoryService()
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("POST", "/categories", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()
	handler.CreateCategory(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateCategory_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    CategoryRequest
		wantErr bool
	}{
		{"empty name", CategoryRequest{Type: intPtr(0)}, true},
		{"short name", CategoryRequest{Name: "A", Type: intPtr(0)}, true},
		{"missing type", CategoryRequest{Name: "Groceries"}, true},
		{"invalid type", CategoryRequest{Name: "Groceries", Type: intPtr(99)}, true},
		{"valid request", CategoryRequest{Name: "Groceries", Type: intPtr(0)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCategoryRequest(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCategoryRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCategory_ServiceError(t *testing.T) {
	service := newMockCategoryService()
	service.createErr = errors.New("user not authenticated")
	handler := &Handler{categoryService: service}

	typeValue := 0
	body := CategoryRequest{
		Name: "Groceries",
		Type: &typeValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateCategory(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestGetCategory_Success(t *testing.T) {
	service := newMockCategoryService()
	service.category = &queries.CategoryResponse{
		ID:        "cat1",
		Name:      "Groceries",
		Type:      0,
		TypeName:  "Expense",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "system",
		UpdatedBy: "system",
	}
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("GET", "/categories/cat1", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.GetCategory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response CategoryResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != "cat1" {
		t.Errorf("expected ID 'cat1', got %s", response.ID)
	}
}

func TestGetCategory_NotFound(t *testing.T) {
	service := newMockCategoryService()
	service.getByIDErr = errors.New("category not found")
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("GET", "/categories/nonexistent", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.GetCategory(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestGetCategory_Forbidden(t *testing.T) {
	service := newMockCategoryService()
	service.getByIDErr = errors.New("unauthorized access to category")
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("GET", "/categories/cat1", nil)
	req = req.WithContext(createContextWithUserID("user2"))
	rr := httptest.NewRecorder()
	handler.GetCategory(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	if response["error"] != "You do not have permission to access this category" {
		t.Errorf("expected forbidden message, got %s", response["error"])
	}
}

func TestUpdateCategory_Success(t *testing.T) {
	service := newMockCategoryService()
	handler := &Handler{categoryService: service}

	typeValue := 1
	body := CategoryRequest{
		Name: "Food",
		Type: &typeValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/categories/cat1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.UpdateCategory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestUpdateCategory_Forbidden(t *testing.T) {
	service := newMockCategoryService()
	service.updateErr = errors.New("unauthorized access to category")
	handler := &Handler{categoryService: service}

	typeValue := 1
	body := CategoryRequest{
		Name: "Food",
		Type: &typeValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/categories/cat1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createContextWithUserID("user2"))
	rr := httptest.NewRecorder()
	handler.UpdateCategory(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	service := newMockCategoryService()
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("DELETE", "/categories/cat1", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.DeleteCategory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestDeleteCategory_Forbidden(t *testing.T) {
	service := newMockCategoryService()
	service.deleteErr = errors.New("unauthorized access to category")
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("DELETE", "/categories/cat1", nil)
	req = req.WithContext(createContextWithUserID("user2"))
	rr := httptest.NewRecorder()
	handler.DeleteCategory(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}
}

func TestListCategories_Success(t *testing.T) {
	service := newMockCategoryService()
	service.categories = []*queries.CategoryResponse{
		{
			ID:        "cat1",
			Name:      "Groceries",
			Type:      0,
			TypeName:  "Expense",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: "system",
			UpdatedBy: "system",
		},
	}
	handler := &Handler{categoryService: service}

	req := httptest.NewRequest("GET", "/categories", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.ListCategories(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response []CategoryResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if len(response) != 1 {
		t.Errorf("expected 1 category, got %d", len(response))
	}
}

func intPtr(i int) *int {
	return &i
}

