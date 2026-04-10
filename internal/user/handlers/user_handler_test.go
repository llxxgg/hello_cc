package handlers

import (
	"bytes"
	"encoding/json"
	"hello_cc/internal/user/models"
	"hello_cc/internal/user/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTest() (*UserHandler, *services.UserService) {
	userService := services.NewUserService()
	handler := NewUserHandler(userService)
	return handler, userService
}

func TestUserHandler_RegisterHandler(t *testing.T) {
	handler, _ := setupTest()

	tests := []struct {
		name           string
		input          models.UserRegister
		expectedStatus int
	}{
		{
			name: "successful registration",
			input: models.UserRegister{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "registration with existing username",
			input: models.UserRegister{
				Username: "testuser",
				Email:    "test2@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "registration with invalid email",
			input: models.UserRegister{
				Username: "newuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.RegisterHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
				t.Logf("response body: %s", w.Body.String())
			}
		})
	}
}

func TestUserHandler_LoginHandler(t *testing.T) {
	handler, userService := setupTest()

	// Pre-register a user
	registerData := models.UserRegister{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "password123",
	}
	_, err := userService.Register(registerData)
	if err != nil {
		t.Fatalf("failed to pre-register user: %v", err)
	}

	tests := []struct {
		name           string
		input          models.UserLogin
		expectedStatus int
	}{
		{
			name: "successful login",
			input: models.UserLogin{
				Username: "loginuser",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "login with invalid password",
			input: models.UserLogin{
				Username: "loginuser",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "login with non-existent user",
			input: models.UserLogin{
				Username: "nonexistent",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.LoginHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
				t.Logf("response body: %s", w.Body.String())
			}
		})
	}
}

func TestUserHandler_ProfileHandler(t *testing.T) {
	handler, userService := setupTest()

	// Pre-register a user
	registerData := models.UserRegister{
		Username: "profileuser",
		Email:    "profile@example.com",
		Password: "password123",
	}
	_, err := userService.Register(registerData)
	if err != nil {
		t.Fatalf("failed to pre-register user: %v", err)
	}

	t.Run("get user profile", func(t *testing.T) {
		type RequestBody struct {
			Username string `json:"username"`
		}

		payload, _ := json.Marshal(RequestBody{Username: "profileuser"})
		req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handler.ProfileHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			t.Logf("response body: %s", w.Body.String())
		}
	})

	t.Run("get non-existent user profile", func(t *testing.T) {
		type RequestBody struct {
			Username string `json:"username"`
		}

		payload, _ := json.Marshal(RequestBody{Username: "nonexistent"})
		req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handler.ProfileHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestUserHandler_AllUsernamesHandler(t *testing.T) {
	handler, userService := setupTest()

	// Pre-register a few users
	testUsers := []struct {
		username string
		email    string
		password string
	}{
		{"user1", "user1@example.com", "password1"},
		{"user2", "user2@example.com", "password2"},
	}

	for _, userData := range testUsers {
		registerData := models.UserRegister{
			Username: userData.username,
			Email:    userData.email,
			Password: userData.password,
		}
		_, err := userService.Register(registerData)
		if err != nil {
			t.Fatalf("failed to pre-register user %s: %v", userData.username, err)
		}
	}

	t.Run("get all usernames", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/usernames", nil)
		w := httptest.NewRecorder()
		handler.AllUsernamesHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string][]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		usernames := response["usernames"]
		if len(usernames) < 2 {
			t.Errorf("expected at least 2 usernames, got %d", len(usernames))
		}
	})
}
