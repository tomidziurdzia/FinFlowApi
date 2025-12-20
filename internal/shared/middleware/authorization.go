package middleware

import (
	"context"
	"net/http"
	"strings"

	basehandler "fin-flow-api/internal/shared/http"
	"fin-flow-api/internal/shared/interface/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(jwtService jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				basehandler.WriteError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				basehandler.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]

			userID, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				basehandler.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}