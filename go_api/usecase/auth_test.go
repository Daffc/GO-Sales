package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/Daffc/GO-Sales/internal/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {

	validJwtSigningKey := []byte("testJwtSigningKey")
	validJwtSessionDuration := uint(2)

	validLoginCredentials := &dto.LoginInputDTO{
		Email:    "user1@example.com",
		Password: "Password@1",
	}

	validHashedPassword, err := bcrypt.GenerateFromPassword([]byte(validLoginCredentials.Password), 0)
	if err != nil {
		t.Fatal(err)
	}

	validUser := &domain.User{
		ID:        1,
		Name:      "User1",
		Email:     "user1@example.com",
		Password:  string(validHashedPassword),
		CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	validAccessToken, err := util.NewAccessToken(validUser, validJwtSigningKey, validJwtSessionDuration)
	if err != nil {
		t.Fatal(err)
	}

	mockUserRepository := new(mockUserRepository)

	testCases := []struct {
		name                     string
		hashedPassword           []byte
		JwtSigningKey            []byte
		JwtSessionDuration       uint
		loginInput               *dto.LoginInputDTO
		mockUserRepositoryReturn *domain.User
		mockUserRepositoryError  error
		expectedOutput           *dto.LoginOutputDTO
		expectedError            error
	}{
		{
			name:                     "Success",
			hashedPassword:           validHashedPassword,
			JwtSigningKey:            validJwtSigningKey,
			JwtSessionDuration:       validJwtSessionDuration,
			loginInput:               validLoginCredentials,
			mockUserRepositoryReturn: validUser,
			mockUserRepositoryError:  nil,
			expectedOutput: &dto.LoginOutputDTO{
				ID:        validUser.ID,
				Name:      validUser.Name,
				Email:     validUser.Email,
				CreatedAt: validUser.CreatedAt,
				UpdatedAt: validUser.UpdatedAt,
				Token:     validAccessToken,
			},
			expectedError: nil,
		},
		{
			name:                     "Credentials Not Found",
			hashedPassword:           validHashedPassword,
			JwtSigningKey:            validJwtSigningKey,
			JwtSessionDuration:       validJwtSessionDuration,
			loginInput:               validLoginCredentials,
			mockUserRepositoryReturn: nil,
			mockUserRepositoryError:  gorm.ErrRecordNotFound,
			expectedOutput:           nil,
			expectedError:            errors.New("wrong credentials"),
		},
		{
			name:                     "Others Recover User Errors",
			hashedPassword:           validHashedPassword,
			JwtSigningKey:            validJwtSigningKey,
			JwtSessionDuration:       validJwtSessionDuration,
			loginInput:               validLoginCredentials,
			mockUserRepositoryReturn: nil,
			mockUserRepositoryError:  gorm.ErrInvalidData,
			expectedOutput:           nil,
			expectedError:            gorm.ErrInvalidData,
		},
		{
			name:               "Wrong Password",
			hashedPassword:     validHashedPassword,
			JwtSigningKey:      validJwtSigningKey,
			JwtSessionDuration: validJwtSessionDuration,
			loginInput: &dto.LoginInputDTO{
				Email:    "user1@example.com",
				Password: "WrondPassword",
			},
			mockUserRepositoryReturn: validUser,
			mockUserRepositoryError:  nil,
			expectedOutput:           nil,
			expectedError:            errors.New("wrong credentials"),
		},
		{
			name:               "Wrong Password",
			hashedPassword:     validHashedPassword,
			JwtSigningKey:      validJwtSigningKey,
			JwtSessionDuration: validJwtSessionDuration,
			loginInput: &dto.LoginInputDTO{
				Email:    "user1@example.com",
				Password: "WrondPassword",
			},
			mockUserRepositoryReturn: validUser,
			mockUserRepositoryError:  nil,
			expectedOutput:           nil,
			expectedError:            errors.New("wrong credentials"),
		},
		{
			name:                     "Erro Token Generation",
			hashedPassword:           validHashedPassword,
			JwtSigningKey:            []byte(""),
			JwtSessionDuration:       validJwtSessionDuration,
			loginInput:               validLoginCredentials,
			mockUserRepositoryReturn: validUser,
			mockUserRepositoryError:  nil,
			expectedOutput:           nil,
			expectedError:            errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepository.ExpectedCalls = nil
			mockUserRepository.On("FindUserByEmail", tc.loginInput.Email).Return(tc.mockUserRepositoryReturn, tc.mockUserRepositoryError)

			authUseCase := NewAuthUseCase(mockUserRepository, tc.JwtSigningKey, tc.JwtSessionDuration)

			lod, err := authUseCase.Login(tc.loginInput)

			assert.Equal(t, lod, tc.expectedOutput, "Expected Logind output to match.")
			assert.Equal(t, err, tc.expectedError, "Expected Logind error to match.")

			mockUserRepository.AssertExpectations(t)
		})
	}
}
