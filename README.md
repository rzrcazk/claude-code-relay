# Claude Code Relay

一个基于 Go + Gin 框架的 Claude Code 的中转服务, 用于将官方 Claude code账号分发管理以及各类国内支持 Claude code 的代理服务中转

## 🚀 特性

- **简洁架构**: 采用经典分层架构，代码结构清晰
- **完整功能**: 用户管理、任务调度、权限控制、API日志
- **双数据库**: SQLite + Redis，满足不同场景需求
- **中间件**: 认证、限流、CORS、日志等完整中间件
- **RESTful API**: 标准化的API设计
- **自动调度**: 内置任务调度器，自动执行定时任务

## 📁 项目结构

```
claude-code-relay/
├── common/          # 通用工具包
│   ├── logger.go    # 日志工具
│   ├── redis.go     # Redis连接
│   └── utils.go     # 工具函数
├── constant/        # 常量定义
│   └── constant.go  # 系统常量
├── controller/      # 控制器层
│   ├── user.go      # 用户控制器
│   ├── task.go      # 任务控制器
│   └── system.go    # 系统控制器
├── middleware/      # 中间件
│   ├── auth.go      # 认证中间件
│   ├── cors.go      # 跨域中间件
│   ├── logger.go    # 日志中间件
│   ├── rate_limit.go # 限流中间件
│   └── request_id.go # 请求ID中间件
├── model/           # 数据模型
│   ├── database.go  # 数据库初始化
│   ├── user.go      # 用户模型
│   ├── task.go      # 任务模型
│   └── api_log.go   # API日志模型
├── router/          # 路由配置
│   └── router.go    # 路由定义
├── service/         # 服务层
│   └── task_scheduler.go # 任务调度服务
├── web/             # 前端文件
│   └── index.html   # 简单的API文档页面
├── main.go          # 程序入口
├── go.mod           # 依赖管理
└── .env.example     # 环境变量示例
```

## 🛠 快速开始

### 1. 环境要求

- Go 1.23+
- Redis (可选，不配置会跳过)

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库和Redis连接
```

### 4. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:10081` 启动

### 5. 访问系统

- 主页: http://localhost:10081
- API文档: 查看主页上的接口列表
- 默认管理员账户: `admin` / `admin123`

## 📚 API 文档

### 认证接口

#### 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "admin123"
}
```

#### 用户注册
```
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
}
```

### 任务接口

#### 创建任务
```
POST /api/v1/tasks
Authorization: 需要登录
Content-Type: application/json

{
    "title": "测试任务",
    "description": "任务描述",
    "priority": 2,
    "schedule_at": "2024-01-01 10:00:00"
}
```

#### 获取任务列表
```
GET /api/v1/tasks?page=1&limit=10
Authorization: 需要登录
```

## 🔧 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| PORT | 服务端口 | 8080 |
| GIN_MODE | Gin运行模式 | debug |
| DB_PATH | SQLite数据库路径 | ./data/scheduler.db |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| SESSION_SECRET | Session密钥 | your-secret-key-here |

### 数据库

系统使用 SQLite 作为主数据库，自动创建以下表：
- `users` - 用户表
- `tasks` - 任务表
- `api_logs` - API日志表

数据库文件默认存储在 `./data/scheduler.db`

### Redis缓存

Redis用于：
- 限流控制
- Session存储
- 缓存数据

如果未配置Redis，系统会跳过相关功能但不影响正常运行。

## 🎯 核心功能

### 用户管理
- 用户注册、登录、退出
- 基于Session的认证
- 角色权限控制（admin/user）

### 任务调度
- 创建定时任务
- 自动调度执行
- 任务状态管理（pending/running/completed/failed）
- 优先级设置

### 系统管理
- API请求日志记录
- 限流保护
- 系统状态监控
- 管理员仪表板

## 🏗 架构设计

### 分层架构
- **Router层**: 路由定义和中间件配置
- **Controller层**: HTTP请求处理和参数验证
- **Service层**: 业务逻辑处理
- **Model层**: 数据访问和持久化
- **Middleware层**: 横切关注点（认证、日志、限流等）

### 设计原则
- **简单至上**: 代码简洁易懂，避免过度设计
- **职责单一**: 每个模块职责明确
- **易于扩展**: 预留扩展接口，便于功能增强
- **错误处理**: 完善的错误处理和日志记录

## 📝 开发说明

### 添加新接口
1. 在 `model/` 中定义数据模型
2. 在 `controller/` 中实现业务逻辑
3. 在 `router/router.go` 中添加路由
4. 如需中间件，在 `middleware/` 中实现

### 日志管理
- 系统日志: `common.SysLog()` 和 `common.SysError()`
- API日志: 自动记录到数据库
- 日志文件: 存储在 `./logs/app.log`

### 错误处理
统一使用HTTP状态码和错误码：
```json
{
    "error": "错误信息",
    "code": 40001
}
```

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License