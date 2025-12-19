package http

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	cfg := Config{
		Addr: "8080",
	}

	server := NewServer(cfg)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}

	if server.httpServer == nil {
		t.Fatal("httpServer is nil")
	}

	if server.httpServer.Addr != ":8080" {
		t.Errorf("expected addr :8080, got %s", server.httpServer.Addr)
	}
}

func TestNormalizeAddr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ":8080"},
		{"port only", "8080", ":8080"},
		{"with colon", ":8080", ":8080"},
		{"double colon", "::8080", ":8080"},
		{"different port", "3000", ":3000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeAddr(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeAddr(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestServerShutdown(t *testing.T) {
	cfg := Config{
		Addr: "0", // Use random port
	}

	server := NewServer(cfg)

	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.httpServer.ListenAndServe()
	}()

	// Give server time to start
	time.Sleep(10 * time.Millisecond)

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := server.httpServer.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Check that server stopped
	select {
	case err := <-errChan:
		if err != http.ErrServerClosed {
			t.Errorf("expected ErrServerClosed, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("server did not stop in time")
	}
}

