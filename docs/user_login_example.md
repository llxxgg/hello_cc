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