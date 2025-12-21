package domain

import (
	"testing"
	"time"
)

func TestNewCategory(t *testing.T) {
	id := "category-123"
	userID := "user-456"
	name := "Groceries"
	categoryType := CategoryTypeExpense
	createdBy := "system"

	category := NewCategory(id, userID, name, categoryType, createdBy)

	if category.ID != id {
		t.Errorf("Expected ID %s, got %s", id, category.ID)
	}

	if category.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, category.UserID)
	}

	if category.Name != name {
		t.Errorf("Expected Name %s, got %s", name, category.Name)
	}

	if category.Type != categoryType {
		t.Errorf("Expected Type %v, got %v", categoryType, category.Type)
	}

	if category.CreatedBy != createdBy {
		t.Errorf("Expected CreatedBy %s, got %s", createdBy, category.CreatedBy)
	}

	if category.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set, got zero time")
	}

	if category.ModifiedAt.IsZero() {
		t.Error("Expected ModifiedAt to be set, got zero time")
	}

	if category.CreatedAt.After(time.Now()) {
		t.Error("Expected CreatedAt to be in the past")
	}
}

func TestCategoryType_String(t *testing.T) {
	tests := []struct {
		categoryType CategoryType
		expected     string
	}{
		{CategoryTypeExpense, "Expense"},
		{CategoryTypeIncome, "Income"},
		{CategoryTypeInvestment, "Investment"},
		{CategoryType(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.categoryType.String(); got != tt.expected {
				t.Errorf("CategoryType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCategoryType_Value(t *testing.T) {
	tests := []struct {
		categoryType CategoryType
		expected     int
	}{
		{CategoryTypeExpense, 0},
		{CategoryTypeIncome, 1},
		{CategoryTypeInvestment, 2},
	}

	for _, tt := range tests {
		t.Run(tt.categoryType.String(), func(t *testing.T) {
			if got := tt.categoryType.Value(); got != tt.expected {
				t.Errorf("CategoryType.Value() = %v, want %v", got, tt.expected)
			}
		})
	}
}

