package http

import (
	"net/http"

	usershttp "fin-flow-api/internal/modules/users/interfaces/http"
	"fin-flow-api/internal/shared/interface/jwt"
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