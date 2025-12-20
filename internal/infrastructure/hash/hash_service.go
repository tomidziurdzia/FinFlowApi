package hash

import (
	"fin-flow-api/internal/shared/interface/hash"

	"golang.org/x/crypto/bcrypt"
)

type Service struct{}

func NewService() hash.Service {
	return &Service{}
}

func (s *Service) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *Service) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}