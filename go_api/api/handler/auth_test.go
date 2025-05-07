package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthUseCase struct {
	mock.Mock
}

// CreateUser implements usecase.UserRepository.
func (m *mockAuthUseCase) Login(input *dto.LoginInputDTO) (*dto.LoginOutputDTO, error) {
	args := m.Called(input)
	return args.Get(0).(*dto.LoginOutputDTO), args.Error(1)
}

func TestLogin(t *testing.T) {

	mockAuthUseCase := new(mockAuthUseCase)

	testCases := []struct {
		name           string
		body           string
		mockInput      *dto.LoginInputDTO
		mockReturn     *dto.LoginOutputDTO
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success",
			body: `{"email": "user1@example.com", "password": "Password@1"}`,
			mockInput: &dto.LoginInputDTO{
				Email:    "user1@example.com",
				Password: "Password@1",
			},
			mockReturn: &dto.LoginOutputDTO{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
				Token: "",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: &dto.LoginOutputDTO{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
				Token: "",
			},
		},
		{
			name:           "Invalid JSON",
			body:           `{"Invalid": "JSON"`,
			mockInput:      nil,
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unexpected EOF",
		},
		{
			name: "Wrong Credentials",
			body: `{"email": "user1@example.com", "password": "Password@1"}`,
			mockInput: &dto.LoginInputDTO{
				Email:    "user1@example.com",
				Password: "Password@1",
			},
			mockReturn:     nil,
			mockError:      errors.New("wrong credentials"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "wrong credentials",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock to return the user when ListUsers is called
			mockAuthUseCase.ExpectedCalls = nil

			// Set up the mock only if mockInput is not nil (error before calling mocked function)
			if tc.mockInput != nil {
				mockAuthUseCase.On("Login", tc.mockInput).Return(tc.mockReturn, tc.mockError)
			}

			authHandler := NewAuthHandler(mockAuthUseCase)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			authHandler.Login(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Expected status code to match.")

			switch rr.Code {
			case http.StatusOK:
				lu := &dto.LoginOutputDTO{}
				err := json.NewDecoder(rr.Body).Decode(lu)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.expectedBody, lu)
			case http.StatusBadRequest:
				var r string
				err := json.NewDecoder(rr.Body).Decode(&r)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.expectedBody, r)
			default:
				t.Fatalf("Unexpected status code: %d", rr.Code)
			}

			// Assert that the mock was called
			mockAuthUseCase.AssertExpectations(t)
		})
	}
}
