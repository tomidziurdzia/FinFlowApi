package http

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"fin-flow-api/internal/modules/users/application/contracts/commands"
	userservices "fin-flow-api/internal/modules/users/application/services"
	basehandler "fin-flow-api/internal/shared/http"
	"fin-flow-api/internal/shared/middleware"
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

	var reqDTO UserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid JSON format in request body")
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
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "email already exists") {
			statusCode = http.StatusConflict
			errorMsg = "An account with this email address already exists"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "User account created successfully")
}

func validateCreateUserRequest(req UserRequest) error {
	req.FirstName = strings.TrimSpace(req.FirstName)
	if req.FirstName == "" {
		return &ValidationError{Field: "first_name", Message: "First name is required"}
	}
	if len(req.FirstName) < 2 {
		return &ValidationError{Field: "first_name", Message: "First name must be at least 2 characters long"}
	}
	if len(req.FirstName) > 255 {
		return &ValidationError{Field: "first_name", Message: "First name must not exceed 255 characters"}
	}

	req.LastName = strings.TrimSpace(req.LastName)
	if req.LastName == "" {
		return &ValidationError{Field: "last_name", Message: "Last name is required"}
	}
	if len(req.LastName) < 2 {
		return &ValidationError{Field: "last_name", Message: "Last name must be at least 2 characters long"}
	}
	if len(req.LastName) > 255 {
		return &ValidationError{Field: "last_name", Message: "Last name must not exceed 255 characters"}
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		return &ValidationError{Field: "email", Message: "Email address is required"}
	}
	if len(req.Email) > 255 {
		return &ValidationError{Field: "email", Message: "Email address must not exceed 255 characters"}
	}
	if !isValidEmail(req.Email) {
		return &ValidationError{Field: "email", Message: "Invalid email address format"}
	}

	if req.Password == "" {
		return &ValidationError{Field: "password", Message: "Password is required"}
	}
	if len(req.Password) < 8 {
		return &ValidationError{Field: "password", Message: "Password must be at least 8 characters long"}
	}
	if len(req.Password) > 255 {
		return &ValidationError{Field: "password", Message: "Password must not exceed 255 characters"}
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

	authenticatedUserID, _ := middleware.GetUserIDFromContext(r.Context())

	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "User ID is required in the URL path")
		return
	}

	if authenticatedUserID != "" && authenticatedUserID != id {
		basehandler.WriteError(w, http.StatusForbidden, "You can only view your own profile")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		statusCode := http.StatusNotFound
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "user not found") {
			errorMsg = "User not found"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
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

	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		basehandler.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "User ID is required in the URL path")
		return
	}

	if authenticatedUserID != id {
		basehandler.WriteError(w, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var reqDTO UserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid JSON format in request body")
		return
	}

	cmd := commands.UpdateUserRequest{
		FirstName: reqDTO.FirstName,
		LastName:  reqDTO.LastName,
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
	}

	if err := h.userService.Update(id, cmd); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "email already exists") {
			statusCode = http.StatusConflict
			errorMsg = "An account with this email address already exists"
		} else if strings.Contains(errorMsg, "user not found") {
			statusCode = http.StatusNotFound
			errorMsg = "User not found"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "User profile updated successfully")
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authenticatedUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		basehandler.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "User ID is required in the URL path")
		return
	}

	if authenticatedUserID != id {
		basehandler.WriteError(w, http.StatusForbidden, "You can only delete your own account")
		return
	}

	if err := h.userService.Delete(id); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "user not found") {
			statusCode = http.StatusNotFound
			errorMsg = "User not found"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "User account deleted successfully")
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	users, err := h.userService.List()
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
