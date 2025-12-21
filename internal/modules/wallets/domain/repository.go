package domain

type WalletRepository interface {
	Create(wallet *Wallet) error
	GetByID(id string, userID string) (*Wallet, error)
	List(userID string) ([]*Wallet, error)
	Update(wallet *Wallet) error
	Delete(id string, userID string) error
}