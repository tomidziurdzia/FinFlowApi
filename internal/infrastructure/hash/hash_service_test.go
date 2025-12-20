package hash

import (
	"testing"
)

func TestHash(t *testing.T) {
	service := NewService()
	password := "testpassword123"

	hashed, err := service.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if hashed == "" {
		t.Error("Hash should not be empty")
	}

	if hashed == password {
		t.Error("Hash should not be the same as the password")
	}
}

func TestHash_DifferentPasswords(t *testing.T) {
	service := NewService()
	password1 := "password1"
	password2 := "password2"

	hashed1, err1 := service.Hash(password1)
	if err1 != nil {
		t.Fatalf("Hash failed: %v", err1)
	}

	hashed2, err2 := service.Hash(password2)
	if err2 != nil {
		t.Fatalf("Hash failed: %v", err2)
	}

	if hashed1 == hashed2 {
		t.Error("Different passwords should produce different hashes")
	}
}

func TestHash_SamePasswordDifferentHashes(t *testing.T) {
	service := NewService()
	password := "samepassword"

	hashed1, err1 := service.Hash(password)
	if err1 != nil {
		t.Fatalf("Hash failed: %v", err1)
	}

	hashed2, err2 := service.Hash(password)
	if err2 != nil {
		t.Fatalf("Hash failed: %v", err2)
	}

	if hashed1 == hashed2 {
		t.Error("Same password should produce different hashes (due to salt)")
	}
}

func TestVerify(t *testing.T) {
	service := NewService()
	password := "testpassword123"

	hashed, err := service.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if !service.Verify(password, hashed) {
		t.Error("Verify should return true for correct password")
	}
}

func TestVerify_WrongPassword(t *testing.T) {
	service := NewService()
	password := "correctpassword"
	wrongPassword := "wrongpassword"

	hashed, err := service.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if service.Verify(wrongPassword, hashed) {
		t.Error("Verify should return false for wrong password")
	}
}

func TestVerify_EmptyPassword(t *testing.T) {
	service := NewService()
	password := "testpassword"

	hashed, err := service.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if service.Verify("", hashed) {
		t.Error("Verify should return false for empty password")
	}
}