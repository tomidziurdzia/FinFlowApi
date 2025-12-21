package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fin-flow-api/internal/modules/users/domain"
)

func TestLogin_Success(t *testing.T) {
	userRepo := newMockUserRepository()
	user := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed-password123", "system")
	userRepo.users["user-1"] = user

	hashService := newMockHashService()
	hashService.verifyFunc = func(password, hash string) bool {
		return password == "password123" && hash == "hashed-password123"
	}

	jwtService := newMockJWTService()
	jwtService.generateTokenFunc = func(userID string) (string, error) {
		return "jwt-token-123", nil
	}

	handler := NewAuthHandler(userRepo, hashService, jwtService)

	body := LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response LoginResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Token == "" {
		t.Error("Token should not be empty")
	}
	if response.User.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", response.User.Email)
	}
}

func TestLogin_InvalidMethod(t *testing.T) {
	userRepo := newMockUserRepository()
	hashService := newMockHashService()
	jwtService := newMockJWTService()

	handler := NewAuthHandler(userRepo, hashService, jwtService)

	req := httptest.NewRequest("GET", "/auth/login", nil)
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestLogin_InvalidBody(t *testing.T) {
	userRepo := newMockUserRepository()
	hashService := newMockHashService()
	jwtService := newMockJWTService()

	handler := NewAuthHandler(userRepo, hashService, jwtService)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := newMockUserRepository()
	userRepo.getByEmailFunc = func(email string) (*domain.User, error) {
		return nil, errors.New("user not found")
	}

	hashService := newMockHashService()
	jwtService := newMockJWTService()

	handler := NewAuthHandler(userRepo, hashService, jwtService)

	body := LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := newMockUserRepository()
	user := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed-password123", "system")
	userRepo.users["user-1"] = user

	hashService := newMockHashService()
	hashService.verifyFunc = func(password, hash string) bool {
		return false
	}

	jwtService := newMockJWTService()

	handler := NewAuthHandler(userRepo, hashService, jwtService)

	body := LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestLogin_TokenGenerationError(t *testing.T) {
	userRepo := newMockUserRepository()
	user := domain.NewUser("user-1", "John", "Doe", "john@example.com", "hashed-password123", "system")
	userRepo.users["user-1"] = user

	hashService := newMockHashService()
	hashService.verifyFunc = func(password, hash string) bool {
		return true
	}

	jwtService := newMockJWTService()
	jwtService.generateTokenFunc = func(userID string) (string, error) {
		return "", errors.New("token generation failed")
	}

	handler := NewAuthHandler(userRepo, hashService, jwtService)

	body := LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
}