package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/Daffc/GO-Sales/domain/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserUseCase struct {
	mock.Mock
}

// CreateUser implements usecase.UserUseCase.
func (m *mockUserUseCase) CreateUser(input *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	args := m.Called(input)
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

// FindUserById implements usecase.UserUseCase.
func (m *mockUserUseCase) FindUserById(input uint) (*dto.UserOutputDTO, error) {
	args := m.Called(input)
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

// UpdateUserPassword implements usecase.UserUseCase.
func (m *mockUserUseCase) UpdateUserPassword(input dto.UpdateUserPasswordInputDTO) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockUserUseCase) ListUsers() ([]*dto.UserOutputDTO, error) {
	args := m.Called()
	return args.Get(0).([]*dto.UserOutputDTO), args.Error(1)
}

func TestCreateUser(t *testing.T) {

	// Create a new mock
	mockUserUseCase := new(mockUserUseCase)

	// Defining test cases mocks and expected results
	testCases := []struct {
		name           string
		requestBody    string
		mockInput      interface{}
		mockReturn     *dto.UserOutputDTO
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:        "Success",
			requestBody: `{"name": "User1", "email": "user1@example.com", "password": "Password@1"}`,
			mockInput: &dto.UserInputDTO{
				Name:     "User1",
				Password: "Password@1",
				Email:    "user1@example.com",
			},
			mockReturn: &dto.UserOutputDTO{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: &dto.UserOutputDTO{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    "{\"invalid\": \"JSON\"",
			mockInput:      nil, // No input should be passed to the mock
			mockReturn:     nil, // The mock should not return anything
			mockError:      nil, // The mock should not be called
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unexpected EOF", // Match the error message returned by the handler
		},
		{
			name:        "Invalid User Name",
			requestBody: "{\n\t\"name\": \"\",\n\t\"email\": \"user1@example.com\",\n\t\"password\": \"Password@1\"\n}",
			mockInput: &dto.UserInputDTO{
				Name:     "",
				Password: "Password@1",
				Email:    "user1@example.com",
			},
			mockReturn:     nil,
			mockError:      fmt.Errorf("invalid user name"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid user name",
		},
	}

	// Runnings tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock to return the user when ListUsers is called
			mockUserUseCase.ExpectedCalls = nil

			// Set up the mock only if mockInput is not nil (error before calling mocked function)
			if tc.mockInput != nil {
				mockUserUseCase.On("CreateUser", tc.mockInput).Return(tc.mockReturn, tc.mockError)
			}

			userHandler := NewUserHandler(mockUserUseCase)

			req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			userHandler.CreateUser(rr, req)

			// Asserting results
			assert.Equal(t, tc.expectedStatus, rr.Code, "Expected status code to match")
			switch rr.Code {
			case http.StatusOK:
				// Check the response body
				var uo *dto.UserOutputDTO
				err = json.NewDecoder(rr.Body).Decode(&uo)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Check the response body
				assert.Equal(t, tc.expectedBody, uo, "Expected user to match")

			case http.StatusBadRequest:
				var response string

				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Check the response error message
				assert.Equal(t, tc.expectedBody, response, "Expected error message to match")
			default:
				t.Fatalf("Unexpected status code: %d", rr.Code)
			}

			// Assert that the mock was called
			mockUserUseCase.AssertExpectations(t)
		})
	}
}

// TestListUsers tests the ListUsers function of the UserHandler.
func TestListUsers(t *testing.T) {
	// Create a new mock user use case
	mockUserUseCase := new(mockUserUseCase)

	// Defining test cases mocks and expected results
	testCases := []struct {
		name           string
		mockReturn     interface{}
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success",
			mockReturn: []*dto.UserOutputDTO{
				{ID: 1, Name: "User1", Email: "user1@example.com"},
				{ID: 2, Name: "User2", Email: "user2@example.com"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: []*dto.UserOutputDTO{
				{ID: 1, Name: "User1", Email: "user1@example.com"},
				{ID: 2, Name: "User2", Email: "user2@example.com"},
			},
		},
		{
			name:           "Bad Request",
			mockReturn:     []*dto.UserOutputDTO{},
			mockError:      fmt.Errorf("Bad Request"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request",
		},
	}

	// Running tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock to return the user when ListUsers is called
			mockUserUseCase.ExpectedCalls = nil
			mockUserUseCase.On("ListUsers").Return(tc.mockReturn, tc.mockError)
			userHandler := NewUserHandler(mockUserUseCase)

			// Create a new HTTP request and response recorder
			req, err := http.NewRequest("GET", "/users", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			userHandler.ListUsers(rr, req)

			// Asserting results
			assert.Equal(t, tc.expectedStatus, rr.Code, "Expected status code to match")
			switch rr.Code {
			case http.StatusOK:
				// Check the response body
				var uos []*dto.UserOutputDTO
				err = json.NewDecoder(rr.Body).Decode(&uos)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Check the response body
				assert.Equal(t, tc.expectedBody, uos, "Expected list of users to match")

			case http.StatusBadRequest:
				var response string
				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Check the response error message
				assert.Equal(t, tc.expectedBody, response, "Expected error message to match")
			default:
				t.Fatalf("Unexpected status code: %d", rr.Code)
			}

			// Assert that the mock was called
			mockUserUseCase.AssertExpectations(t)
		})
	}
}

func TestFindUserById(t *testing.T) {

	// Create a new mock user use case
	mockUserUseCase := new(mockUserUseCase)

	// Defining tests inputs and expected outputs.
	testCases := []struct {
		name           string
		url            string
		mockInput      uint
		mockReturn     *dto.UserOutputDTO
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Success",
			url:       "/users/1",
			mockInput: 1,
			mockReturn: &dto.UserOutputDTO{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
			},
			mockError:      nil,
			expectedStatus: 200,
			expectedBody: &dto.UserOutputDTO{
				ID:    1,
				Name:  "User1",
				Email: "user1@example.com",
			},
		},
		{
			name:           "Invalid Id",
			url:            "/users/X",
			mockInput:      1,
			mockReturn:     nil,
			mockError:      fmt.Errorf("Invalid user id"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid user id",
		},
		{
			name:           "Invalid URL path",
			url:            "/users",
			mockInput:      1,
			mockReturn:     nil,
			mockError:      fmt.Errorf("Invalid URL path"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid URL path",
		},
		{
			name:           "User not found",
			url:            "/users/1",
			mockInput:      1,
			mockReturn:     nil,
			mockError:      fmt.Errorf("record not found"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "record not found",
		},
	}

	// Executing tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Mocking UserUseCase according to test.
			mockUserUseCase.ExpectedCalls = nil
			mockUserUseCase.On("FindUserById", tc.mockInput).Return(tc.mockReturn, tc.mockError)
			userHandler := NewUserHandler(mockUserUseCase)

			// Create a new HTTP request and response recorder
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			userHandler.FindUserById(rr, req)

			// Assert the response status code
			assert.Equal(t, tc.expectedStatus, rr.Code, "Expected status code to match")

			// Assert the response body
			switch rr.Code {
			case http.StatusOK:
				var uo dto.UserOutputDTO
				err = json.NewDecoder(rr.Body).Decode(&uo)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				assert.Equal(t, tc.expectedBody, &uo, "Expected user to match")
			case http.StatusBadRequest:
				var response string
				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Check the response error message
				assert.Equal(t, tc.expectedBody, response, "Expected error message to match")
			default:
				t.Fatalf("Unexpected status code: %d", rr.Code)
			}

			// Assert that the mock was called
			mockUserUseCase.AssertExpectations(t)

		})
	}
}
