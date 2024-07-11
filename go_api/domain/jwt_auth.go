package domain

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	ID    uint
	Name  string
	Email string
	jwt.StandardClaims
}
