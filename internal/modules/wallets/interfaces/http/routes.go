package http

import (
	"net/http"

	"fin-flow-api/internal/shared/interface/jwt"
	"fin-flow-api/internal/shared/middleware"
)

var walletHandler *Handler

func SetupRoutes(mux *http.ServeMux, jwtService jwt.Service) {
	mountWallets(mux, jwtService)
}

func mountWallets(mux *http.ServeMux, jwtService jwt.Service) {
	mux.HandleFunc("/wallets/currencies", walletHandler.GetCurrencies)
	mux.HandleFunc("/wallets/types", walletHandler.GetWalletTypes)

	mux.HandleFunc("/wallets", handleWalletsCollection(jwtService))

	protectedHandler := middleware.RequireAuth(jwtService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/wallets/currencies" || path == "/wallets/types" {
			http.NotFound(w, r)
			return
		}
		handleWalletsResource(w, r)
	}))
	mux.Handle("/wallets/", protectedHandler)
}

func handleWalletsCollection(jwtService jwt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.RequireAuth(jwtService)(http.HandlerFunc(walletHandler.ListWallets)).ServeHTTP(w, r)
		case http.MethodPost:
			middleware.RequireAuth(jwtService)(http.HandlerFunc(walletHandler.CreateWallet)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleWalletsResource(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		walletHandler.GetWallet(w, r)
	case http.MethodPut:
		walletHandler.UpdateWallet(w, r)
	case http.MethodDelete:
		walletHandler.DeleteWallet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func SetHandler(handler *Handler) {
	walletHandler = handler
}