package util

import (
	"testing"
	"time"

	"github.com/Daffc/GO-Sales/domain"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {

	testCases := []struct {
		name               string
		user               *domain.User
		jwtSigningKey      []byte
		jwtSessionDuration uint
		expectError        bool
	}{
		{
			name: "Success",
			user: &domain.User{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
			},
			jwtSigningKey:      []byte("TestSigningKey"),
			jwtSessionDuration: uint(1),
			expectError:        false,
		},
		{
			name: "Zero session duration",
			user: &domain.User{
				ID:    2,
				Name:  "User2",
				Email: "user2@example.com",
			},
			jwtSigningKey:      []byte("TestSigningKey"),
			jwtSessionDuration: 0,
			expectError:        false,
		},
		{
			name:               "Nil User",
			user:               nil,
			jwtSigningKey:      []byte(""),
			jwtSessionDuration: 1,
			expectError:        true,
		},
		{
			name: "Empty signing key",
			user: &domain.User{
				ID:    3,
				Name:  "User3",
				Email: "user3@example.com",
			},
			jwtSigningKey:      nil,
			jwtSessionDuration: 0,
			expectError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			token, err := NewAccessToken(tc.user, tc.jwtSigningKey, tc.jwtSessionDuration)
			if tc.expectError {
				assert.Error(t, err, "Expected error")
				assert.Empty(t, token, "Expected token to be empty")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.NotEmpty(t, token, "Expected token to be non-empty")
				parsedToken, err := jwt.ParseWithClaims(
					token,
					&domain.UserClaims{},
					func(token *jwt.Token) (interface{}, error) {
						return tc.jwtSigningKey, nil
					})
				if err != nil {
					t.Fatalf("Failed to parse token: %v", err)
				}

				claims, ok := parsedToken.Claims.(*domain.UserClaims)
				if !ok || !parsedToken.Valid {
					t.Fatalf("Failed to parse token to claims: %v", err)
				}

				assert.Equal(t, tc.user.ID, claims.ID, "Expected user ID to match.")
				assert.Equal(t, tc.user.Name, claims.Name, "Expected user Name to match.")
				assert.Equal(t, tc.user.Email, claims.Email, "Expected user Email to match.")
				assert.Equal(t, uint(time.Duration(int64(time.Second)*(claims.ExpiresAt-claims.IssuedAt)).Hours()), tc.jwtSessionDuration, "Expected token duration to match")
			}
		})
	}
}

func TestRecoverUserFromToken(t *testing.T) {

	jwtSigningKey := []byte("TestSigningKey")
	validUser := &domain.User{
		ID:    1,
		Name:  "User1",
		Email: "user1@example.com",
	}

	claims := domain.UserClaims{
		ID:    validUser.ID,
		Name:  validUser.Name,
		Email: validUser.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := accessToken.SignedString(jwtSigningKey)
	if err != nil {
		t.Fatal()
	}

	expiredClaims := domain.UserClaims{
		ID:    validUser.ID,
		Name:  validUser.Name,
		Email: validUser.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
		},
	}
	accessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, err := accessToken.SignedString(jwtSigningKey)
	if err != nil {
		t.Fatal()
	}

	type InvalidClaims struct {
		InvalidField string
		jwt.StandardClaims
	}
	invalidClaims := InvalidClaims{
		InvalidField: "Invalid value",
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	invalidUserClaimsAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidClaims)
	invalidUserClaimsToken, err := invalidUserClaimsAccessToken.SignedString(jwtSigningKey)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		jwtSigningKey  []byte
		token          string
		expectedOutput interface{}
		expectError    bool
	}{
		{
			name:           "Success",
			token:          validToken,
			jwtSigningKey:  jwtSigningKey,
			expectError:    false,
			expectedOutput: validUser,
		},
		{
			name:           "Expired Token",
			jwtSigningKey:  jwtSigningKey,
			token:          expiredToken,
			expectError:    true,
			expectedOutput: nil,
		},
		{
			name:           "Invalid Token",
			jwtSigningKey:  jwtSigningKey,
			token:          "Invalid Token",
			expectError:    true,
			expectedOutput: nil,
		},
		{
			name:           "Invalid Token User Claim",
			jwtSigningKey:  jwtSigningKey,
			token:          invalidUserClaimsToken,
			expectError:    true,
			expectedOutput: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := RecoverUserFromToken(tc.token, tc.jwtSigningKey)
			if tc.expectError {
				assert.Error(t, err, "Expect Error")
				assert.Nil(t, u, "Expect user to be nil")
			} else {
				assert.Equal(t, tc.expectedOutput, u, "Expected user to match")
			}
		})
	}
}
