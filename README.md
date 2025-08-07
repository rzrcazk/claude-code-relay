# Claude Code Relay

一个基于 Go + Gin 框架的 Claude Code 中转服务，支持多账号管理、负载均衡、API Key 分发、使用统计和成本计算。提供完整的账号池管理功能，支持 Claude 官方 API 和 Claude Console 双平台接入。

## 🚀 特性

- **账号池管理**: 多 Claude 账号统一管理，自动负载均衡
- **双平台支持**: 支持 Claude 官方 API 和 Claude Console
- **API Key 分发**: 灵活的 API Key 管理和分组功能
- **智能调度**: 基于权重和优先级的智能账号调度
- **使用统计**: 详细的 Token 使用量统计和成本计算
- **实时监控**: 账号状态监控、异常检测和自动故障转移
- **完整中间件**: 认证、限流、CORS、日志、请求追踪
- **分层架构**: 清晰的代码结构，易于维护和扩展
- **定时任务**: 自动数据重置、日志清理等定时功能
- **代理支持**: 为每个账号配置独立的代理设置

## 📁 项目结构

```
claude-code-relay/
├── common/              # 通用工具包
│   ├── cost_calculator.go # 成本计算工具
│   ├── jwt.go           # JWT处理工具
│   ├── logger.go        # 日志工具
│   ├── oauth.go         # OAuth认证工具
│   ├── redis.go         # Redis连接工具
│   ├── token_parser.go  # Token解析工具
│   └── utils.go         # 通用工具函数
├── constant/            # 常量定义
│   └── constant.go      # 系统常量和错误码
├── controller/          # 控制器层
│   ├── account.go       # 账号管理控制器
│   ├── api_key.go       # API Key管理控制器
│   ├── claude_code.go   # Claude Code中转控制器
│   ├── group.go         # 分组管理控制器
│   ├── logs.go          # 日志管理控制器
│   ├── system.go        # 系统管理控制器
│   └── user.go          # 用户管理控制器
├── middleware/          # 中间件层
│   ├── auth.go          # 认证中间件
│   ├── claude.go        # Claude专用中间件
│   ├── cors.go          # 跨域中间件
│   ├── logger.go        # 日志记录中间件
│   ├── rate_limit.go    # 限流中间件
│   └── request_id.go    # 请求ID中间件
├── model/               # 数据模型层
│   ├── account.go       # 账号数据模型
│   ├── api_key.go       # API Key数据模型
│   ├── api_log.go       # API日志模型
│   ├── database.go      # 数据库初始化
│   ├── group.go         # 分组数据模型
│   ├── logs.go          # 日志数据模型
│   ├── task.go          # 任务数据模型
│   ├── time.go          # 时间工具模型
│   └── user.go          # 用户数据模型
├── relay/               # Claude API中转层
│   ├── claude.go        # Claude官方API中转
│   └── claude_console.go # Claude Console中转
├── router/              # 路由配置层
│   ├── api_router.go    # RESTful API路由
│   └── claude_router.go # Claude中转路由
├── service/             # 业务逻辑层
│   ├── account.go       # 账号管理服务
│   ├── api_key.go       # API Key管理服务
│   ├── cron_service.go  # 定时任务服务
│   ├── group.go         # 分组管理服务
│   ├── logs.go          # 日志管理服务
│   └── user.go          # 用户管理服务
├── web/                 # 前端文件
│   └── index.html       # 简单的API文档页面
├── main.go              # 程序入口
├── go.mod               # Go模块依赖
├── go.sum               # 依赖校验和
└── .env.example         # 环境变量示例
```

## 🛠 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+ (推荐) 或 SQLite (默认)
- Redis (可选，用于限流和缓存)

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 生成安全密钥（生产环境必须）
openssl rand -base64 32  # 用于SESSION_SECRET
openssl rand -base64 32  # 用于JWT_SECRET  
openssl rand -base64 16  # 用于SALT
```

编辑 `.env` 文件，至少配置以下必需参数：
```bash
SESSION_SECRET=your-session-secret
JWT_SECRET=your-jwt-secret
SALT=your-salt-value
```

### 4. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动（可通过 PORT 环境变量修改）

### 5. 访问系统

- 健康检查: http://localhost:8080/health
- API接口: http://localhost:8080/api/v1/*
- 默认管理员账户: `admin` / `admin123`

## 📚 API 文档

### 认证接口

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "admin123"
}
```

