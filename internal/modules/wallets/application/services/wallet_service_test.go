package services

import (
	"context"
	"errors"
	"fin-flow-api/internal/modules/wallets/application/contracts/commands"
	"fin-flow-api/internal/modules/wallets/domain"
	"fin-flow-api/internal/shared/middleware"
	"testing"
)

type mockWalletRepository struct {
	wallets  map[string]*domain.Wallet
	createErr  error
	getByIDErr error
	updateErr  error
	deleteErr  error
	listErr    error
}

func newMockWalletRepository() *mockWalletRepository {
	return &mockWalletRepository{
		wallets: make(map[string]*domain.Wallet),
	}
}

func (m *mockWalletRepository) Create(wallet *domain.Wallet) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.wallets[wallet.ID] = wallet
	return nil
}

func (m *mockWalletRepository) GetByID(id string, userID string) (*domain.Wallet, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	wallet, exists := m.wallets[id]
	if !exists {
		return nil, errors.New("wallet not found")
	}
	if wallet.UserID != userID {
		return nil, errors.New("unauthorized access to wallet")
	}
	return wallet, nil
}

func (m *mockWalletRepository) List(userID string) ([]*domain.Wallet, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var result []*domain.Wallet
	for _, wallet := range m.wallets {
		if wallet.UserID == userID {
			result = append(result, wallet)
		}
	}
	return result, nil
}

func (m *mockWalletRepository) Update(wallet *domain.Wallet) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	existing, exists := m.wallets[wallet.ID]
	if !exists {
		return errors.New("wallet not found")
	}
	if existing.UserID != wallet.UserID {
		return errors.New("unauthorized access to wallet")
	}
	m.wallets[wallet.ID] = wallet
	return nil
}

func (m *mockWalletRepository) Delete(id string, userID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	wallet, exists := m.wallets[id]
	if !exists {
		return errors.New("wallet not found")
	}
	if wallet.UserID != userID {
		return errors.New("unauthorized access to wallet")
	}
	delete(m.wallets, id)
	return nil
}

type mockContext struct {
	context.Context
	userID string
	hasID  bool
}

func (m *mockContext) Value(key interface{}) interface{} {
	if key == middleware.UserIDKey {
		if m.hasID {
			return m.userID
		}
		return nil
	}
	return nil
}

func TestNewWalletService(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	if service == nil {
		t.Fatal("NewWalletService returned nil")
	}
	if service.repository != repo {
		t.Error("repository not set correctly")
	}
	if service.systemUser != "system" {
		t.Error("systemUser not set correctly")
	}
}

func TestWalletService_Create(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.WalletRequest{
		Name:     "Main Account",
		Type:     0,
		Balance:  1000.50,
		Currency: "USD",
	}

	err := service.Create(ctx, req)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if len(repo.wallets) != 1 {
		t.Errorf("expected 1 wallet, got %d", len(repo.wallets))
	}
}

func TestWalletService_Create_InvalidType(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.WalletRequest{
		Name:     "Main Account",
		Type:     99,
		Balance:  1000.50,
		Currency: "USD",
	}

	err := service.Create(ctx, req)
	if err == nil {
		t.Error("Create should fail with invalid type")
	}
	if err != domain.ErrInvalidWalletType {
		t.Errorf("expected ErrInvalidWalletType, got %v", err)
	}
}

func TestWalletService_Create_InvalidCurrency(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.WalletRequest{
		Name:     "Main Account",
		Type:     0,
		Balance:  1000.50,
		Currency: "INVALID",
	}

	err := service.Create(ctx, req)
	if err == nil {
		t.Error("Create should fail with invalid currency")
	}
	if err != domain.ErrInvalidCurrency {
		t.Errorf("expected ErrInvalidCurrency, got %v", err)
	}
}

func TestWalletService_Create_NotAuthenticated(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{hasID: false}

	req := commands.WalletRequest{
		Name:     "Main Account",
		Type:     0,
		Balance:  1000.50,
		Currency: "USD",
	}

	err := service.Create(ctx, req)
	if err == nil {
		t.Error("Create should fail when user not authenticated")
	}
	if err.Error() != "user not authenticated" {
		t.Errorf("expected 'user not authenticated', got %v", err)
	}
}

func TestWalletService_GetByID(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	wallet := domain.NewWallet("wallet1", "user1", "Main Account", domain.WalletTypeBank, 1000.50, domain.CurrencyUSD, "system")
	repo.wallets["wallet1"] = wallet

	ctx := &mockContext{userID: "user1", hasID: true}

	result, err := service.GetByID(ctx, "wallet1")
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if result.ID != "wallet1" {
		t.Errorf("expected ID 'wallet1', got %s", result.ID)
	}
}

func TestWalletService_GetByID_NotFound(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	_, err := service.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Error("GetByID should fail when wallet not found")
	}
	if err.Error() != "wallet not found" {
		t.Errorf("expected 'wallet not found', got %v", err)
	}
}

