package usecase

import (
	"github.com/Daffc/GO-Sales/domain"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(u *domain.User) (*domain.User, error) {
	args := m.Called(u)
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserRepository) ListUsers() ([]*domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *mockUserRepository) FindUserById(id uint) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) FindUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) UpdateUserPassword(u *domain.User) error {
	args := m.Called(u)
	return args.Error(1)
}