#### 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
}
```

### 账号管理接口 (管理员)

#### 获取账号列表
```http
GET /api/v1/admin/accounts?page=1&limit=10&platform_type=claude
Authorization: Bearer <admin_token>
```

#### 创建账号
```http
POST /api/v1/admin/accounts
Authorization: Bearer <admin_token>
Content-Type: application/json

{
    "name": "Claude账号1",
    "platform_type": "claude",
    "request_url": "https://api.anthropic.com",
    "secret_key": "sk-xxx",
    "group_id": 1,
    "priority": 100,
    "weight": 100
}
```

### API Key管理接口

#### 创建API Key
```http
POST /api/v1/api-keys
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "我的API Key",
    "group_id": 1,
    "expires_at": "2025-01-01 00:00:00"
}
```

#### 获取API Key列表
```http
GET /api/v1/api-keys?page=1&limit=10
Authorization: Bearer <token>
```

### 分组管理接口

#### 创建分组
```http
POST /api/v1/groups
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "默认分组",
    "remark": "分组描述"
}
```

## 🔧 配置说明

### 环境变量

#### 基础配置
| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| PORT | 服务端口 | 8080 | ❌ |
| GIN_MODE | Gin运行模式 (debug/release) | debug | ❌ |
| HTTP_CLIENT_TIMEOUT | HTTP客户端超时时间(秒) | 120 | ❌ |

#### 安全配置（生产环境必须）
| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| SESSION_SECRET | Session密钥 | - | ✅ |
| JWT_SECRET | JWT密钥 | - | ✅ |
| SALT | 密码加密盐值 | - | ✅ |

#### MySQL数据库配置（推荐）
| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| MYSQL_HOST | MySQL主机地址 | localhost | ❌ |
| MYSQL_PORT | MySQL端口 | 3306 | ❌ |
| MYSQL_USER | MySQL用户名 | root | ❌ |
| MYSQL_PASSWORD | MySQL密码 | - | ❌ |
| MYSQL_DATABASE | MySQL数据库名 | claude_code_relay | ❌ |
| MYSQL_MAX_OPEN_CONNS | 最大连接数 | 100 | ❌ |
| MYSQL_MAX_IDLE_CONNS | 最大空闲连接数 | 10 | ❌ |
| MYSQL_MAX_LIFETIME_MINUTES | 连接最大生命周期(分钟) | 60 | ❌ |
| MYSQL_MAX_IDLE_TIME_MINUTES | 连接最大空闲时间(分钟) | 30 | ❌ |

#### Redis配置（可选）
| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| REDIS_HOST | Redis主机地址 | localhost | ❌ |
| REDIS_PORT | Redis端口 | 6379 | ❌ |
| REDIS_PASSWORD | Redis密码 | - | ❌ |
| REDIS_DB | Redis数据库编号 | 0 | ❌ |

#### 日志配置
| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| LOG_LEVEL | 日志级别 | info | ❌ |
| LOG_FILE | 日志文件路径 | ./logs/app.log | ❌ |
| LOG_RECORD_API | 是否记录API日志到数据库 | false | ❌ |
| LOG_RETENTION_MONTHS | 日志保留月数 | 3 | ❌ |

#### 默认管理员配置
| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| DEFAULT_ADMIN_USERNAME | 默认管理员用户名 | admin | ❌ |
| DEFAULT_ADMIN_PASSWORD | 默认管理员密码 | admin123 | ❌ |
| DEFAULT_ADMIN_EMAIL | 默认管理员邮箱 | admin@example.com | ❌ |

### 数据库

系统支持MySQL（推荐）和SQLite两种数据库：

**MySQL模式（生产环境推荐）:**
- 配置MySQL相关环境变量即可启用
- 支持高并发和大数据量
- 提供更好的性能和可靠性

**SQLite模式（开发环境）:**
- 未配置MySQL时自动使用SQLite
- 数据库文件默认存储在 `./data/data.db`
- 无需额外配置，开箱即用

**数据表结构:**
- `users` - 用户表（用户账户、角色权限）
- `accounts` - Claude账号表（账号信息、使用统计、状态监控）
- `api_keys` - API Key表（密钥管理、使用统计、过期时间）
- `groups` - 分组表（账号分组、API Key分组）
- `tasks` - 任务表（任务调度、状态管理）
- `api_logs` - API日志表（请求日志、响应数据、统计信息）

### Redis缓存

Redis用于：
- 限流控制
- Session存储
- 缓存数据

如果未配置Redis，系统会跳过相关功能但不影响正常运行。

## 🎯 核心功能

### 账号池管理
- 多 Claude 账号统一管理
- 支持 Claude 官方 API 和 Claude Console
- 账号状态实时监控（正常/接口异常/账号异常）
- 自动故障转移和负载均衡
- 基于权重和优先级的智能调度
- 代理配置支持

### 使用统计与成本控制
- 详细的 Token 使用量统计
- 实时成本计算（支持多种模型）
- 日/月/年度使用报告
- 缓存Token使用统计
- 账号和API Key级别的使用监控

### API Key 分发管理
- 灵活的 API Key 创建和管理
- 支持过期时间设置
- 分组管理功能
- 使用量限制和监控
- 密钥安全存储

### 用户与权限管理
- 用户注册、登录、权限控制
- 支持 Session 和 JWT 双重认证
- 管理员和普通用户角色区分
- 细粒度权限控制

### 日志与监控
- 完整的 API 请求日志记录
- 实时系统状态监控
- 自动日志清理和归档
- 请求追踪和错误诊断

### 定时任务调度
- 自动重置日使用统计
- 定期日志清理
- 系统维护任务
- 支持 Cron 表达式配置

## 🏗 架构设计

### 分层架构
- **Router层**: RESTful API路由和Claude中转路由配置
- **Controller层**: HTTP请求处理、参数验证、响应格式化
- **Service层**: 核心业务逻辑、账号调度、使用统计
- **Model层**: 数据模型定义、数据库操作、CRUD接口
- **Relay层**: Claude API中转、请求代理、响应处理
- **Middleware层**: 认证、限流、日志、CORS等横切关注点
- **Common层**: 工具函数、成本计算、Token解析等通用功能

### 核心设计模式
- **账号池模式**: 多账号负载均衡和故障转移
- **中转代理模式**: 统一API接口，支持多平台
- **权重调度算法**: 基于优先级和权重的智能选择
- **状态监控模式**: 实时账号状态检测和异常处理
- **分组管理模式**: 灵活的资源分组和权限控制

### 设计原则
- **高可用性**: 多账号备份，自动故障转移
- **可扩展性**: 模块化设计，易于功能扩展
- **安全性**: 密钥加密存储，完整的权限控制
- **可观测性**: 完整的日志记录和监控统计
- **易维护性**: 清晰的分层架构，标准化的代码规范

## 📝 开发说明

### 添加新功能
1. 在 `model/` 中定义数据模型和数据结构
2. 在 `service/` 中实现核心业务逻辑
3. 在 `controller/` 中处理HTTP请求和响应
4. 在 `router/` 中配置路由规则
5. 如需中间件，在 `middleware/` 中实现
6. 如涉及定时任务，在 `service/cron_service.go` 中添加

### 账号调度算法
系统采用基于权重和优先级的智能调度算法：
1. **优先级排序**: 数字越小优先级越高
2. **权重选择**: 在同优先级中按权重比例选择
3. **状态过滤**: 仅选择正常状态的账号
4. **故障转移**: 自动跳过异常账号

### 成本计算
内置支持多种Claude模型的成本计算：
- Input/Output Token分别计费
- 缓存Token（Cache Read/Write）单独计费
- 支持自定义模型价格配置
- 实时USD成本计算

### 日志管理
- **系统日志**: `common.SysLog()` 和 `common.SysError()`
- **API日志**: 自动记录请求/响应到数据库
- **文件日志**: 存储在 `./logs/app.log`
- **日志清理**: 根据 `LOG_RETENTION_MONTHS` 自动清理

### 错误处理
统一的错误响应格式：
```json
{
    "error": "错误描述信息",
    "code": 40001
}
```

常见错误码：
- `40001`: 参数错误
- `40101`: 认证失败
- `40301`: 权限不足
- `40401`: 资源不存在
- `50001`: 服务器内部错误

### 安全考虑
- 密码使用盐值加密存储
- API Key使用安全随机生成
- 支持JWT和Session双重认证
- 敏感信息（Token、密钥）不在日志中记录
- 请求限流防止滥用

## 🚀 部署指南

### Docker Compose部署（推荐）

项目提供两种Docker部署方式：

#### 方式一：使用现有MySQL和Redis服务

如果你已经有MySQL和Redis服务，使用 `docker-compose.yml`：

```bash
# 1. 复制环境变量文件并配置
cp .env.example .env
# 编辑 .env 文件，配置MySQL和Redis连接信息