func TestWalletService_GetByID_Unauthorized(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	wallet := domain.NewWallet("wallet1", "user1", "Main Account", domain.WalletTypeBank, 1000.50, domain.CurrencyUSD, "system")
	repo.wallets["wallet1"] = wallet

	ctx := &mockContext{userID: "user2", hasID: true}

	_, err := service.GetByID(ctx, "wallet1")
	if err == nil {
		t.Error("GetByID should fail when user doesn't own wallet")
	}
	if err.Error() != "unauthorized access to wallet" {
		t.Errorf("expected 'unauthorized access to wallet', got %v", err)
	}
}

func TestWalletService_Update(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	wallet := domain.NewWallet("wallet1", "user1", "Main Account", domain.WalletTypeBank, 1000.50, domain.CurrencyUSD, "system")
	repo.wallets["wallet1"] = wallet

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.WalletRequest{
		Name:     "Savings Account",
		Type:     4,
		Balance:  5000.00,
		Currency: "EUR",
	}

	err := service.Update(ctx, "wallet1", req)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	updated := repo.wallets["wallet1"]
	if updated.Name != "Savings Account" {
		t.Errorf("expected name 'Savings Account', got %s", updated.Name)
	}
	if updated.Balance != 5000.00 {
		t.Errorf("expected balance 5000.00, got %f", updated.Balance)
	}
}

func TestWalletService_Update_NotFound(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	req := commands.WalletRequest{
		Name:     "Savings Account",
		Type:     4,
		Balance:  5000.00,
		Currency: "EUR",
	}

	err := service.Update(ctx, "nonexistent", req)
	if err == nil {
		t.Error("Update should fail when wallet not found")
	}
	if err.Error() != "wallet not found" {
		t.Errorf("expected 'wallet not found', got %v", err)
	}
}

func TestWalletService_Update_Unauthorized(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	wallet := domain.NewWallet("wallet1", "user1", "Main Account", domain.WalletTypeBank, 1000.50, domain.CurrencyUSD, "system")
	repo.wallets["wallet1"] = wallet

	ctx := &mockContext{userID: "user2", hasID: true}

	req := commands.WalletRequest{
		Name:     "Savings Account",
		Type:     4,
		Balance:  5000.00,
		Currency: "EUR",
	}

	err := service.Update(ctx, "wallet1", req)
	if err == nil {
		t.Error("Update should fail when user doesn't own wallet")
	}
	if err.Error() != "unauthorized access to wallet" {
		t.Errorf("expected 'unauthorized access to wallet', got %v", err)
	}
}

func TestWalletService_Delete(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	wallet := domain.NewWallet("wallet1", "user1", "Main Account", domain.WalletTypeBank, 1000.50, domain.CurrencyUSD, "system")
	repo.wallets["wallet1"] = wallet

	ctx := &mockContext{userID: "user1", hasID: true}

	err := service.Delete(ctx, "wallet1")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	if _, exists := repo.wallets["wallet1"]; exists {
		t.Error("wallet should be deleted")
	}
}

func TestWalletService_Delete_NotFound(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	err := service.Delete(ctx, "nonexistent")
	if err == nil {
		t.Error("Delete should fail when wallet not found")
	}
	if err.Error() != "wallet not found" {
		t.Errorf("expected 'wallet not found', got %v", err)
	}
}

func TestWalletService_Delete_Unauthorized(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	wallet := domain.NewWallet("wallet1", "user1", "Main Account", domain.WalletTypeBank, 1000.50, domain.CurrencyUSD, "system")
	repo.wallets["wallet1"] = wallet

	ctx := &mockContext{userID: "user2", hasID: true}

	err := service.Delete(ctx, "wallet1")
	if err == nil {
		t.Error("Delete should fail when user doesn't own wallet")
	}
	if err.Error() != "unauthorized access to wallet" {
		t.Errorf("expected 'unauthorized access to wallet', got %v", err)
	}
}

func TestWalletService_List(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	wallet1 := domain.NewWallet("wallet1", "user1", "Main Account", domain.WalletTypeBank, 1000.50, domain.CurrencyUSD, "system")
	wallet2 := domain.NewWallet("wallet2", "user1", "Savings", domain.WalletTypeSavings, 5000.00, domain.CurrencyEUR, "system")
	wallet3 := domain.NewWallet("wallet3", "user2", "Cash", domain.WalletTypeCash, 100.00, domain.CurrencyUSD, "system")
	repo.wallets["wallet1"] = wallet1
	repo.wallets["wallet2"] = wallet2
	repo.wallets["wallet3"] = wallet3

	ctx := &mockContext{userID: "user1", hasID: true}

	result, err := service.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 wallets, got %d", len(result))
	}
}

func TestWalletService_List_Empty(t *testing.T) {
	repo := newMockWalletRepository()
	service := NewWalletService(repo, "system")

	ctx := &mockContext{userID: "user1", hasID: true}

	result, err := service.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 wallets, got %d", len(result))
	}
}