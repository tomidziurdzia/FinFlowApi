package http

import (
	"encoding/json"
	"net/http"

	basehandler "fin-flow-api/internal/shared/http"
	"fin-flow-api/internal/shared/interface/hash"
	"fin-flow-api/internal/shared/interface/jwt"
	"fin-flow-api/internal/users/domain"
)

type AuthHandler struct {
	userRepo   domain.UserRepository
	hashService hash.Service
	jwtService jwt.Service
}

func NewAuthHandler(userRepo domain.UserRepository, hashService hash.Service, jwtService jwt.Service) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		hashService: hashService,
		jwtService: jwtService,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		basehandler.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !h.hashService.Verify(req.Password, user.Password) {
		basehandler.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := LoginResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	basehandler.WriteJSON(w, http.StatusOK, response)
}