package domain

import "errors"

type WalletType int

const (
	WalletTypeBank WalletType = iota
	WalletTypeCash
	WalletTypeCreditCard
	WalletTypeDebitCard
	WalletTypeSavings
	WalletTypeInvestment
	WalletTypeOther
)

var ErrInvalidWalletType = errors.New("invalid wallet type")

func (wt WalletType) String() string {
	switch wt {
	case WalletTypeBank:
		return "Bank"
	case WalletTypeCash:
		return "Cash"
	case WalletTypeCreditCard:
		return "CreditCard"
	case WalletTypeDebitCard:
		return "DebitCard"
	case WalletTypeSavings:
		return "Savings"
	case WalletTypeInvestment:
		return "Investment"
	case WalletTypeOther:
		return "Other"
	default:
		return "Unknown"
	}
}

func (wt WalletType) Value() int {
	return int(wt)
}

func IsValidWalletType(value int) bool {
	wt := WalletType(value)
	return wt >= WalletTypeBank && wt <= WalletTypeOther
}

func GetAllWalletTypes() []map[string]interface{} {
	return []map[string]interface{}{
		{"value": WalletTypeBank.Value(), "name": WalletTypeBank.String()},
		{"value": WalletTypeCash.Value(), "name": WalletTypeCash.String()},
		{"value": WalletTypeCreditCard.Value(), "name": WalletTypeCreditCard.String()},
		{"value": WalletTypeDebitCard.Value(), "name": WalletTypeDebitCard.String()},
		{"value": WalletTypeSavings.Value(), "name": WalletTypeSavings.String()},
		{"value": WalletTypeInvestment.Value(), "name": WalletTypeInvestment.String()},
		{"value": WalletTypeOther.Value(), "name": WalletTypeOther.String()},
	}
}