# Claude Code Relay

基于Go + Gin的Claude Code 镜像中转服务，采用分层架构设计。支持多账号池管理、智能负载均衡、API Key分发、使用统计和成本计算。包含完整的Vue 3管理界面，基于TDesign组件库。   

## 💡 你能得到什么?

与三五个好友一起拼车使用 `Claude Code` 账号, 同时也可以作为多个 "车主" 为不同的用户提供服务.   

支持任意符合 `Claude Code` API规范的账号池, 如: `GLM4.5` `Qwen3-Code`等, 甚至一些国内的中转镜像站的 `Claude Code` 专属分组均可, 这样就能实现在账号限流的时候智能切换.   

同时提供了完整的使用统计和成本计算, 让你清楚了解每个账号的使用情况和费用支出, 以及为每个Api Key设置每日限额.   

## ✨ 核心特性

**后端服务**
- 多账号池统一管理，智能负载均衡
- 支持Claude官方API和Claude Console双平台
- 基于权重和优先级的智能调度算法
- 完整的Token使用统计和成本计算
- 分层架构设计（Controller-Service-Model）
- 完整中间件链（Auth、CORS、限流、日志等）
- 账号请求异常自动禁用, 定时检测自动恢复
- API Key支持每日限额和可用模型配置

**前端界面** 
- Vue 3 + TypeScript + TDesign组件库
- 响应式管理界面，支持暗黑模式
- 实时数据统计和可视化图表
- 完整的权限管理和用户系统
- 单独的API KEY的用量查询

## 🏗 项目架构

**后端分层结构**
```
├── controller/     # HTTP请求处理、参数验证、响应格式化  
├── service/        # 核心业务逻辑、账号调度、使用统计
├── model/          # 数据模型定义、GORM操作
├── middleware/     # 认证、限流、日志、CORS
├── relay/          # Claude API中转层
├── common/         # 工具函数、成本计算、JWT处理
└── router/         # 路由配置
```

**前端项目结构**
```
web/
├── src/
│   ├── pages/      # 业务页面组件
│   ├── components/ # 公共组件
│   ├── api/        # API请求封装
│   ├── store/      # Pinia状态管理
│   ├── router/     # Vue Router路由
│   └── utils/      # 工具函数
├── package.json    # 依赖配置
└── vite.config.ts  # Vite配置
```

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Node.js 18.18.0+ (前端开发)
- MySQL 8.0+
- Redis (可选，用于限流和缓存)

### 后端启动

```bash
# 1. 安装依赖
go mod tidy

# 2. 配置环境变量
cp .env.example .env
# 生成必需的安全密钥
SESSION_SECRET=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)  
SALT=$(openssl rand -base64 16)

# 3. 启动后端服务（端口8080）
go run main.go
```

### 前端启动

```bash
# 1. 进入前端目录
cd web

# 2. 安装依赖
npm install

# 3. 启动开发服务器（端口3005）
npm run dev

# 或使用Mock数据开发
npm run dev:mock
```

### 访问系统
- 后端API：http://localhost:8080/api/v1/
- 前端界面：http://localhost:3005
- 健康检查：http://localhost:8080/health
- 默认管理员：`admin` / `admin123`

## 📋 核心API

### 认证接口
```bash
# 用户登录
POST /api/v1/auth/login
{"username":"admin","password":"admin123"}

# 用户注册  
POST /api/v1/auth/register
{"username":"user","email":"user@example.com","password":"123456"}
```

### 管理接口（需管理员权限）
```bash
# 账号管理
GET    /api/v1/admin/accounts
POST   /api/v1/admin/accounts
PUT    /api/v1/admin/accounts/{id}
DELETE /api/v1/admin/accounts/{id}

# API Key管理
GET    /api/v1/api-keys
POST   /api/v1/api-keys
DELETE /api/v1/api-keys/{id}

# 分组管理
GET    /api/v1/groups
POST   /api/v1/groups

# 使用统计
GET    /api/v1/admin/logs
```

### Claude中转接口
```bash
# Claude API中转（兼容官方格式）
POST /v1/messages
Authorization: Bearer YOUR_API_KEY

{
  "model": "claude-3-sonnet-20240229",
  "messages": [{"role":"user","content":"Hello"}],
  "max_tokens": 1000
}
```

## ⚙️ 配置说明

### 必需环境变量（生产环境）
```bash
SESSION_SECRET=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)  
SALT=$(openssl rand -base64 16)
```

### 数据库配置
**MySQL（推荐生产环境）**
```bash
MYSQL_HOST=localhost
MYSQL_USER=root
MYSQL_PASSWORD=your-password
MYSQL_DATABASE=claude_code_relay
```

**MySQL数据库**
- 必须配置MySQL相关环境变量
- 支持高并发和大数据量
- 提供更好的性能和可靠性

