package http

import (
	"net/http"

	"fin-flow-api/internal/shared/interface/jwt"
	"fin-flow-api/internal/shared/middleware"
)

var categoryHandler *Handler

func SetupRoutes(mux *http.ServeMux, jwtService jwt.Service) {
	mountCategories(mux, jwtService)
}

func mountCategories(mux *http.ServeMux, jwtService jwt.Service) {
	mux.HandleFunc("/categories", handleCategoriesCollection(jwtService))

	protectedHandler := middleware.RequireAuth(jwtService)(http.HandlerFunc(handleCategoriesResource))
	mux.Handle("/categories/", protectedHandler)
}

func handleCategoriesCollection(jwtService jwt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.RequireAuth(jwtService)(http.HandlerFunc(categoryHandler.ListCategories)).ServeHTTP(w, r)
		case http.MethodPost:
			middleware.RequireAuth(jwtService)(http.HandlerFunc(categoryHandler.CreateCategory)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleCategoriesResource(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categoryHandler.GetCategory(w, r)
	case http.MethodPut:
		categoryHandler.UpdateCategory(w, r)
	case http.MethodDelete:
		categoryHandler.DeleteCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func SetHandler(handler *Handler) {
	categoryHandler = handler
}