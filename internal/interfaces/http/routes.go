package http

import (
	"net/http"

	usershttp "fin-flow-api/internal/users/interfaces/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", HealthHandler)
	
	usershttp.SetupRoutes(mux)
}