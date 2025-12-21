package services

import (
	"context"
	"errors"
	"fin-flow-api/internal/modules/categories/application/contracts/commands"
	"fin-flow-api/internal/modules/categories/application/contracts/queries"
	"fin-flow-api/internal/modules/categories/domain"
	"fin-flow-api/internal/shared/middleware"

	"github.com/google/uuid"
)

type CategoryService struct {
	repository domain.CategoryRepository
	systemUser string
}

func NewCategoryService(repository domain.CategoryRepository, systemUser string) *CategoryService {
	return &CategoryService{
		repository: repository,
		systemUser: systemUser,
	}
}

func (s *CategoryService) getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	return userID, nil
}

func (s *CategoryService) Create(ctx context.Context, req commands.CategoryRequest) error {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	categoryType := domain.CategoryType(req.Type)
	if !isValidCategoryType(categoryType) {
		return domain.ErrInvalidCategoryType
	}

	id := uuid.New().String()

	category := domain.NewCategory(
		id,
		userID,
		req.Name,
		categoryType,
		s.systemUser,
	)

	return s.repository.Create(category)
}

func (s *CategoryService) Update(ctx context.Context, categoryID string, req commands.CategoryRequest) error {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	category, err := s.repository.GetByID(categoryID, userID)
	if err != nil {
		return err
	}

	categoryType := domain.CategoryType(req.Type)
	if !isValidCategoryType(categoryType) {
		return domain.ErrInvalidCategoryType
	}

	category.Name = req.Name
	category.Type = categoryType
	category.Entity.UpdateModified(s.systemUser)

	return s.repository.Update(category)
}

func (s *CategoryService) Delete(ctx context.Context, categoryID string) error {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	return s.repository.Delete(categoryID, userID)
}

func (s *CategoryService) GetByID(ctx context.Context, categoryID string) (*queries.CategoryResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	category, err := s.repository.GetByID(categoryID, userID)
	if err != nil {
		return nil, err
	}

	return &queries.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Type:      category.Type.Value(),
		TypeName:  category.Type.String(),
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.ModifiedAt,
		CreatedBy: category.CreatedBy,
		UpdatedBy: category.ModifiedBy,
	}, nil
}

func (s *CategoryService) List(ctx context.Context) ([]*queries.CategoryResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	categories, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*queries.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = &queries.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Type:      category.Type.Value(),
			TypeName:  category.Type.String(),
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.ModifiedAt,
			CreatedBy: category.CreatedBy,
			UpdatedBy: category.ModifiedBy,
		}
	}

	return responses, nil
}

func isValidCategoryType(categoryType domain.CategoryType) bool {
	return categoryType == domain.CategoryTypeExpense ||
		categoryType == domain.CategoryTypeIncome ||
		categoryType == domain.CategoryTypeInvestment
}