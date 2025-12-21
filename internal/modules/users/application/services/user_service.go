package services

import (
	"fin-flow-api/internal/shared/interface/hash"
	"fin-flow-api/internal/modules/users/application/contracts/commands"
	"fin-flow-api/internal/modules/users/application/contracts/queries"
	"fin-flow-api/internal/modules/users/domain"

	"github.com/google/uuid"
)

type UserService struct {
	repository domain.UserRepository
	hashService hash.Service
	systemUser  string
}

func NewUserService(repository domain.UserRepository, hashService hash.Service, systemUser string) *UserService {
	return &UserService{
		repository: repository,
		hashService: hashService,
		systemUser: systemUser,
	}
}

func (s *UserService) Create(req commands.CreateUserRequest) error {
	hashedPassword, err := s.hashService.Hash(req.Password)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	user := domain.NewUser(
		id,
		req.FirstName,
		req.LastName,
		req.Email,
		hashedPassword,
		s.systemUser,
	)

	return s.repository.Create(user)
}

func (s *UserService) Update(req commands.UpdateUserRequest) error {
	user, err := s.repository.GetByID(req.ID)
	if err != nil {
		return err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.Entity.UpdateModified(s.systemUser)

	return s.repository.Update(user)
}

func (s *UserService) Delete(req commands.DeleteUserRequest) error {
	return s.repository.Delete(req.ID)
}

func (s *UserService) GetByID(req queries.GetUserRequest) (*queries.UserResponse, error) {
	user, err := s.repository.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	return &queries.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.ModifiedAt,
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.ModifiedBy,
	}, nil
}

func (s *UserService) List(req queries.ListUsersRequest) ([]*queries.UserResponse, error) {
	users, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	responses := make([]*queries.UserResponse, len(users))
	for i, user := range users {
		responses[i] = &queries.UserResponse{
			ID:        user.ID,
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

func (s *UserService) SyncByAuthID(authID, firstName, lastName, email string) (*queries.UserResponse, error) {
	user, err := s.repository.GetByAuthID(authID)
	if err != nil {
		if err.Error() == "user not found" {
			id := uuid.New().String()
			newUser := domain.NewUserWithAuthID(
				id,
				authID,
				firstName,
				lastName,
				email,
				"",
				s.systemUser,
			)
			
			if createErr := s.repository.Create(newUser); createErr != nil {
				return nil, createErr
			}
			
			return &queries.UserResponse{
				ID:        newUser.ID,
				FirstName: newUser.FirstName,
				LastName:  newUser.LastName,
				Email:     newUser.Email,
				CreatedAt: newUser.CreatedAt,
				UpdatedAt: newUser.ModifiedAt,
				CreatedBy: newUser.CreatedBy,
				UpdatedBy: newUser.ModifiedBy,
			}, nil
		}
		return nil, err
	}

	return &queries.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.ModifiedAt,
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.ModifiedBy,
	}, nil
}