# 2. 启动应用
docker-compose up -d

# 3. 查看日志
docker-compose logs -f app

# 4. 停止服务
docker-compose down
```

#### 方式二：一键部署全套服务

如果需要同时部署MySQL和Redis，使用 `docker-compose-all.yml`：

```bash
# 1. 复制并配置环境变量（可选，会使用默认值）
cp .env.example .env

# 2. 启动全套服务（包含MySQL、Redis、应用）
docker-compose -f docker-compose-all.yml up -d

# 3. 查看所有服务状态
docker-compose -f docker-compose-all.yml ps

# 4. 查看应用日志
docker-compose -f docker-compose-all.yml logs -f app

# 5. 停止所有服务
docker-compose -f docker-compose-all.yml down
```

**默认配置说明：**
- 应用端口：`10081` (映射到容器内的8080)
- MySQL数据库：`claude_code_relay`
- MySQL用户：`claude` / 密码：`claude123456`
- Redis：无密码，端口6379
- 默认管理员：`admin` / `admin123`

### 手动Docker部署

```bash
# 1. 构建镜像
docker build -t claude-code-relay .

# 2. 运行容器（使用SQLite）
docker run -d \
  --name claude-code-relay \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  -e SESSION_SECRET=$(openssl rand -base64 32) \
  -e JWT_SECRET=$(openssl rand -base64 32) \
  -e SALT=$(openssl rand -base64 16) \
  claude-code-relay

