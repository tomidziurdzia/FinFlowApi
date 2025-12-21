package http

import (
	"net/http"

	basehandler "fin-flow-api/internal/shared/http"
	"fin-flow-api/internal/shared/middleware"
)

func (h *Handler) SyncUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authID, ok := middleware.GetClerkAuthIDFromContext(r)
	if !ok {
		basehandler.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	firstName, lastName, nameOk := middleware.GetClerkNameFromContext(r)
	if !nameOk {
		firstName = ""
		lastName = ""
	}

	email, emailOk := middleware.GetClerkEmailFromContext(r)
	if !emailOk {
		email = ""
	}

	user, err := h.userService.SyncByAuthID(authID, firstName, lastName, email)
	if err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	basehandler.WriteJSON(w, http.StatusOK, response)
}