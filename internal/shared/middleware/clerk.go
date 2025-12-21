package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	basehandler "fin-flow-api/internal/shared/http"
)

const ClerkAuthIDKey contextKey = "clerkAuthID"
const ClerkFirstNameKey contextKey = "clerkFirstName"
const ClerkLastNameKey contextKey = "clerkLastName"
const ClerkEmailKey contextKey = "clerkEmail"

func RequireClerkAuth(next http.Handler) http.Handler {
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

		authID, firstName, lastName, email, err := extractClerkClaims(tokenString)
		if err != nil {
			basehandler.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), ClerkAuthIDKey, authID)
		ctx = context.WithValue(ctx, ClerkFirstNameKey, firstName)
		ctx = context.WithValue(ctx, ClerkLastNameKey, lastName)
		ctx = context.WithValue(ctx, ClerkEmailKey, email)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func extractClerkClaims(tokenString string) (authID, firstName, lastName, email string, err error) {
	claims, err := parseJWTClaims(tokenString)
	if err != nil {
		return "", "", "", "", err
	}

	if sub, ok := claims["sub"].(string); ok {
		authID = sub
	}
	if fn, ok := claims["first_name"].(string); ok {
		firstName = fn
	} else if fn, ok := claims["given_name"].(string); ok {
		firstName = fn
	}
	if ln, ok := claims["last_name"].(string); ok {
		lastName = ln
	} else if ln, ok := claims["family_name"].(string); ok {
		lastName = ln
	}
	if em, ok := claims["email"].(string); ok {
		email = em
	}

	if authID == "" {
		return "", "", "", "", http.ErrAbortHandler
	}

	return authID, firstName, lastName, email, nil
}

func parseJWTClaims(tokenString string) (map[string]interface{}, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, http.ErrAbortHandler
	}

	payload := parts[1]
	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func GetClerkAuthIDFromContext(r *http.Request) (string, bool) {
	authID, ok := r.Context().Value(ClerkAuthIDKey).(string)
	return authID, ok
}

func GetClerkNameFromContext(r *http.Request) (firstName, lastName string, ok bool) {
	firstName, firstNameOk := r.Context().Value(ClerkFirstNameKey).(string)
	lastName, lastNameOk := r.Context().Value(ClerkLastNameKey).(string)
	return firstName, lastName, firstNameOk && lastNameOk
}

func GetClerkEmailFromContext(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(ClerkEmailKey).(string)
	return email, ok
}