# 3. 运行容器（使用外部MySQL和Redis）
docker run -d \
  --name claude-code-relay \
  -p 8080:8080 \
  -v $(pwd)/logs:/app/logs \
  -e SESSION_SECRET=$(openssl rand -base64 32) \
  -e JWT_SECRET=$(openssl rand -base64 32) \
  -e SALT=$(openssl rand -base64 16) \
  -e MYSQL_HOST=your-mysql-host \
  -e MYSQL_USER=your-mysql-user \
  -e MYSQL_PASSWORD=your-mysql-password \
  -e MYSQL_DATABASE=claude_code_relay \
  -e REDIS_HOST=your-redis-host \
  claude-code-relay
```

### 二进制文件部署

#### 使用Makefile构建

```bash
# 构建多平台二进制文件
make build

# 查看构建产物
ls out/
# claude-code-relay-linux-amd64
# claude-code-relay-windows-amd64.exe
# claude-code-relay-darwin-amd64
# claude-code-relay-linux-arm64
# claude-code-relay-darwin-arm64

# 清理构建产物
make clean
```

#### 生产环境配置

```bash
# 1. 设置环境变量
export GIN_MODE=release
export LOG_LEVEL=info
export SESSION_SECRET=$(openssl rand -base64 32)
export JWT_SECRET=$(openssl rand -base64 32)
export SALT=$(openssl rand -base64 16)

# 2. 配置MySQL（推荐）
export MYSQL_HOST=your-mysql-host
export MYSQL_USER=your-mysql-user
export MYSQL_PASSWORD=your-mysql-password
export MYSQL_DATABASE=claude_code_relay

# 3. 配置Redis（可选）
export REDIS_HOST=your-redis-host
export REDIS_PORT=6379

