package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"fin-flow-api/internal/modules/wallets/application/contracts/commands"
	"fin-flow-api/internal/modules/wallets/application/contracts/queries"
	"fin-flow-api/internal/modules/wallets/domain"
	basehandler "fin-flow-api/internal/shared/http"
)

type walletService interface {
	Create(ctx context.Context, req commands.WalletRequest) error
	GetByID(ctx context.Context, id string) (*queries.WalletResponse, error)
	Update(ctx context.Context, id string, req commands.WalletRequest) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*queries.WalletResponse, error)
}

type Handler struct {
	walletService walletService
}

func NewHandler(walletService walletService) *Handler {
	return &Handler{
		walletService: walletService,
	}
}

func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqDTO WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid JSON format in request body")
		return
	}

	if err := validateWalletRequest(reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := commands.WalletRequest{
		Name:     reqDTO.Name,
		Type:     *reqDTO.Type,
		Balance:  *reqDTO.Balance,
		Currency: *reqDTO.Currency,
	}

	if err := h.walletService.Create(r.Context(), cmd); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "invalid wallet type") {
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid wallet type. Must be 0 (Bank), 1 (Cash), 2 (CreditCard), 3 (DebitCard), 4 (Savings), 5 (Investment), or 6 (Other)"
		} else if strings.Contains(errorMsg, "invalid currency") {
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid currency code"
		} else if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "duplicate") || strings.Contains(errorMsg, "unique") || strings.Contains(errorMsg, "name already exists") {
			statusCode = http.StatusConflict
			errorMsg = "A wallet with this name already exists"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "Wallet created successfully")
}

func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/wallets/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "Wallet ID is required in the URL path")
		return
	}

	wallet, err := h.walletService.GetByID(r.Context(), id)
	if err != nil {
		statusCode := http.StatusNotFound
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "unauthorized access") {
			statusCode = http.StatusForbidden
			errorMsg = "You do not have permission to access this wallet"
		} else if strings.Contains(errorMsg, "wallet not found") {
			errorMsg = "Wallet not found"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	response := WalletResponse{
		ID:        wallet.ID,
		Name:      wallet.Name,
		Type:      wallet.Type,
		TypeName:  wallet.TypeName,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
		CreatedBy: wallet.CreatedBy,
		UpdatedBy: wallet.UpdatedBy,
	}

	basehandler.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/wallets/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "Wallet ID is required in the URL path")
		return
	}

	var reqDTO WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, "Invalid JSON format in request body")
		return
	}

	if err := validateWalletRequest(reqDTO); err != nil {
		basehandler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := commands.WalletRequest{
		Name:     reqDTO.Name,
		Type:     *reqDTO.Type,
		Balance:  *reqDTO.Balance,
		Currency: *reqDTO.Currency,
	}

	if err := h.walletService.Update(r.Context(), id, cmd); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "invalid wallet type") {
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid wallet type. Must be 0 (Bank), 1 (Cash), 2 (CreditCard), 3 (DebitCard), 4 (Savings), 5 (Investment), or 6 (Other)"
		} else if strings.Contains(errorMsg, "invalid currency") {
			statusCode = http.StatusBadRequest
			errorMsg = "Invalid currency code"
		} else if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "unauthorized access") {
			statusCode = http.StatusForbidden
			errorMsg = "You do not have permission to update this wallet"
		} else if strings.Contains(errorMsg, "wallet not found") {
			statusCode = http.StatusNotFound
			errorMsg = "Wallet not found"
		} else if strings.Contains(errorMsg, "duplicate") || strings.Contains(errorMsg, "unique") || strings.Contains(errorMsg, "name already exists") {
			statusCode = http.StatusConflict
			errorMsg = "A wallet with this name already exists"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "Wallet updated successfully")
}

func (h *Handler) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/wallets/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		basehandler.WriteError(w, http.StatusBadRequest, "Wallet ID is required in the URL path")
		return
	}

	if err := h.walletService.Delete(r.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()
		
		if strings.Contains(errorMsg, "user not authenticated") {
			statusCode = http.StatusUnauthorized
			errorMsg = "Authentication required"
		} else if strings.Contains(errorMsg, "unauthorized access") {
			statusCode = http.StatusForbidden
			errorMsg = "You do not have permission to delete this wallet"
		} else if strings.Contains(errorMsg, "wallet not found") {
			statusCode = http.StatusNotFound
			errorMsg = "Wallet not found"
		}
		
		basehandler.WriteError(w, statusCode, errorMsg)
		return
	}

	basehandler.WriteSuccess(w, "Wallet deleted successfully")
}

func (h *Handler) ListWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	wallets, err := h.walletService.List(r.Context())
	if err != nil {
		basehandler.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]WalletResponse, len(wallets))
	for i, wallet := range wallets {
		responses[i] = WalletResponse{
			ID:        wallet.ID,
			Name:      wallet.Name,
			Type:      wallet.Type,
			TypeName:  wallet.TypeName,
			Balance:   wallet.Balance,
			Currency:  wallet.Currency,
			CreatedAt: wallet.CreatedAt,
			UpdatedAt: wallet.UpdatedAt,
			CreatedBy: wallet.CreatedBy,
			UpdatedBy: wallet.UpdatedBy,
		}
	}

	basehandler.WriteJSON(w, http.StatusOK, responses)
}

func (h *Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	currencies := domain.GetAllCurrencies()
	response := CurrenciesResponse{
		Fiat:   currencies["fiat"],
		Crypto: currencies["crypto"],
	}

	basehandler.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) GetWalletTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		basehandler.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	types := []WalletTypeResponse{
		{Value: domain.WalletTypeBank.Value(), Name: domain.WalletTypeBank.String()},
		{Value: domain.WalletTypeCash.Value(), Name: domain.WalletTypeCash.String()},
		{Value: domain.WalletTypeCreditCard.Value(), Name: domain.WalletTypeCreditCard.String()},
		{Value: domain.WalletTypeDebitCard.Value(), Name: domain.WalletTypeDebitCard.String()},
		{Value: domain.WalletTypeSavings.Value(), Name: domain.WalletTypeSavings.String()},
		{Value: domain.WalletTypeInvestment.Value(), Name: domain.WalletTypeInvestment.String()},
		{Value: domain.WalletTypeOther.Value(), Name: domain.WalletTypeOther.String()},
	}

	basehandler.WriteJSON(w, http.StatusOK, types)
}

func validateWalletRequest(req WalletRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return &ValidationError{Field: "name", Message: "Wallet name is required"}
	}
	if len(req.Name) < 2 {
		return &ValidationError{Field: "name", Message: "Wallet name must be at least 2 characters long"}
	}
	if len(req.Name) > 255 {
		return &ValidationError{Field: "name", Message: "Wallet name must not exceed 255 characters"}
	}

	if req.Type == nil {
		return &ValidationError{Field: "type", Message: "Wallet type is required"}
	}

	if !isValidWalletType(*req.Type) {
		return &ValidationError{Field: "type", Message: "Wallet type must be 0 (Bank), 1 (Cash), 2 (CreditCard), 3 (DebitCard), 4 (Savings), 5 (Investment), or 6 (Other)"}
	}

	if req.Balance == nil {
		return &ValidationError{Field: "balance", Message: "Balance is required"}
	}

	if req.Currency == nil {
		return &ValidationError{Field: "currency", Message: "Currency is required"}
	}

	if !domain.IsValidCurrency(*req.Currency) {
		return &ValidationError{Field: "currency", Message: "Invalid currency code"}
	}

	return nil
}

func isValidWalletType(typeValue int) bool {
	return typeValue >= 0 && typeValue <= 6
}


type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}