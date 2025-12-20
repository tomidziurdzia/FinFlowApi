package jwt

import (
	"errors"
	"os"
	"time"

	jwtinterface "fin-flow-api/internal/shared/interface/jwt"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretKey []byte
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewService() jwtinterface.Service {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production" // Default, debe cambiarse
	}
	return &Service{
		secretKey: []byte(secretKey),
	}
}

func (s *Service) GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}