package jwt

import (
	"os"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("NewService returned nil")
	}
}

func TestNewService_WithEnvSecret(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	testSecret := "test-secret-key-12345"
	os.Setenv("JWT_SECRET", testSecret)

	service := NewService()
	if service == nil {
		t.Fatal("NewService returned nil")
	}
}

func TestGenerateToken(t *testing.T) {
	service := NewService()
	userID := "test-user-id"

	token, err := service.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestGenerateToken_DifferentUsers(t *testing.T) {
	service := NewService()
	userID1 := "user-1"
	userID2 := "user-2"

	token1, err1 := service.GenerateToken(userID1)
	if err1 != nil {
		t.Fatalf("GenerateToken failed: %v", err1)
	}

	token2, err2 := service.GenerateToken(userID2)
	if err2 != nil {
		t.Fatalf("GenerateToken failed: %v", err2)
	}

	if token1 == token2 {
		t.Error("Different users should produce different tokens")
	}
}

func TestValidateToken(t *testing.T) {
	service := NewService()
	userID := "test-user-id"

	token, err := service.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	validatedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("expected userID %s, got %s", userID, validatedUserID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	service := NewService()
	invalidToken := "invalid.token.here"

	_, err := service.ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken should fail for invalid token")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	service := NewService()

	_, err := service.ValidateToken("")
	if err == nil {
		t.Error("ValidateToken should fail for empty token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	service1 := NewService()
	userID := "test-user-id"

	token, err := service1.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "different-secret")
	service2 := NewService()

	_, err = service2.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken should fail for token signed with different secret")
	}
}

func TestTokenExpiration(t *testing.T) {
	service := NewService()
	userID := "test-user-id"

	token, err := service.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	validatedUserID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("expected userID %s, got %s", userID, validatedUserID)
	}

	time.Sleep(100 * time.Millisecond)

	validatedUserID2, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed after short delay: %v", err)
	}

	if validatedUserID2 != userID {
		t.Errorf("Token should still be valid after short delay")
	}
}