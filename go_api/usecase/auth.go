package usecase

import (
	"errors"

	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/Daffc/GO-Sales/internal/util"
	"github.com/Daffc/GO-Sales/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Login(input *dto.LoginInputDTO) (*dto.LoginOutputDTO, error)
}

type authUseCase struct {
	userRepository     repository.UserRepository
	JwtSigningKey      []byte
	JwtSessionDuration uint
}

func NewAuthUseCase(userRepository repository.UserRepository, jwtSigningKey []byte, jwtSessionDuration uint) AuthUseCase {
	auc := &authUseCase{
		userRepository:     userRepository,
		JwtSigningKey:      jwtSigningKey,
		JwtSessionDuration: jwtSessionDuration,
	}
	return auc
}

func (ac *authUseCase) Login(input *dto.LoginInputDTO) (*dto.LoginOutputDTO, error) {
	user, err := ac.userRepository.FindUserByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.New("wrong credentials")
		default:
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("wrong credentials")
	}

	ss, err := util.NewAccessToken(user, ac.JwtSigningKey, ac.JwtSessionDuration)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	loginOutputDTO := dto.LoginOutputDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: ss,
	}

	return &loginOutputDTO, nil
}