### Redis缓存（可选）
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your-password  # 可选
```

### 前端环境变量
```bash
# .env.development
VITE_API_URL=http://127.0.0.1:8080
VITE_API_URL_PREFIX=/api/v1
```

## 🔧 核心功能

### 智能账号池
- 多Claude账号统一管理和负载均衡
- 基于权重和优先级的智能调度算法
- 实时状态监控和自动故障转移
- 支持代理配置和平台切换

### 使用统计分析
- 详细的Token使用量统计和成本计算
- 支持多种Claude模型的精确计费
- 实时数据可视化和报告生成
- 账号和API Key级别的使用监控

### 权限管理系统
- 完整的用户注册、登录和权限控制
- Session + JWT双重认证机制
- 灵活的API Key创建和分组管理
- 细粒度的访问权限控制

## 🏗 设计架构

### 后端分层设计
- **Controller层**: 请求处理、参数验证、响应格式化
- **Service层**: 业务逻辑、账号调度、统计计算
- **Model层**: 数据模型、数据库操作、CRUD接口
- **Middleware层**: 认证、限流、CORS、日志记录

### 智能调度算法
1. **优先级排序**: 数字越小优先级越高
2. **权重选择**: 同优先级中按权重比例选择
3. **状态过滤**: 仅选择正常状态的账号
4. **故障转移**: 自动跳过异常账号

### 技术栈
**后端**: Go 1.21+, Gin, GORM, Redis(可选)
**前端**: Vue 3.5+, TypeScript, TDesign, Vite 6+
**数据库**: MySQL 8.0+

## 💻 开发说明

### 后端开发规范
- **分层原则**: Controller → Service → Model
- **错误处理**: 使用 `common.SysLog()` 和 `common.SysError()`
- **用户信息获取**: 使用 `user := c.MustGet("user").(*model.User)`
- **依赖管理**: 添加新依赖后运行 `go mod tidy`

### 前端开发规范
- **组件样式**: 必须声明 `<style scoped>`
- **API请求**: 统一使用 `@/utils/request` 封装的axios
- **代码检查**: 提交前自动运行lint检查
- **开发端口**: 前端3005，自动代理后端8080

### 数据库规范
**核心数据表**
- `users` - 用户账户和角色权限
- `accounts` - Claude账号池和使用统计  
- `api_keys` - API密钥管理和使用监控
- `groups` - 分组管理和权限控制
- `api_logs` - API请求日志和统计数据

## 🐳 部署指南

### Docker部署（推荐）

**一键部署全套服务**
```bash
# 启动MySQL + Redis + 应用
docker-compose -f docker-compose-all.yml up -d

# 查看服务状态
docker-compose -f docker-compose-all.yml ps

# 访问地址
echo "应用地址: http://localhost:10081"
echo "默认管理员: admin / admin123"
```

**使用现有数据库**
```bash
# 复制并编辑环境变量
cp .env.example .env

# 启动应用
docker-compose up -d
```

### 二进制部署

**构建多平台版本**
```bash
# 使用Makefile构建
make build

# 查看构建产物
ls out/
```

**生产环境启动**
```bash
# 设置必需环境变量
export SESSION_SECRET=$(openssl rand -base64 32)
export JWT_SECRET=$(openssl rand -base64 32)
export SALT=$(openssl rand -base64 16)

# 配置数据库（可选）
export MYSQL_HOST=your-host
export MYSQL_USER=your-user
export MYSQL_PASSWORD=your-password

# 启动服务
./claude-code-relay
```

### 反向代理

**Nginx配置示例**
```nginx
server {
    listen 80;
    server_name your-domain.com;
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Caddy配置示例**
```caddyfile
your-domain.com {
    reverse_proxy 127.0.0.1:8080
}
```

## 📋 使用示例

### 1. 管理员登录获取Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 2. 添加Claude账号到账号池
```bash
curl -X POST http://localhost:8080/api/v1/admin/accounts \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Claude官方账号",
    "platform_type": "claude",
    "request_url": "https://api.anthropic.com",
    "secret_key": "sk-your-claude-api-key",
    "priority": 100,
    "weight": 100
  }'
```

### 3. 创建API Key用于中转
```bash
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的中转Key",
    "expires_at": "2025-12-31 23:59:59"
  }'
```

### 4. 通过中转服务调用Claude
```bash
curl -X POST http://localhost:8080/v1/messages \
  -H "Authorization: Bearer YOUR_RELAY_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-sonnet-20240229",
    "messages": [{"role":"user","content":"你好，Claude！"}],
    "max_tokens": 1000
  }'
```

## ❓ 常见问题

**Q: 如何重置管理员密码？**
A: 删除数据库中的admin用户记录，重启服务会自动重新创建默认管理员。

**Q: Redis连接失败是否影响正常使用？**  
A: 不影响核心功能，但会跳过限流和缓存特性。

**Q: 支持哪些Claude模型？**
A: 支持所有Claude模型，包括Claude-3.5系列，成本计算会自动适配不同模型。

**Q: 如何查看账号使用统计？**
A: 通过前端管理界面或API接口查看详细的使用统计和成本分析。

**Q: 前端开发时如何处理跨域？**  
A: 前端开发服务器已配置代理，会自动转发API请求到后端8080端口。

---

## 📄 许可证

MIT License - 欢迎贡献代码和提交Issue！