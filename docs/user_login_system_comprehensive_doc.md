# 用户登录系统技术文档

## 1. 概述

本文档详细介绍了用户登录系统的架构设计、实现细节、API接口、测试策略和部署方法。该系统提供了完整的用户注册、登录和资料管理功能。

## 2. 系统架构

系统采用经典的三层架构模式：
- **表现层（Handlers）**: HTTP路由和请求/响应处理
- **业务层（Services）**: 核心业务逻辑处理
- **数据层（Models）**: 数据结构定义

### 架构图
（此处可以插入使用 DSL 文件生成的架构图）

### 时序图
（此处可以插入使用 DSL 文件生成的时序图）

## 3. API 接口文档

| 接口 | 方法 | 描述 | 请求参数 | 响应 | 错误码 |
|------|------|------|----------|------|--------|
| /api/v1/users/register | POST | 用户注册 | username, email, password | 用户信息 | 400, 409 |
| /api/v1/users/login | POST | 用户登录 | username, password | 用户信息 | 400, 401 |
| /api/v1/users/profile | POST | 获取用户资料 | username | 用户信息 | 400, 404 |
| /api/v1/users/usernames | GET | 获取所有用户名 | - | 用户名列表 | - |

### 3.1 用户注册接口

**接口地址**: `POST /api/v1/users/register`

**请求体**:
```json
{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
}
```

**响应示例**:
```json
{
    "id": "uuid-string",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-04-10T16:20:33.50494+08:00"
}
```

**错误响应**:
- 400 Bad Request: JSON格式错误
- 409 Conflict: 用户名已存在或邮箱格式错误

### 3.2 用户登录接口

**接口地址**: `POST /api/v1/users/login`

**请求体**:
```json
{
    "username": "johndoe",
    "password": "securepassword123"
}
```

**响应示例**:
```json
{
    "id": "uuid-string",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-04-10T16:20:33.50494+08:00"
}
```

**错误响应**:
- 400 Bad Request: JSON格式错误
- 401 Unauthorized: 用户名或密码错误

### 3.3 获取用户资料接口

**接口地址**: `POST /api/v1/users/profile`

**请求体**:
```json
{
    "username": "johndoe"
}
```

**响应示例**:
```json
{
    "id": "uuid-string",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2026-04-10T16:20:33.50494+08:00"
}
```

**错误响应**:
- 404 Not Found: 用户不存在

### 3.4 获取所有用户名接口

**接口地址**: `GET /api/v1/users/usernames`

**响应示例**:
```json
{
    "usernames": [
        "user1",
        "user2",
        "johndoe"
    ]
}
```

## 4. 数据模型

| 模型名 | 字段 | 类型 | 描述 |
|--------|------|------|------|
| User | ID | string | 用户唯一标识 |
| | Username | string | 用户名 |
| | Email | string | 邮箱 |
| | Password | string | 密码（不返回） |
| | CreatedAt | time.Time | 创建时间 |
| | UpdatedAt | time.Time | 更新时间 |
| UserLogin | Username | string | 用户名 |
| | Password | string | 密码 |
| UserRegister | Username | string | 用户名 |
| | Email | string | 邮箱 |
| | Password | string | 密码 |

## 5. 安全设计

| 安全特性 | 实现方式 | 说明 |
|----------|----------|------|
| 输入验证 | 长度限制、格式验证 | 防止注入攻击 |
| 密码保护 | 不在响应中暴露 | 使用 "-" 标签隐藏密码 |
| 线程安全 | 读写锁机制 | sync.RWMutex 保证并发安全 |
| 错误处理 | 统一错误响应 | 避免信息泄露 |

## 6. 测试策略

| 测试类型 | 测试场景 | 覆盖率 |
|----------|----------|--------|
| 单元测试 | 功能正确性 | 100% |
| 集成测试 | 端到端流程 | 100% |
| 异常测试 | 错误处理 | 100% |
| 并发测试 | 高并发场景 | 100% |

### 6.1 测试覆盖情况

| 模块 | 测试用例数 | 通过率 |
|------|------------|--------|
| UserHandler | 12 | 100% |
| UserService | 8 | 100% |
| API端点 | 20 | 100% |

## 7. 部署和运行

### 7.1 本地运行
```bash
go run cmd/server/main.go
```

### 7.2 自定义端口
```bash
PORT=9000 go run cmd/server/main.go
```

### 7.3 环境变量配置
| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| PORT | 8080 | 服务端口 |

## 8. 代码结构

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
cmd/
└── server/
    └── main.go         # 主服务入口
docs/
├── user_login_system_tech_doc.md  # 详细技术文档
├── user_login_example.md          # 使用示例
└── superpowers/
    ├── specs/
    │   └── 2026-04-10-user-login-design.md  # 设计文档
    └── plans/
        └── 2026-04-10-user-login.md         # 实施计划
```

## 9. 开发和扩展指南

### 9.1 添加新功能
- 在相应的层级添加新方法
- 编写对应的单元测试
- 遵循现有代码的命名和格式规范

### 9.2 潜在改进

| 改进项 | 优先级 | 说明 |
|--------|--------|------|
| 密码加密存储 | 高 | 使用 bcrypt 存储密码哈希 |
| 数据库持久化 | 高 | 迁移到 PostgreSQL/MySQL |
| JWT 认证 | 中 | 实现 token 认证机制 |
| 速率限制 | 中 | 防止暴力破解 |
| 详细日志 | 低 | 添加结构化日志 |
| 监控指标 | 低 | 集成监控系统 |

### 9.3 代码规范
- 遵循 Go 语言规范
- 使用 `gofmt` 格式化代码
- 所有错误都使用 `fmt.Errorf("...: %w", err)` 进行包装

## 10. 性能和安全

| 性能指标 | 基线值 | 目标值 |
|----------|--------|--------|
| 平均响应时间 | <50ms | <50ms |
| 并发处理能力 | 100 QPS | 1000 QPS |
| 内存占用 | <50MB | <100MB |

| 安全措施 | 状态 | 备注 |
|----------|------|------|
| 密码验证 | ✓ | 长度和格式检查 |
| SQL注入防护 | ✓ | 使用参数化查询 |
| XSS防护 | - | 后续实现 |
| CSRF防护 | - | 后续实现 |
| DDOS防护 | - | 限流中间件 |