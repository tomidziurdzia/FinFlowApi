package http

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	basehandler "fin-flow-api/internal/shared/http"
	"fin-flow-api/internal/shared/middleware"
	"fin-flow-api/internal/users/application/contracts/commands"
	"fin-flow-api/internal/users/application/contracts/queries"
	userservices "fin-flow-api/internal/users/application/services"
)

type Handler struct {
	userService *userservices.UserService
}

func NewHandler(userService *userservices.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqDTO CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validateCreateUserRequest(reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := commands.CreateUserRequest{
		FirstName: reqDTO.FirstName,
		LastName:  reqDTO.LastName,
		Email:      reqDTO.Email,
		Password:   reqDTO.Password,
	}

	if err := h.userService.Create(cmd); err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	basehandler.WriteSuccess(w, "User created successfully")
}

func validateCreateUserRequest(req CreateUserRequest) error {
	req.FirstName = strings.TrimSpace(req.FirstName)
	if req.FirstName == "" {
		return &ValidationError{Field: "first_name", Message: "First name is required"}
	}
	if len(req.FirstName) < 2 {
		return &ValidationError{Field: "first_name", Message: "First name must be at least 2 characters"}
	}

	req.LastName = strings.TrimSpace(req.LastName)
	if req.LastName == "" {
		return &ValidationError{Field: "last_name", Message: "Last name is required"}
	}
	if len(req.LastName) < 2 {
		return &ValidationError{Field: "last_name", Message: "Last name must be at least 2 characters"}
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	if !isValidEmail(req.Email) {
		return &ValidationError{Field: "email", Message: "Invalid email format"}
	}

	if req.Password == "" {
		return &ValidationError{Field: "password", Message: "Password is required"}
	}
	if len(req.Password) < 8 {
		return &ValidationError{Field: "password", Message: "Password must be at least 8 characters"}
	}

	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authenticatedUserID, _ := middleware.GetUserIDFromContext(r)

	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if authenticatedUserID != "" && authenticatedUserID != id {
		basehandler.WriteError(w, http.StatusForbidden, "You can only view your own profile")
		return
	}

	query := queries.GetUserRequest{
		ID: id,
	}

	user, err := h.userService.GetByID(query)
	if err != nil {
		basehandler.WriteError(w, http.StatusNotFound, err.Error())
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

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authenticatedUserID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		basehandler.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if authenticatedUserID != id {
		basehandler.WriteError(w, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var reqDTO UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.UpdateUserRequest{
		ID:        id,
		FirstName: reqDTO.FirstName,
		LastName:  reqDTO.LastName,
		Email:     reqDTO.Email,
	}

	if err := h.userService.Update(cmd); err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

		basehandler.WriteSuccess(w, "User updated successfully")
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authenticatedUserID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		basehandler.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if authenticatedUserID != id {
		basehandler.WriteError(w, http.StatusForbidden, "You can only delete your own account")
		return
	}

	cmd := commands.DeleteUserRequest{
		ID: id,
	}

	if err := h.userService.Delete(cmd); err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

		basehandler.WriteSuccess(w, "User deleted successfully")
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := queries.ListUsersRequest{}

	users, err := h.userService.List(query)
	if err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
	}

		basehandler.WriteJSON(w, http.StatusOK, responses)
}
