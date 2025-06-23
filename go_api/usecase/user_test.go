package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {

	mockUserRepository := new(mockUserRepository)

	testCases := []struct {
		name                     string
		input                    *dto.UserInputDTO
		mockUserRepositoryInput  *domain.User
		mockUserRepositoryReturn *domain.User
		mockUserRepositoryError  error
		expectedOutput           interface{}
		expectedError            interface{}
	}{
		{
			name: "Success",
			input: &dto.UserInputDTO{
				ID:       0,
				Name:     "User1",
				Email:    "user1@example.com",
				Password: "Password@1",
			},
			mockUserRepositoryInput: &domain.User{
				Name:      "User1",
				Email:     "user1@example.com",
				Password:  "Password@1",
				CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockUserRepositoryReturn: &domain.User{
				ID:        1,
				Name:      "User1",
				Email:     "user1@example.com",
				Password:  "Password@1",
				CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockUserRepositoryError: nil,
			expectedOutput: &dto.UserOutputDTO{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
			},
			expectedError: nil,
		},
		{
			name: "Invalid User",
			input: &dto.UserInputDTO{
				ID:       0,
				Name:     "",
				Email:    "user1@example.com",
				Password: "Password@1",
			},
			mockUserRepositoryInput:  nil,
			mockUserRepositoryReturn: nil,
			mockUserRepositoryError:  nil,
			expectedOutput:           nil,
			expectedError:            errors.New("invalid user name"),
		},
		{
			name: "Invalid user data error",
			input: &dto.UserInputDTO{
				ID:       0,
				Name:     "User1",
				Email:    "user1@example.com",
				Password: "Password@1",
			},
			mockUserRepositoryInput: &domain.User{
				Name:      "User1",
				Email:     "user1@example.com",
				Password:  "Password@1",
				CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockUserRepositoryReturn: nil,
			mockUserRepositoryError:  gorm.ErrInvalidData,
			expectedOutput:           nil,
			expectedError:            gorm.ErrInvalidData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository.ExpectedCalls = nil

			mockUserRepository.On("CreateUser", mock.MatchedBy(func(u *domain.User) bool {
				return u.Name == tc.mockUserRepositoryInput.Name &&
					u.Email == tc.mockUserRepositoryInput.Email &&
					bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tc.mockUserRepositoryInput.Password)) == nil
			})).Return(tc.mockUserRepositoryReturn, tc.mockUserRepositoryError)

			userUseCase := NewUserUseCase(mockUserRepository)

			uo, err := userUseCase.CreateUser(tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err, "Expected error")
				assert.Empty(t, uo, "Expected User to be nil")
				assert.Equal(t, err, tc.expectedError, "Expected CreateUser error to match.")

			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.NotEmpty(t, uo, "Expected User not to be nil")
				assert.Equal(t, uo, tc.expectedOutput, "Expected CreateUser output to match.")
			}
		})
	}
}

func TestListUsers(t *testing.T) {

	mockUserRepository := new(mockUserRepository)

	testCases := []struct {
		name                     string
		mockUserRepositoryReturn []*domain.User
		mockUserRepositoryError  error
		expectedOutput           interface{}
		expectedError            interface{}
	}{
		{
			name: "Success",
			mockUserRepositoryReturn: []*domain.User{
				&domain.User{
					ID:        1,
					Name:      "User1",
					Email:     "user1@example.com",
					Password:  "Password@1",
					CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				&domain.User{
					ID:        2,
					Name:      "User2",
					Email:     "user2@example.com",
					Password:  "Password@2",
					CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			mockUserRepositoryError: nil,
			expectedOutput: []*dto.UserOutputDTO{
				&dto.UserOutputDTO{
					ID:        1,
					Name:      "User1",
					Email:     "user1@example.com",
					CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				&dto.UserOutputDTO{
					ID:        2,
					Name:      "User2",
					Email:     "user2@example.com",
					CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedError: nil,
		},
		{
			name:                     "Recover users error",
			mockUserRepositoryReturn: nil,
			mockUserRepositoryError:  gorm.ErrInvalidDB,
			expectedOutput:           nil,
			expectedError:            gorm.ErrInvalidDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository.ExpectedCalls = nil

			mockUserRepository.On("ListUsers").Return(tc.mockUserRepositoryReturn, tc.mockUserRepositoryError)

			userUseCase := NewUserUseCase(mockUserRepository)

			ulo, err := userUseCase.ListUsers()

			if tc.expectedError != nil {
				assert.Error(t, err, "Expected error")
				assert.Empty(t, ulo, "Expected Users List to be nil")
				assert.Equal(t, err, tc.expectedError, "Expected ListUsers error to match.")

			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.NotEmpty(t, ulo, "Expected Users list not to be nil")
				assert.Equal(t, ulo, tc.expectedOutput, "Expected ListUsers output to match.")
			}
		})
	}
}

func TestFindUserById(t *testing.T) {

	mockUserRepository := new(mockUserRepository)

	testCases := []struct {
		name                     string
		input                    uint
		mockUserRepositoryInput  uint
		mockUserRepositoryReturn *domain.User
		mockUserRepositoryError  error
		expectedOutput           interface{}
		expectedError            interface{}
	}{
		{
			name:                    "Success",
			input:                   1,
			mockUserRepositoryInput: 1,
			mockUserRepositoryReturn: &domain.User{
				ID:        1,
				Name:      "User1",
				Email:     "user1@example.com",
				Password:  "Password@1",
				CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockUserRepositoryError: nil,
			expectedOutput: &dto.UserOutputDTO{
				ID:        1,
				Name:      "User1",
				Email:     "user1@example.com",
				CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: nil,
		},
		{
			name:                     "User not found",
			input:                    100,
			mockUserRepositoryInput:  100,
			mockUserRepositoryReturn: nil,
			mockUserRepositoryError:  gorm.ErrRecordNotFound,
			expectedOutput:           nil,
			expectedError:            gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository.ExpectedCalls = nil

			mockUserRepository.On("FindUserById", tc.mockUserRepositoryInput).Return(tc.mockUserRepositoryReturn, tc.mockUserRepositoryError)

			userUseCase := NewUserUseCase(mockUserRepository)

			uo, err := userUseCase.FindUserById(tc.mockUserRepositoryInput)

			if tc.expectedError != nil {
				assert.Error(t, err, "Expected error")
				assert.Empty(t, uo, "Expected Users List to be nil")
				assert.Equal(t, err, tc.expectedError, "Expected FindUserById error to match.")

			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.NotEmpty(t, uo, "Expected Users list not to be nil")
				assert.Equal(t, uo, tc.expectedOutput, "Expected FindUserById output to match.")
			}
		})
	}
}
