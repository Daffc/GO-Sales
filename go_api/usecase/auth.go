package usecase

import (
	"errors"
	"log"
	"time"

	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/Daffc/GO-Sales/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthLoginInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginOutputDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type AuthUseCase struct {
	userRepository       *repository.UserRepository
	JwtSigningKey        []byte
	HoursSessionInterval int8
}

type UserClaims struct {
	ID    uint
	Name  string
	Email string
	jwt.StandardClaims
}

func NewAuthUseCase(userRepository *repository.UserRepository, jwtSigningKey []byte, hoursSessionInterval int8) *AuthUseCase {
	auc := &AuthUseCase{
		userRepository:       userRepository,
		JwtSigningKey:        jwtSigningKey,
		HoursSessionInterval: hoursSessionInterval,
	}
	return auc
}

func (ac AuthUseCase) NewAccessToken(c *UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := accessToken.SignedString(ac.JwtSigningKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (ac AuthUseCase) ValidateAccessToken(t string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		t,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return ac.JwtSigningKey, nil
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		log.Println(err)
		return nil, err
	}

	return claims, nil
}

func (ac AuthUseCase) Login(input *dto.LoginInputDTO) (*dto.LoginOutputDTO, error) {
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

	userClaims := UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(ac.HoursSessionInterval)).Unix(),
		},
	}

	ss, err := ac.NewAccessToken(&userClaims)
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
