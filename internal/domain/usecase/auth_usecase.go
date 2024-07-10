package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Daffc/GO-Sales/internal/domain/repository"
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
	userRepository *repository.UserRepository
}

type UserClaims struct {
	ID    uint
	Name  string
	Email string
	jwt.StandardClaims
}

func NewAuthUseCase(userRepository *repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepository: userRepository}
}
func (ac AuthUseCase) Login(input AuthLoginInputDTO) (*AuthLoginOutputDTO, error) {
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
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	ss, err := NewAccessToken(&userClaims)
	if err != nil {
		return nil, err
	}

	userDTO := AuthLoginOutputDTO{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Token: ss,
	}

	return &userDTO, nil
}

func NewAccessToken(c *UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := accessToken.SignedString([]byte("SEGREDO123"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateAccessToken(t string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		t,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("SEGREDO123"), nil
		})
	if err != nil {
		log.Fatal(err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		log.Fatal("unknown claims type, cannot proceed")
	}

	fmt.Println(claims.ID, claims.Email, claims.Name, claims.StandardClaims)
	return claims, nil
}
