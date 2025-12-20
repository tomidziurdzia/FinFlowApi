package http

import "fin-flow-api/internal/shared/interface/jwt"

type mockJWTService struct{}

func (m *mockJWTService) GenerateToken(userID string) (string, error) {
	return "mock-token", nil
}

func (m *mockJWTService) ValidateToken(tokenString string) (string, error) {
	return "mock-user-id", nil
}

func newMockJWTService() jwt.Service {
	return &mockJWTService{}
}