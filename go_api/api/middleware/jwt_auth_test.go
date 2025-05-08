package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/Daffc/GO-Sales/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {

	mockAuthenticationHandler := func(w http.ResponseWriter, r *http.Request, u *domain.User) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	mockJwtSigningKey := []byte("test-signing-key")
	mockAuthenticator := NewJwtAuthenticator(mockAuthenticationHandler, mockJwtSigningKey)

	validUser := &domain.User{
		ID:    1,
		Name:  "User1",
		Email: "user1.example.com",
	}

	validToken, _ := util.NewAccessToken(validUser, mockJwtSigningKey, 1)

	testCases := []struct {
		name            string
		authHeader      string
		expectedStatus  int
		mockRecoverUser *domain.User
	}{
		{
			name:            "Valid Token",
			authHeader:      "bearer " + validToken,
			expectedStatus:  http.StatusOK,
			mockRecoverUser: validUser,
		},
		{
			name:            "Missing Authorization Header",
			authHeader:      "",
			expectedStatus:  http.StatusUnauthorized,
			mockRecoverUser: nil,
		},
		{
			name:            "Invalid Authorization Header Format",
			authHeader:      "invalid-header",
			expectedStatus:  http.StatusUnauthorized,
			mockRecoverUser: nil,
		},
		{
			name:            "Invalid Token Type",
			authHeader:      "invalid type",
			expectedStatus:  http.StatusUnauthorized,
			mockRecoverUser: nil,
		},
		{
			name:            "Invalid Token",
			authHeader:      "bearer invalid-token",
			expectedStatus:  http.StatusUnauthorized,
			mockRecoverUser: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Authorization", tc.authHeader)
			rr := httptest.NewRecorder()

			mockAuthenticator.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

		})
	}
}