# 4. 启动服务
./claude-code-relay
```

### 反向代理配置

#### Nginx配置

```nginx
upstream claude_code_relay {
    server 127.0.0.1:8080;
    # 如果有多个实例，可以添加更多server
    # server 127.0.0.1:8081;
}

server {
    listen 80;
    server_name your-domain.com;
    
    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    # SSL证书配置
    ssl_certificate /path/to/your/cert.pem;
    ssl_certificate_key /path/to/your/key.pem;
    
    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # 限制请求体大小
    client_max_body_size 10m;
    
    location / {
        proxy_pass http://claude_code_relay;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # 缓冲区设置
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }
    
    # 健康检查端点
    location /health {
        proxy_pass http://claude_code_relay;
        access_log off;
    }
    
    # 静态文件缓存
    location ~* \.(css|js|jpg|jpeg|png|gif|ico|svg|woff|woff2)$ {
        proxy_pass http://claude_code_relay;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

#### Caddy配置

```caddyfile
your-domain.com {
    reverse_proxy 127.0.0.1:8080
    
    # 安全头
    header {
        X-Frame-Options "SAMEORIGIN"
        X-Content-Type-Options "nosniff"
        X-XSS-Protection "1; mode=block"
    }
    
    # 请求大小限制
    request_body {
        max_size 10MB
    }
    
    # 健康检查
    handle /health {
        reverse_proxy 127.0.0.1:8080
    }
}
```

### 高可用部署

#### 多实例负载均衡

```yaml
# docker-compose-ha.yml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your-root-password
      MYSQL_DATABASE: claude_code_relay
      MYSQL_USER: claude
      MYSQL_PASSWORD: your-password
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

  redis:
    image: redis:7.0-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped

  app1:
    image: claude-code-relay:latest
    ports:
      - "8080:8080"
    env_file: .env
    depends_on: [mysql, redis]
    restart: unless-stopped

  app2:
    image: claude-code-relay:latest
    ports:
      - "8081:8080"
    env_file: .env
    depends_on: [mysql, redis]
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on: [app1, app2]
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

### 监控和维护

#### 健康检查

```bash
# 检查服务状态
curl -f http://localhost:8080/health || exit 1

# 检查数据库连接
curl -s http://localhost:8080/api/v1/status | jq .database_status

# 检查Redis连接
curl -s http://localhost:8080/api/v1/status | jq .redis_status
```

#### 日志管理

```bash
# Docker环境查看日志
docker-compose logs -f app

# 查看应用日志文件
tail -f ./logs/app.log

# 日志轮转配置（logrotate）
cat > /etc/logrotate.d/claude-code-relay << EOF
/path/to/claude-code-relay/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    create 0644 root root
    postrotate
        docker-compose restart app 2>/dev/null || true
    endscript
}
EOF
```

## 🔧 使用示例

### 1. 管理员首次登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 2. 添加Claude账号
```bash
curl -X POST http://localhost:8080/api/v1/admin/accounts \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Claude账号1",
    "platform_type": "claude",
    "request_url": "https://api.anthropic.com",
    "secret_key": "sk-your-claude-key",
    "priority": 100,
    "weight": 100
  }'
```

### 3. 创建API Key
```bash
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试API Key",
    "expires_at": "2025-12-31 23:59:59"
  }'
```

### 4. 使用中转服务
```bash
# 通过中转服务调用Claude API
curl -X POST http://localhost:8080/v1/messages \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-sonnet-20240229",
    "messages": [{"role": "user", "content": "Hello, Claude!"}],
    "max_tokens": 1000
  }'
```

## ❓ 常见问题

### Q: 如何重置管理员密码？
A: 删除数据库中的admin用户，重启服务会自动创建默认管理员账户。

### Q: Redis连接失败是否影响使用？
A: 不影响基本功能，但会跳过限流和缓存功能。

### Q: 如何监控账号使用情况？
A: 通过管理员接口查看账号列表，包含详细的使用统计信息。

### Q: 支持哪些Claude模型？
A: 支持所有Claude模型，包括Claude-3系列，成本计算会根据不同模型自动调整。

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License