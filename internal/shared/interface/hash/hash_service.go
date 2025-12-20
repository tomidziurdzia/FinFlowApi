package hash

type Service interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}