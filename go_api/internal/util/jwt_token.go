package util

import (
	"log"
	"time"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/golang-jwt/jwt"
)

func NewAccessToken(user *domain.User, jwtSigningKey []byte, JwtSessionDuration uint) (string, error) {

	claims := domain.UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(JwtSessionDuration)).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessToken.SignedString(jwtSigningKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return token, nil
}

func IsAuthorized(requestToken string, jwtSigningKey []byte) (bool, error) {
	_, err := jwt.Parse(
		requestToken,
		func(token *jwt.Token) (interface{}, error) {
			return jwtSigningKey, nil
		})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func RecoverUserFromToken(t string, jwtSigningKey []byte) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(
		t,
		&domain.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSigningKey, nil
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userClaims, ok := token.Claims.(*domain.UserClaims)
	if !ok {
		log.Println(err)
		return nil, err
	}

	user := &domain.User{
		ID:    userClaims.ID,
		Name:  userClaims.Name,
		Email: userClaims.Email,
	}

	return user, nil
}
