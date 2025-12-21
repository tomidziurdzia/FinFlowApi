package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fin-flow-api/internal/shared/middleware"
	"fin-flow-api/internal/modules/users/application/services"
	"fin-flow-api/internal/modules/users/domain"
)

func TestCreateUser_Success(t *testing.T) {
	repo := newMockUserRepository()
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	body := CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	if response["message"] != "User created successfully" {
		t.Errorf("expected success message, got %s", response["message"])
	}
}

func TestCreateUser_InvalidMethod(t *testing.T) {
	repo := newMockUserRepository()
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestCreateUser_InvalidBody(t *testing.T) {
	repo := newMockUserRepository()
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateUser_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    CreateUserRequest
		wantErr bool
	}{
		{"empty first name", CreateUserRequest{LastName: "Doe", Email: "john@example.com", Password: "password123"}, true},
		{"short first name", CreateUserRequest{FirstName: "J", LastName: "Doe", Email: "john@example.com", Password: "password123"}, true},
		{"empty last name", CreateUserRequest{FirstName: "John", Email: "john@example.com", Password: "password123"}, true},
		{"short last name", CreateUserRequest{FirstName: "John", LastName: "D", Email: "john@example.com", Password: "password123"}, true},
		{"invalid email", CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "invalid-email", Password: "password123"}, true},
		{"empty password", CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "john@example.com"}, true},
		{"short password", CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "short"}, true},
		{"valid request", CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "password123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateUserRequest(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateUserRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUser_ServiceError(t *testing.T) {
	repo := newMockUserRepository()
	repo.createFunc = func(user *domain.User) error {
		return errors.New("repository error")
	}
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	body := CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateUser(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
}

func TestGetUser_Success(t *testing.T) {
	repo := newMockUserRepository()
	repo.getByIDFunc = func(id string) (*domain.User, error) {
		return domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed", "system"), nil
	}
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("GET", "/users/user-1", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response UserResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != "user-1" {
		t.Errorf("expected ID user-1, got %s", response.ID)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	repo := newMockUserRepository()
	repo.getByIDFunc = func(id string) (*domain.User, error) {
		return nil, errors.New("user not found")
	}
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("GET", "/users/nonexistent", nil)
	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestGetUser_Forbidden(t *testing.T) {
	repo := newMockUserRepository()
	repo.getByIDFunc = func(id string) (*domain.User, error) {
		return domain.NewUser("user-2", "Jane", "Doe", "jane@example.com", "hashed", "system"), nil
	}
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("GET", "/users/user-2", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	user := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed", "system")
	repo := newMockUserRepository()
	repo.getByIDFunc = func(id string) (*domain.User, error) {
		return user, nil
	}
	repo.updateFunc = func(u *domain.User) error {
		return nil
	}
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	body := UpdateUserRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/users/user-1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	repo := newMockUserRepository()
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("PUT", "/users/user-1", nil)
	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestUpdateUser_Forbidden(t *testing.T) {
	repo := newMockUserRepository()
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("PUT", "/users/user-2", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	repo := newMockUserRepository()
	repo.deleteFunc = func(id string) error {
		return nil
	}
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("DELETE", "/users/user-1", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestDeleteUser_Unauthorized(t *testing.T) {
	repo := newMockUserRepository()
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("DELETE", "/users/user-1", nil)
	rr := httptest.NewRecorder()
	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestListUsers_Success(t *testing.T) {
	repo := newMockUserRepository()
	repo.listFunc = func() ([]*domain.User, error) {
		return []*domain.User{
			domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed", "system"),
			domain.NewUser("user-2", "Jane", "Smith", "jane@example.com", "hashed", "system"),
		}, nil
	}
	hashService := newMockHashService()
	userService := services.NewUserService(repo, hashService, "system")
	handler := NewHandler(userService)

	req := httptest.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	handler.ListUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var responses []UserResponse
	json.NewDecoder(rr.Body).Decode(&responses)
	if len(responses) != 2 {
		t.Errorf("expected 2 users, got %d", len(responses))
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.com", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"user@", false},
		{"user@example", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}
}