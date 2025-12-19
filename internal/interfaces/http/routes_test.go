package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	mux := http.NewServeMux()
	SetupRoutes(mux)

	// Test health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("health endpoint returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test non-existent endpoint
	req404, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr404 := httptest.NewRecorder()
	mux.ServeHTTP(rr404, req404)

	if status := rr404.Code; status != http.StatusNotFound {
		t.Errorf("non-existent endpoint returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

