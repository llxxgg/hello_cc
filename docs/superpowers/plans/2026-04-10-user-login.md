# 用户登录功能实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现一个基于内存存储的用户登录系统，包括注册、登录和用户资料功能

**Architecture:** 采用三层架构（handler/service/model），使用标准库 net/http 和 JSON 格式

**Tech Stack:** Go (stdlib), memory storage, JSON

---

## 文件结构

- `internal/user/models/user.go` - 定义用户数据结构
- `internal/user/services/user_service.go` - 实现用户业务逻辑
- `internal/user/handlers/user_handler.go` - 实现HTTP处理器
- `internal/user/handlers/user_handler_test.go` - 用户处理器测试
- `cmd/server/main.go` - 主服务入口（示例）

## 任务分解

### Task 1: 创建用户数据模型

**Files:**
- Create: `internal/user/models/user.go`

- [ ] **Step 1: 创建User相关结构体定义**

```go
package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Never expose password in JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserLogin represents login credentials
type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserRegister represents registration data
type UserRegister struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse represents public user information (for API responses)
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User model to a public UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
```

- [ ] **Step 2: 运行 go fmt 格式化代码**

运行: `go fmt ./internal/user/models/user.go`
预期: 格式化完成，无错误

- [ ] **Step 3: 提交代码**

```bash
git add internal/user/models/user.go
git commit -m "feat: add user models with login/register structs

Co-authored-by: Claude Opus 4.6 <noreply@anthropic.com>"
```

### Task 2: 创建用户服务

**Files:**
- Create: `internal/user/services/user_service.go`

- [ ] **Step 1: 创建用户服务结构体和接口**

```go
package services

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"hello_cc/internal/user/models"
	"hello_cc/utils"
)

// UserService manages user-related business logic
type UserService struct {
	users map[string]*models.User // In-memory storage, key is username
	mutex sync.RWMutex           // Thread-safe access to users map
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
```

- [ ] **Step 2: 运行 go fmt 格式化代码**

运行: `go fmt ./internal/user/services/user_service.go`
预期: 格式化完成，无错误

- [ ] **Step 3: 提交代码**

```bash
git add internal/user/services/user_service.go
git commit -m "feat: add user service with register/login logic

Co-authored-by: Claude Opus 4.6 <noreply@anthropic.com>"
```

### Task 3: 创建用户处理器

**Files:**
- Create: `internal/user/handlers/user_handler.go`

- [ ] **Step 1: 创建用户HTTP处理器**

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"hello_cc/internal/user/models"
	"hello_cc/internal/user/services"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterHandler handles user registration
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userReg models.UserRegister
	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userResp, err := h.userService.Register(userReg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResp)
}

// LoginHandler handles user login
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userLogin models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userResp, err := h.userService.Login(userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResp)
}

// ProfileHandler returns user profile (for demonstration, we'll just return based on username from body)
func (h *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Using POST for simplicity, in real app this would be GET with auth
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a real application, this would extract user from JWT or session
	// For this example, we'll take username from request body
	type RequestBody struct {
		Username string `json:"username"`
	}
	
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userResp, err := h.userService.GetUserByUsername(reqBody.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResp)
}

// AllUsernamesHandler returns all registered usernames
func (h *UserHandler) AllUsernamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	usernames := h.userService.GetAllUsernames()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]string{"usernames": usernames})
}
```

- [ ] **Step 2: 运行 go fmt 格式化代码**

运行: `go fmt ./internal/user/handlers/user_handler.go`
预期: 格式化完成，无错误

- [ ] **Step 3: 提交代码**

```bash
git add internal/user/handlers/user_handler.go
git commit -m "feat: add user handlers for register/login/profile

Co-authored-by: Claude Opus 4.6 <noreply@anthropic.com>"
```

### Task 4: 创建用户处理器测试

**Files:**
- Create: `internal/user/handlers/user_handler_test.go`

- [ ] **Step 1: 创建用户处理器测试**

```go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"hello_cc/internal/user/models"
	"hello_cc/internal/user/services"
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
```

- [ ] **Step 2: 运行 go fmt 格式化代码**

运行: `go fmt ./internal/user/handlers/user_handler_test.go`
预期: 格式化完成，无错误

- [ ] **Step 3: 提交代码**

```bash
git add internal/user/handlers/user_handler_test.go
git commit -m "test: add comprehensive tests for user handlers

