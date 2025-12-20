package http

import (
	"net/http"

	"fin-flow-api/internal/shared/interface/jwt"
	usershttp "fin-flow-api/internal/users/interfaces/http"
)

func SetupRoutes(mux *http.ServeMux, jwtService jwt.Service) {
	mount(mux, jwtService)
}

func mount(mux *http.ServeMux, jwtService jwt.Service) {
	mux.HandleFunc("/health", HealthHandler)
	
	mountV1(mux, jwtService)
}

func mountV1(mux *http.ServeMux, jwtService jwt.Service) {
	usershttp.SetupRoutes(mux, jwtService)
}