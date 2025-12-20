package domain

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	id := "test-id"
	firstName := "John"
	lastName := "Doe"
	email := "john@example.com"
	password := "hashedpassword"
	createdBy := "system"

	user := NewUser(id, firstName, lastName, email, password, createdBy)

	if user == nil {
		t.Fatal("NewUser returned nil")
	}

	if user.ID != id {
		t.Errorf("expected ID %s, got %s", id, user.ID)
	}

	if user.FirstName != firstName {
		t.Errorf("expected FirstName %s, got %s", firstName, user.FirstName)
	}

	if user.LastName != lastName {
		t.Errorf("expected LastName %s, got %s", lastName, user.LastName)
	}

	if user.Email != email {
		t.Errorf("expected Email %s, got %s", email, user.Email)
	}

	if user.Password != password {
		t.Errorf("expected Password %s, got %s", password, user.Password)
	}

	if user.CreatedBy != createdBy {
		t.Errorf("expected CreatedBy %s, got %s", createdBy, user.CreatedBy)
	}

	if user.ModifiedBy != createdBy {
		t.Errorf("expected ModifiedBy %s, got %s", createdBy, user.ModifiedBy)
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if user.ModifiedAt.IsZero() {
		t.Error("ModifiedAt should not be zero")
	}

	if user.CreatedAt.After(time.Now()) {
		t.Error("CreatedAt should not be in the future")
	}
}