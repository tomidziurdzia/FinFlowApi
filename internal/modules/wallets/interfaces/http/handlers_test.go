package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fin-flow-api/internal/modules/wallets/application/contracts/commands"
	"fin-flow-api/internal/modules/wallets/application/contracts/queries"
	"fin-flow-api/internal/shared/middleware"
)

type mockWalletService struct {
	createErr  error
	getByIDErr error
	updateErr  error
	deleteErr  error
	listErr    error
	wallet     *queries.WalletResponse
	wallets    []*queries.WalletResponse
}

func newMockWalletService() *mockWalletService {
	return &mockWalletService{
		wallets: []*queries.WalletResponse{},
	}
}

func (m *mockWalletService) Create(ctx context.Context, req commands.WalletRequest) error {
	return m.createErr
}

func (m *mockWalletService) GetByID(ctx context.Context, id string) (*queries.WalletResponse, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.wallet, nil
}

func (m *mockWalletService) Update(ctx context.Context, id string, req commands.WalletRequest) error {
	return m.updateErr
}

func (m *mockWalletService) Delete(ctx context.Context, id string) error {
	return m.deleteErr
}

func (m *mockWalletService) List(ctx context.Context) ([]*queries.WalletResponse, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.wallets, nil
}

func createContextWithUserID(userID string) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, middleware.UserIDKey, userID)
}

func TestCreateWallet_Success(t *testing.T) {
	service := newMockWalletService()
	handler := &Handler{walletService: service}

	typeValue := 0
	balanceValue := 1000.50
	currencyValue := "USD"
	body := WalletRequest{
		Name:     "Main Account",
		Type:     &typeValue,
		Balance:  &balanceValue,
		Currency: &currencyValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createContextWithUserID("user1"))

	rr := httptest.NewRecorder()
	handler.CreateWallet(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	if response["message"] != "Wallet created successfully" {
		t.Errorf("expected success message, got %s", response["message"])
	}
}

func TestCreateWallet_InvalidMethod(t *testing.T) {
	service := newMockWalletService()
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("GET", "/wallets", nil)
	rr := httptest.NewRecorder()
	handler.CreateWallet(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestCreateWallet_InvalidBody(t *testing.T) {
	service := newMockWalletService()
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("POST", "/wallets", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()
	handler.CreateWallet(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestCreateWallet_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    WalletRequest
		wantErr bool
	}{
		{"empty name", WalletRequest{Type: intPtr(0), Balance: floatPtr(100.0), Currency: stringPtr("USD")}, true},
		{"short name", WalletRequest{Name: "A", Type: intPtr(0), Balance: floatPtr(100.0), Currency: stringPtr("USD")}, true},
		{"missing type", WalletRequest{Name: "Main Account", Balance: floatPtr(100.0), Currency: stringPtr("USD")}, true},
		{"invalid type", WalletRequest{Name: "Main Account", Type: intPtr(99), Balance: floatPtr(100.0), Currency: stringPtr("USD")}, true},
		{"missing balance", WalletRequest{Name: "Main Account", Type: intPtr(0), Currency: stringPtr("USD")}, true},
		{"missing currency", WalletRequest{Name: "Main Account", Type: intPtr(0), Balance: floatPtr(100.0)}, true},
		{"invalid currency", WalletRequest{Name: "Main Account", Type: intPtr(0), Balance: floatPtr(100.0), Currency: stringPtr("INVALID")}, true},
		{"valid request", WalletRequest{Name: "Main Account", Type: intPtr(0), Balance: floatPtr(1000.50), Currency: stringPtr("USD")}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWalletRequest(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateWalletRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateWallet_ServiceError(t *testing.T) {
	service := newMockWalletService()
	service.createErr = errors.New("user not authenticated")
	handler := &Handler{walletService: service}

	typeValue := 0
	balanceValue := 1000.50
	currencyValue := "USD"
	body := WalletRequest{
		Name:     "Main Account",
		Type:     &typeValue,
		Balance:  &balanceValue,
		Currency: &currencyValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallets", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.CreateWallet(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestGetWallet_Success(t *testing.T) {
	service := newMockWalletService()
	service.wallet = &queries.WalletResponse{
		ID:        "wallet1",
		Name:      "Main Account",
		Type:      0,
		TypeName:  "Bank",
		Balance:   1000.50,
		Currency:  "USD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "system",
		UpdatedBy: "system",
	}
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("GET", "/wallets/wallet1", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.GetWallet(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response WalletResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != "wallet1" {
		t.Errorf("expected ID 'wallet1', got %s", response.ID)
	}
}

func TestGetWallet_NotFound(t *testing.T) {
	service := newMockWalletService()
	service.getByIDErr = errors.New("wallet not found")
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("GET", "/wallets/nonexistent", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.GetWallet(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestGetWallet_Forbidden(t *testing.T) {
	service := newMockWalletService()
	service.getByIDErr = errors.New("unauthorized access to wallet")
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("GET", "/wallets/wallet1", nil)
	req = req.WithContext(createContextWithUserID("user2"))
	rr := httptest.NewRecorder()
	handler.GetWallet(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	if response["error"] != "You do not have permission to access this wallet" {
		t.Errorf("expected forbidden message, got %s", response["error"])
	}
}

func TestUpdateWallet_Success(t *testing.T) {
	service := newMockWalletService()
	handler := &Handler{walletService: service}

	typeValue := 4
	balanceValue := 5000.00
	currencyValue := "EUR"
	body := WalletRequest{
		Name:     "Savings Account",
		Type:     &typeValue,
		Balance:  &balanceValue,
		Currency: &currencyValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/wallets/wallet1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.UpdateWallet(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestUpdateWallet_Forbidden(t *testing.T) {
	service := newMockWalletService()
	service.updateErr = errors.New("unauthorized access to wallet")
	handler := &Handler{walletService: service}

	typeValue := 4
	balanceValue := 5000.00
	currencyValue := "EUR"
	body := WalletRequest{
		Name:     "Savings Account",
		Type:     &typeValue,
		Balance:  &balanceValue,
		Currency: &currencyValue,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/wallets/wallet1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createContextWithUserID("user2"))
	rr := httptest.NewRecorder()
	handler.UpdateWallet(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}
}

func TestDeleteWallet_Success(t *testing.T) {
	service := newMockWalletService()
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("DELETE", "/wallets/wallet1", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.DeleteWallet(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestDeleteWallet_Forbidden(t *testing.T) {
	service := newMockWalletService()
	service.deleteErr = errors.New("unauthorized access to wallet")
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("DELETE", "/wallets/wallet1", nil)
	req = req.WithContext(createContextWithUserID("user2"))
	rr := httptest.NewRecorder()
	handler.DeleteWallet(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rr.Code)
	}
}

func TestListWallets_Success(t *testing.T) {
	service := newMockWalletService()
	service.wallets = []*queries.WalletResponse{
		{
			ID:        "wallet1",
			Name:      "Main Account",
			Type:      0,
			TypeName:  "Bank",
			Balance:   1000.50,
			Currency:  "USD",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: "system",
			UpdatedBy: "system",
		},
	}
	handler := &Handler{walletService: service}

	req := httptest.NewRequest("GET", "/wallets", nil)
	req = req.WithContext(createContextWithUserID("user1"))
	rr := httptest.NewRecorder()
	handler.ListWallets(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var response []WalletResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if len(response) != 1 {
		t.Errorf("expected 1 wallet, got %d", len(response))
	}
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}