Co-authored-by: Claude Opus 4.6 <noreply@anthropic.com>"
```

### Task 5: 创建示例主服务

**Files:**
- Create: `cmd/server/main.go`

- [ ] **Step 1: 创建示例服务器**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"hello_cc/internal/user/handlers"
	"hello_cc/internal/user/services"
)

func main() {
	// Initialize services
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	// Define routes
	http.HandleFunc("/api/v1/users/register", userHandler.RegisterHandler)
	http.HandleFunc("/api/v1/users/login", userHandler.LoginHandler)
	http.HandleFunc("/api/v1/users/profile", userHandler.ProfileHandler)
	http.HandleFunc("/api/v1/users/usernames", userHandler.AllUsernamesHandler)

	// Default port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on :%s\n", port)
	fmt.Println("Endpoints:")
	fmt.Println("  POST   /api/v1/users/register - Register a new user")
	fmt.Println("  POST   /api/v1/users/login    - Login a user")
	fmt.Println("  POST   /api/v1/users/profile  - Get user profile")
	fmt.Println("  GET    /api/v1/users/usernames- Get all usernames")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

- [ ] **Step 2: 运行 go fmt 格式化代码**

运行: `go fmt ./cmd/server/main.go`
预期: 格式化完成，无错误

- [ ] **Step 3: 提交代码**

```bash
git add cmd/server/main.go
git commit -m "feat: add example server to demonstrate user login system

Co-authored-by: Claude Opus 4.6 <noreply@anthropic.com>"
```

### Task 6: 创建使用示例和文档

**Files:**
- Create: `docs/user_login_example.md`

- [ ] **Step 1: 创建使用示例文档**

```markdown
# 用户登录系统使用示例

## 启动服务器

```bash
go run cmd/server/main.go
```

或者指定端口：

```bash
PORT=9000 go run cmd/server/main.go
```

## API 端点

### 注册新用户
- **URL**: `/api/v1/users/register`
- **方法**: `POST`
- **内容类型**: `application/json`

示例请求：
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### 用户登录
- **URL**: `/api/v1/users/login`
- **方法**: `POST`
- **内容类型**: `application/json`

示例请求：
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

### 获取用户资料
- **URL**: `/api/v1/users/profile`
- **方法**: `POST` （注意：这里为了简化用POST，实际应使用带认证的GET）
- **内容类型**: `application/json`

示例请求：
```bash
curl -X POST http://localhost:8080/api/v1/users/profile \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe"
  }'
```

### 获取所有用户名
- **URL**: `/api/v1/users/usernames`
- **方法**: `GET`

示例请求：
```bash
curl http://localhost:8080/api/v1/users/usernames
```

## 测试

运行所有测试：
```bash
go test ./...
```

运行特定的处理器测试：
```bash
go test ./internal/user/handlers/
```

## 重要说明

此实现在内存中存储用户数据，仅适用于演示和测试目的。在生产环境中，应：

1. 使用加密哈希（如bcrypt）存储密码
2. 使用数据库（如PostgreSQL、MySQL）存储用户数据
3. 实现适当的认证机制（如JWT或会话）
4. 添加更多输入验证和安全措施
```

- [ ] **Step 2: 运行 go fmt 格式化代码**

运行: `go fmt ./docs/user_login_example.md`
预期: 命令执行成功（虽然Markdown文件无需格式化）

- [ ] **Step 3: 提交代码**

```bash
git add docs/user_login_example.md cmd/server/main.go
git commit -m "docs: add usage examples and documentation for user login system

Co-authored-by: Claude Opus 4.6 <noreply@anthropic.com>"
```

### Task 7: 运行所有测试验证功能

**Files:**
- Test: `./...`

- [ ] **Step 1: 运行所有包的测试**

运行: `go test ./...`
预期: 所有测试通过

- [ ] **Step 2: 如果测试失败，修复代码**

如果测试失败，分析失败原因并修正代码，直到所有测试通过

- [ ] **Step 3: 提交最终版本**

```bash
git add .
git commit -m "feat: complete user login system with tests and documentation

Co-authored-by: Claude Opus 4.6 <noreply@anthropic.com>"
```