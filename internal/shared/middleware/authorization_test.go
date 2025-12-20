package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockJWTService struct {
	generateTokenFunc func(userID string) (string, error)
	validateTokenFunc func(tokenString string) (string, error)
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
	if tokenString == "valid-token" {
		return "user-123", nil
	}
	return "", http.ErrAbortHandler
}

func TestRequireAuth_ValidToken(t *testing.T) {
	jwtService := &mockJWTService{
		validateTokenFunc: func(tokenString string) (string, error) {
			if tokenString == "valid-token" {
				return "user-123", nil
			}
			return "", http.ErrAbortHandler
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r)
		if !ok {
			t.Error("UserID should be in context")
		}
		if userID != "user-123" {
			t.Errorf("expected userID 'user-123', got %s", userID)
		}
		w.WriteHeader(http.StatusOK)
	})

	authHandler := RequireAuth(jwtService)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestRequireAuth_NoAuthorizationHeader(t *testing.T) {
	jwtService := &mockJWTService{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when authorization header is missing")
	})

	authHandler := RequireAuth(jwtService)(handler)

	req := httptest.NewRequest("GET", "/test", nil)

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestRequireAuth_InvalidFormat(t *testing.T) {
	jwtService := &mockJWTService{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when authorization format is invalid")
	})

	authHandler := RequireAuth(jwtService)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	jwtService := &mockJWTService{
		validateTokenFunc: func(tokenString string) (string, error) {
			return "", http.ErrAbortHandler
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when token is invalid")
	})

	authHandler := RequireAuth(jwtService)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	jwtService := &mockJWTService{
		validateTokenFunc: func(tokenString string) (string, error) {
			return "user-456", nil
		},
	}

	var capturedUserID string
	var capturedOk bool

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID, capturedOk = GetUserIDFromContext(r)
		w.WriteHeader(http.StatusOK)
	})

	authHandler := RequireAuth(jwtService)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()
	authHandler.ServeHTTP(rr, req)

	if !capturedOk {
		t.Error("GetUserIDFromContext should return ok=true")
	}

	if capturedUserID != "user-456" {
		t.Errorf("expected userID 'user-456', got %s", capturedUserID)
	}
}

func TestGetUserIDFromContext_NotInContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	userID, ok := GetUserIDFromContext(req)
	if ok {
		t.Error("GetUserIDFromContext should return ok=false when userID is not in context")
	}

	if userID != "" {
		t.Errorf("expected empty userID, got %s", userID)
	}
}