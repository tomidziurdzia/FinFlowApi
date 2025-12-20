package services

import (
	"fin-flow-api/internal/users/application/contracts/commands"
	"fin-flow-api/internal/users/application/contracts/queries"
	"fin-flow-api/internal/users/domain"

	"github.com/google/uuid"
)

// UserService handles all user-related operations
type UserService struct {
	repository domain.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repository domain.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

// Create creates a new user
func (s *UserService) Create(req commands.CreateUserRequest) error {
	// Generate ID
	id := uuid.New().String()

	// Create user using domain constructor
	user := domain.NewUser(
		id,
		req.AuthID,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		"system", // createdBy - deberías obtenerlo del contexto
	)

	// Save user
	return s.repository.Create(user)
}

// Update updates an existing user
func (s *UserService) Update(req commands.UpdateUserRequest) error {
	// Get existing user
	user, err := s.repository.GetByID(req.UserID)
	if err != nil {
		return err
	}

	// Update user fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.Entity.UpdateModified("system") // modifiedBy - deberías obtenerlo del contexto

	// Save changes
	return s.repository.Update(user)
}

// Delete deletes a user by ID
func (s *UserService) Delete(req commands.DeleteUserRequest) error {
	return s.repository.Delete(req.UserID)
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(req queries.GetUserRequest) (*queries.UserResponse, error) {
	user, err := s.repository.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}

	return &queries.UserResponse{
		ID:        user.ID,
		AuthID:    user.AuthID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.ModifiedAt,
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.ModifiedBy,
	}, nil
}

// List retrieves all users
func (s *UserService) List(req queries.ListUsersRequest) ([]*queries.UserResponse, error) {
	users, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	responses := make([]*queries.UserResponse, len(users))
	for i, user := range users {
		responses[i] = &queries.UserResponse{
			ID:        user.ID,
			AuthID:    user.AuthID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.ModifiedAt,
			CreatedBy: user.CreatedBy,
			UpdatedBy: user.ModifiedBy,
		}
	}

	return responses, nil
}

