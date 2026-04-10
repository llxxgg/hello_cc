package services

import (
	"errors"
	"fmt"
	"net/mail"
	"hello_cc/internal/user/models"
	"hello_cc/utils"
	"sync"
	"time"
)

// UserService manages user-related business logic
type UserService struct {
	users map[string]*models.User // In-memory storage, key is username
	mutex sync.RWMutex            // Thread-safe access to users map
}

// NewUserService creates a new user service instance
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*models.User),
	}
}

// Register registers a new user
func (s *UserService) Register(userReg models.UserRegister) (*models.UserResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if user already exists
	if _, exists := s.users[userReg.Username]; exists {
		return nil, errors.New("username already exists")
	}

	// Validate email format
	if !isValidEmail(userReg.Email) {
		return nil, errors.New("invalid email format")
	}

	// Create new user
	now := time.Now()
	user := &models.User{
		ID:        utils.GenerateUUID(),
		Username:  userReg.Username,
		Email:     userReg.Email,
		Password:  userReg.Password, // In production, this should be hashed
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Store user
	s.users[user.Username] = user

	response := user.ToResponse()
	return &response, nil
}

// Login authenticates a user and returns user info
func (s *UserService) Login(loginReq models.UserLogin) (*models.UserResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Find user by username
	user, exists := s.users[loginReq.Username]
	if !exists {
		return nil, errors.New("invalid username or password")
	}

	// Check password (in memory, comparing directly)
	if user.Password != loginReq.Password {
		return nil, errors.New("invalid username or password")
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*models.UserResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", username)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetAllUsernames returns all registered usernames
func (s *UserService) GetAllUsernames() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	usernames := make([]string, 0, len(s.users))
	for username := range s.users {
		usernames = append(usernames, username)
	}
	return usernames
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
