package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONResponse(t *testing.T) {

	type TestInterface struct {
		A int    `json:"a"`
		B string `json:"b"`
		C []int  `json:"c"`
	}

	testCases := []struct {
		name                      string
		statusCode                int
		input                     *TestInterface
		expectedBody              interface{}
		expectedStatus            int
		expectedHeaderContentType string
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			input: &TestInterface{
				A: 1,
				B: "text",
				C: []int{1, 2, 3},
			},
			expectedBody:              "{\"a\":1,\"b\":\"text\",\"c\":[1,2,3]}\n",
			expectedStatus:            http.StatusOK,
			expectedHeaderContentType: "application/json",
		},
		{
			name:                      "Nil Input",
			statusCode:                http.StatusNoContent,
			input:                     nil,
			expectedBody:              "null\n",
			expectedStatus:            http.StatusNoContent,
			expectedHeaderContentType: "application/json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			JSONResponse(rr, tc.input, tc.statusCode)

			assert.Equal(t, tc.expectedStatus, rr.Code, "Expected status code to match.")
			assert.Equal(t, tc.expectedHeaderContentType, rr.Header().Get("Content-Type"), "Expected header content-type to match.")
			assert.Equal(t, tc.expectedBody, rr.Body.String(), "Expected body to match.")
		})
	}
}

// package util

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// func TestJSONResponse(t *testing.T) {
// 	type TestInterface struct {
// 		A int    `json:"a"`
// 		B string `json:"b"`
// 	}

// 	tests := []struct {
// 		name           string
// 		statusCode     int
// 		input          interface{}
// 		expectedBody   string
// 		expectedStatus int
// 	}{
// 		{
// 			name:       "Success",
// 			statusCode: http.StatusOK,
// 			input: &TestInterface{
// 				A: 1,
// 				B: "text",
// 			},
// 			expectedBody:   `{"a":1,"b":"text"}`,
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "Nil Input",
// 			statusCode:     http.StatusNoContent,
// 			input:          nil,
// 			expectedBody:   `null`,
// 			expectedStatus: httfunc JSONResponse(w http.ResponseWriter, response interface{}, statusCode int) {
// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(statusCode)
// 				json.NewEncoder(w).Encode(response)
// 			}p.StatusNoContent,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Create a ResponseRecorder to capture the response
// 			rr := httptest.NewRecorder()

// 			// Call the function
// 			JSONResponse(rr, tt.input, tt.statusCode)

// 			// Check the status code
// 			if rr.Code != tt.expectedStatus {
// 				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
// 			}

// 			// Check the response body
// 			body := rr.Body.String()
// 			if body != tt.expectedBody {
// 				t.Errorf("expected body %s, got %s", tt.expectedBody, body)
// 			}

// 			// Check the Content-Type header
// 			contentType := rr.Header().Get("Content-Type")
// 			if contentType != "application/json" {
// 				t.Errorf("expected Content-Type application/json, got %s", contentType)
// 			}
// 		})
// 	}
// }
