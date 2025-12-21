package http

import (
	"net/http"

	categorieshttp "fin-flow-api/internal/modules/categories/interfaces/http"
	usershttp "fin-flow-api/internal/modules/users/interfaces/http"
	walletshttp "fin-flow-api/internal/modules/wallets/interfaces/http"
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
	categorieshttp.SetupRoutes(mux, jwtService)
	walletshttp.SetupRoutes(mux, jwtService)
}