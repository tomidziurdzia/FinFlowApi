package http

import (
	"net/http"

	"fin-flow-api/internal/shared/interface/jwt"
	"fin-flow-api/internal/shared/middleware"
)

var userHandler *Handler
var authHandler *AuthHandler

func SetupRoutes(mux *http.ServeMux, jwtService jwt.Service) {
	mountUsers(mux, jwtService)
	mountAuth(mux)
}

func mountUsers(mux *http.ServeMux, jwtService jwt.Service) {
	mux.HandleFunc("/users", handleUsersCollection(jwtService))
	
	protectedHandler := middleware.RequireAuth(jwtService)(http.HandlerFunc(handleUsersResource))
	mux.Handle("/users/", protectedHandler)
	
	syncHandler := middleware.RequireClerkAuth(http.HandlerFunc(userHandler.SyncUser))
	mux.Handle("/users/sync", syncHandler)
}

func handleUsersCollection(jwtService jwt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.RequireAuth(jwtService)(http.HandlerFunc(userHandler.ListUsers)).ServeHTTP(w, r)
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleUsersResource(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userHandler.GetUser(w, r)
	case http.MethodPut:
		userHandler.UpdateUser(w, r)
	case http.MethodDelete:
		userHandler.DeleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func mountAuth(mux *http.ServeMux) {
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func SetHandler(handler *Handler) {
	userHandler = handler
}

func SetAuthHandler(handler *AuthHandler) {
	authHandler = handler
}