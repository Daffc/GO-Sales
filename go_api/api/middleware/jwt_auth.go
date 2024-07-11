package middleware

import (
	"net/http"
	"strings"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/Daffc/GO-Sales/internal/util"
)

type AuthenticationHandler func(http.ResponseWriter, *http.Request, *domain.User)

type JwtAuthenticator struct {
	handler       AuthenticationHandler
	JwtSigningKey []byte
}

func (ja *JwtAuthenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if splitToken[0] != "bearer" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authToken := splitToken[1]

	user, err := util.RecoverUserFromToken(authToken, ja.JwtSigningKey)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ja.handler(w, r, user)
}

func NewJwtAuthenticator(handlerToWrap AuthenticationHandler, jwtSigningKey []byte) *JwtAuthenticator {
	return &JwtAuthenticator{handlerToWrap, jwtSigningKey}
}
