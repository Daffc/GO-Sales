package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigParser(t *testing.T) {

	validEnvContent := `
	DB_NAME=dbname
	DB_HOST=localhost
	DB_PORT=1234
	DB_USERNAME=user
	DB_PASSWORD=password
	DB_MIRGATION_FOLDER=./Migration
	SERVER_PORT=3000
	JWT_SIGNING_KEY=SigningKey
	JWT_SESSION_DURATION=1000
	SERVER_WRITE_TIMEOUT=15
	SERVER_READ_TIMEOUT=15
	SERVER_IDLE_TIMEOUT=60`
	validEnvContentFilePath := "./.test.env"
	err := os.WriteFile(validEnvContentFilePath, []byte(validEnvContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(validEnvContentFilePath)

	invalidEnvContent := ``
	invalidEnvContentFilePath := "./.invalid.test.env"
	err = os.WriteFile(invalidEnvContentFilePath, []byte(invalidEnvContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidEnvContentFilePath)

	testCases := []struct {
		name            string
		mockConfig      *Config
		mockEnvFilePath string
		expectedReturn  interface{}
		expectError     bool
	}{
		{
			name: "Success",
			expectedReturn: &Config{
				Database: Database{
					Name:                 "dbname",
					Host:                 "localhost",
					Port:                 "1234",
					User:                 "user",
					Password:             "password",
					MigrationsFolderPath: "./Migration",
				},
				Server: Server{
					Port:               "3000",
					JwtSigningKey:      []byte("SigningKey"),
					JwtSessionDuration: 1000,
					WriteTimeout:       15,
					ReadTimeout:        15,
					IdleTimeout:        60,
				},
			},
			mockEnvFilePath: validEnvContentFilePath,
			expectError:     false,
		},
		{
			name:            "Invalid File Path",
			expectedReturn:  nil,
			mockEnvFilePath: "inexistent_file",
			expectError:     true,
		},
		{
			name:            "Invalid Env Content",
			expectedReturn:  nil,
			mockEnvFilePath: invalidEnvContentFilePath,
			expectError:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			config, erro := NewConfigParser(tc.mockEnvFilePath)

			if tc.expectError {
				assert.Nil(t, config, "Expected config to be nil")
				assert.Error(t, erro, "Expected error")
			} else {
				assert.Equal(t, config, tc.expectedReturn, "Expected config to match")
				assert.NoError(t, erro, "Expected error to be nil")
			}

			// Clearing environment variables
			os.Clearenv()
		})
	}
}
