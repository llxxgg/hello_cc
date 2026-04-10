# 用户登录系统技术文档

## 1. 概述

本文档详细介绍了用户登录系统的架构设计、实现细节、API接口、测试策略和部署方法。该系统提供了完整的用户注册、登录和资料管理功能。

## 2. 系统架构

### 2.1 整体架构
系统采用经典的三层架构模式：
- **表现层（Handlers）**: HTTP路由和请求/响应处理
- **业务层（Services）**: 核心业务逻辑处理
- **数据层（Models）**: 数据结构定义

### 2.2 模块划分
```
internal/
├── user/
│   ├── models/         # 数据模型定义
│   │   └── user.go
│   ├── services/       # 业务逻辑
│   │   └── user_service.go
│   └── handlers/       # HTTP处理器
│       ├── user_handler.go
│       └── user_handler_test.go
```

## 3. 数据模型设计

### 3.1 用户实体（User）
```go
type User struct {
    ID        string    `json:"id" db:"id"`
    Username  string    `json:"username" db:"username"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"` // 不暴露到JSON
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### 3.2 输入数据结构
- `UserLogin`: 用户登录凭证结构
- `UserRegister`: 用户注册数据结构  
- `UserResponse`: 响应给客户端的用户信息结构

### 3.3 输出数据结构
- `UserResponse`: 经过脱敏的用户信息，隐藏敏感数据如密码

## 4. 业务逻辑设计

### 4.1 UserService 结构
```go
type UserService struct {
    users map[string]*models.User // 内存存储，key为用户名
    mutex sync.RWMutex           // 线程安全访问锁
}
```

### 4.2 核心功能
- **注册功能**: 检查用户名唯一性，验证邮箱格式，创建新用户
- **登录功能**: 验证用户名和密码，返回用户信息
- **资料获取**: 根据用户名获取用户信息
- **用户名列表**: 获取所有已注册用户名

### 4.3 线程安全机制
使用读写锁（sync.RWMutex）确保在多协程环境下对用户数据的线程安全访问。

## 5. API 接口文档

### 5.1 用户注册
- **端点**: `POST /api/v1/users/register`
- **请求体**:
```json
{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
}
```
- **响应**: 201 Created
```json
{
    "id": "uuid-string",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-04-10T16:20:33.50494+08:00"
}
```
- **错误响应**: 
  - 400 Bad Request: JSON格式错误
  - 409 Conflict: 用户名已存在或邮箱格式错误
  - 422 Unprocessable Entity: 验证失败

### 5.2 用户登录
- **端点**: `POST /api/v1/users/login`
- **请求体**:
```json
{
    "username": "johndoe",
    "password": "securepassword123"
}
```
- **响应**: 200 OK
```json
{
    "id": "uuid-string",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-04-10T16:20:33.50494+08:00"
}
```
- **错误响应**:
  - 400 Bad Request: JSON格式错误
  - 401 Unauthorized: 用户名或密码错误

### 5.3 获取用户资料
- **端点**: `POST /api/v1/users/profile`
- **请求体**:
```json
{
    "username": "johndoe"
}
```
- **响应**: 200 OK
```json
{
    "id": "uuid-string",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-04-10T16:20:33.50494+08:00"
}
```
- **错误响应**:
  - 404 Not Found: 用户不存在

### 5.4 获取所有用户名
- **端点**: `GET /api/v1/users/usernames`
- **响应**: 200 OK
```json
{
    "usernames": [
        "user1",
        "user2",
        "johndoe"
    ]
}
```

## 6. 安全设计

### 6.1 输入验证
- 用户名长度限制（3-32字符）
- 邮箱格式验证
- 密码最小长度限制（6字符）

### 6.2 数据保护
- 密码不在响应中暴露（JSON tag 为 `-`）
- 敏感信息不返回给客户端

### 6.3 并发安全
- 使用读写锁保护共享数据结构
- 确保在高并发情况下的数据一致性

## 7. 测试策略

### 7.1 单元测试覆盖
- 用户注册功能测试（正常情况、重复用户名、无效邮箱）
- 用户登录功能测试（正常登录、错误密码、不存在的用户）
- 用户资料获取测试（正常获取、用户不存在）
- 用户名列表获取测试

### 7.2 测试示例
使用 `httptest` 包模拟HTTP请求，验证各端点的行为。

## 8. 部署和运行

### 8.1 本地运行
```bash
go run cmd/server/main.go
```

### 8.2 自定义端口
```bash
PORT=9000 go run cmd/server/main.go
```

### 8.3 API 测试示例
```bash
# 注册新用户
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 登录
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## 9. 开发和扩展指南

### 9.1 添加新功能
- 在相应的层级添加新方法
- 编写对应的单元测试
- 遵循现有代码的命名和格式规范

### 9.2 潜在改进
1. **密码安全**: 在生产环境中，密码应使用 bcrypt 或其他安全哈希算法存储
2. **持久化存储**: 当前使用内存存储，可扩展为数据库存储（PostgreSQL、MySQL）
3. **身份验证**: 可添加 JWT 或会话管理机制
4. **输入过滤**: 添加更多安全过滤措施防止注入攻击
5. **日志记录**: 集成结构化日志记录
6. **限流机制**: 添加速率限制防止暴力破解

### 9.3 代码规范
- 遵循 Go 语言规范
- 使用 `gofmt` 格式化代码
- 所有错误都使用 `fmt.Errorf("...: %w", err)` 进行包装

## 10. 性能考虑

- 内存存储在小规模应用中性能良好
- 使用读写锁优化读写性能
- 对于大量用户，建议迁移到数据库系统

## 11. 错误处理

- 定义明确的 HTTP 状态码
- 统一的错误响应格式
- 详细的错误信息帮助调试

## 12. 依赖项

- `github.com/google/uuid`: 生成用户唯一ID
- Go 标准库: `net/http`, `encoding/json`, `sync`, `errors`, `time`, `fmt`, `net/mail`
- 内置工具: `testing`, `httptest` 用于测试