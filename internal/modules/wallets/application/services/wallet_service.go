package services

import (
	"context"
	"errors"
	"fin-flow-api/internal/modules/wallets/application/contracts/commands"
	"fin-flow-api/internal/modules/wallets/application/contracts/queries"
	"fin-flow-api/internal/modules/wallets/domain"
	"fin-flow-api/internal/shared/middleware"

	"github.com/google/uuid"
)

type WalletService struct {
	repository domain.WalletRepository
	systemUser string
}

func NewWalletService(repository domain.WalletRepository, systemUser string) *WalletService {
	return &WalletService{
		repository: repository,
		systemUser: systemUser,
	}
}

func (s *WalletService) getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	return userID, nil
}

func (s *WalletService) Create(ctx context.Context, req commands.WalletRequest) error {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	walletType := domain.WalletType(req.Type)
	if !domain.IsValidWalletType(req.Type) {
		return domain.ErrInvalidWalletType
	}

	if !domain.IsValidCurrency(req.Currency) {
		return domain.ErrInvalidCurrency
	}

	id := uuid.New().String()

	wallet := domain.NewWallet(
		id,
		userID,
		req.Name,
		walletType,
		req.Balance,
		domain.Currency(req.Currency),
		s.systemUser,
	)

	return s.repository.Create(wallet)
}

func (s *WalletService) Update(ctx context.Context, walletID string, req commands.WalletRequest) error {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	wallet, err := s.repository.GetByID(walletID, userID)
	if err != nil {
		return err
	}

	walletType := domain.WalletType(req.Type)
	if !domain.IsValidWalletType(req.Type) {
		return domain.ErrInvalidWalletType
	}

	if !domain.IsValidCurrency(req.Currency) {
		return domain.ErrInvalidCurrency
	}

	wallet.Name = req.Name
	wallet.Type = walletType
	wallet.Balance = req.Balance
	wallet.Currency = domain.Currency(req.Currency)
	wallet.Entity.UpdateModified(s.systemUser)

	return s.repository.Update(wallet)
}

func (s *WalletService) Delete(ctx context.Context, walletID string) error {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	return s.repository.Delete(walletID, userID)
}

func (s *WalletService) GetByID(ctx context.Context, walletID string) (*queries.WalletResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	wallet, err := s.repository.GetByID(walletID, userID)
	if err != nil {
		return nil, err
	}

	return &queries.WalletResponse{
		ID:        wallet.ID,
		Name:      wallet.Name,
		Type:      wallet.Type.Value(),
		TypeName:  wallet.Type.String(),
		Balance:   wallet.Balance,
		Currency:  wallet.Currency.String(),
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.ModifiedAt,
		CreatedBy: wallet.CreatedBy,
		UpdatedBy: wallet.ModifiedBy,
	}, nil
}

func (s *WalletService) List(ctx context.Context) ([]*queries.WalletResponse, error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	wallets, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*queries.WalletResponse, len(wallets))
	for i, wallet := range wallets {
		responses[i] = &queries.WalletResponse{
			ID:        wallet.ID,
			Name:      wallet.Name,
			Type:      wallet.Type.Value(),
			TypeName:  wallet.Type.String(),
			Balance:   wallet.Balance,
			Currency:  wallet.Currency.String(),
			CreatedAt: wallet.CreatedAt,
			UpdatedAt: wallet.ModifiedAt,
			CreatedBy: wallet.CreatedBy,
			UpdatedBy: wallet.ModifiedBy,
		}
	}

	return responses, nil
}