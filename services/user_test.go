package services

import (
	"api/models"
	"github.com/stretchr/testify/mock"
)

type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) GenerateToken(profile *models.User) (*string, error) {
	args := m.Called(profile)
	return args.Get(0).(*string), args.Error(1)
}
