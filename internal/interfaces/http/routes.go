package http

import (
	"net/http"

	usershttp "fin-flow-api/internal/users/interfaces/http"
)

func SetupRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/health", HealthHandler)
	
	// User routes
	usershttp.SetupRoutes(mux)
}