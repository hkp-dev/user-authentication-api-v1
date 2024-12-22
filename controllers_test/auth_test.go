package controllers_test

import (
	"app/models"
	"app/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Mock data
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "securepassword123",
	}

	// Mock services
	services.ValidateUser = func(user models.User) []string {
		if user.Username == "" || user.Email == "" || user.Password == "" {
			return []string{"All fields are required"}
		}
		return nil
	}

	services.CheckUserExists = func(user models.User) (bool, error) {
		if user.Username == "existinguser" {
			return true, nil
		}
		return false, nil
	}

	services.HashPassword = func(password string) (string, error) {
		return "hashedpassword", nil
	}

	services.CreateUser = func(user *models.User, cart *models.Cart) error {
		user.ID = "mockedUserID"
		cart.ID = "mockedCartID"
		return nil
	}

	services.CreateCart = func(user models.User) error {
		return nil
	}

	// Test cases
	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful Registration",
			payload: models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "securepassword123",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Registration successful! You can now log in on testuser"}`,
		},
		{
			name: "Missing Fields",
			payload: models.User{
				Username: "",
				Email:    "test@example.com",
				Password: "securepassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":["All fields are required"]}`,
		},
		{
			name: "User Already Exists",
			payload: models.User{
				Username: "existinguser",
				Email:    "exists@example.com",
				Password: "securepassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Username or email already exists"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tc.payload)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Call the handler
			handlers.Register(w, req)

			// Check status code
			res := w.Result()
			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			// Check response body
			responseBody := w.Body.String()
			assert.JSONEq(t, tc.expectedBody, responseBody)
		})
	}
}
