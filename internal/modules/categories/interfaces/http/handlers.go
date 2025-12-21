package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"fin-flow-api/internal/modules/categories/application/contracts/commands"
	"fin-flow-api/internal/modules/categories/application/contracts/queries"
	basehandler "fin-flow-api/internal/shared/http"
)

type categoryService interface {
	Create(ctx context.Context, req commands.CategoryRequest) error
	GetByID(ctx context.Context, id string) (*queries.CategoryResponse, error)
	Update(ctx context.Context, id string, req commands.CategoryRequest) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*queries.CategoryResponse, error)
}

type Handler struct {
	categoryService categoryService
}

func NewHandler(categoryService categoryService) *Handler {
	return &Handler{
		categoryService: categoryService,
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqDTO CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid JSON format in request body")
		return
	}

	if err := validateCategoryRequest(reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := commands.CategoryRequest{
		Name: reqDTO.Name,
		Type: *reqDTO.Type,
	}

	if err := h.categoryService.Create(r.Context(), cmd); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "invalid category type") {
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid category type. Must be 0 (Expense), 1 (Income), or 2 (Investment)"
		} else if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "duplicate") || strings.Contains(errorMsg, "unique") {
			statusCode = http.StatusConflict
			errorMsg = "A category with this name already exists"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "Category created successfully")
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/categories/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "Category ID is required in the URL path")
		return
	}

	category, err := h.categoryService.GetByID(r.Context(), id)
	if err != nil {
		statusCode := http.StatusNotFound
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "unauthorized access") {
			statusCode = http.StatusForbidden
			errorMsg = "You do not have permission to access this category"
		} else if strings.Contains(errorMsg, "category not found") {
			errorMsg = "Category not found"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	response := CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Type:      category.Type,
		TypeName:  category.TypeName,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		CreatedBy: category.CreatedBy,
		UpdatedBy: category.UpdatedBy,
	}

	basehandler.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/categories/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "Category ID is required in the URL path")
		return
	}

	var reqDTO CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid JSON format in request body")
		return
	}

	if err := validateCategoryRequest(reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := commands.CategoryRequest{
		Name: reqDTO.Name,
		Type: *reqDTO.Type,
	}

	if err := h.categoryService.Update(r.Context(), id, cmd); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "invalid category type") {
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid category type. Must be 0 (Expense), 1 (Income), or 2 (Investment)"
		} else if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "unauthorized access") {
			statusCode = http.StatusForbidden
			errorMsg = "You do not have permission to update this category"
		} else if strings.Contains(errorMsg, "category not found") {
			statusCode = http.StatusNotFound
			errorMsg = "Category not found"
		} else if strings.Contains(errorMsg, "duplicate") || strings.Contains(errorMsg, "unique") {
			statusCode = http.StatusConflict
			errorMsg = "A category with this name already exists"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "Category updated successfully")
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/categories/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "Category ID is required in the URL path")
		return
	}

	if err := h.categoryService.Delete(r.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "unauthorized access") {
			statusCode = http.StatusForbidden
			errorMsg = "You do not have permission to delete this category"
		} else if strings.Contains(errorMsg, "category not found") {
			statusCode = http.StatusNotFound
			errorMsg = "Category not found"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "Category deleted successfully")
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	categories, err := h.categoryService.List(r.Context())
	if err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Type:      category.Type,
			TypeName:  category.TypeName,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
			CreatedBy: category.CreatedBy,
			UpdatedBy: category.UpdatedBy,
		}
	}

	basehandler.WriteJSON(w, http.StatusOK, responses)
}

func validateCategoryRequest(req CategoryRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return &ValidationError{Field: "name", Message: "Category name is required"}
	}
	if len(req.Name) < 2 {
		return &ValidationError{Field: "name", Message: "Category name must be at least 2 characters long"}
	}
	if len(req.Name) > 255 {
		return &ValidationError{Field: "name", Message: "Category name must not exceed 255 characters"}
	}

	if req.Type == nil {
		return &ValidationError{Field: "type", Message: "Category type is required"}
	}

	if !isValidCategoryType(*req.Type) {
		return &ValidationError{Field: "type", Message: "Category type must be 0 (Expense), 1 (Income), or 2 (Investment)"}
	}

	return nil
}

func isValidCategoryType(typeValue int) bool {
	return typeValue >= 0 && typeValue <= 2
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}