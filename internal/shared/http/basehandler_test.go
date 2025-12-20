package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	data := map[string]string{
		"message": "test",
	}

	rr := httptest.NewRecorder()
	WriteJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["message"] != "test" {
		t.Errorf("expected message 'test', got %s", response["message"])
	}
}

func TestWriteJSON_CustomStatus(t *testing.T) {
	data := map[string]string{"error": "not found"}

	rr := httptest.NewRecorder()
	WriteJSON(rr, http.StatusNotFound, data)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteError(rr, http.StatusBadRequest, "Invalid input")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["error"] != "Invalid input" {
		t.Errorf("expected error 'Invalid input', got %s", response["error"])
	}
}

func TestWriteSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	WriteSuccess(rr, "Operation successful")

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["message"] != "Operation successful" {
		t.Errorf("expected message 'Operation successful', got %s", response["message"])
	}
}

func TestWriteJSON_ComplexData(t *testing.T) {
	data := map[string]interface{}{
		"id":    "123",
		"name":  "John",
		"age":   30,
		"items": []string{"item1", "item2"},
	}

	rr := httptest.NewRecorder()
	WriteJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["id"] != "123" {
		t.Errorf("expected id '123', got %v", response["id"])
	}
}

