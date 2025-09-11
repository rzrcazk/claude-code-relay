# Claude Code Relay


<div align="center">


[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Node.js](https://img.shields.io/badge/Node.js-18+-green.svg)](https://nodejs.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-red.svg)](https://redis.io/)
[![Redis](https://img.shields.io/badge/Mysql-5.7+-yellow.svg)](https://www.mysql.com/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)


A Claude Code API relay service based on Go + Gin with layered architecture design. Supports multi-account pool management, intelligent load balancing, API key distribution, usage statistics, and cost calculation. Includes a complete Vue 3 admin interface built with TDesign component library.   


**If this project helps you, please give me a ⭐️!**   

[Chinese](README.md)  •  [English](README_EN.md)     

</div>


---

![home.png](docs/home.png)

## ⚠️ Important Notice
Please read carefully before using this project:

🚨 Terms of Service Risk: Using this project may violate Anthropic's Terms of Service. Please carefully read Anthropic's user agreement before use. All risks from using this project are borne by the user.

📖 Disclaimer: This project is for technical learning and research purposes only. The author assumes no responsibility for account bans, service interruptions, or other losses caused by using this project.

## 💡 What Can You Get?

Share a `Claude Code` account with friends for cost-effective usage, or serve as multiple "car owners" providing services to different users.

Supports any account pool that complies with `Claude Code` API specifications, such as: `GLM4.5`, `Qwen3-Code`, and even some domestic relay mirror sites' `Claude Code` dedicated groups, enabling intelligent switching when accounts are rate-limited.

Also provides complete usage statistics and cost calculation, allowing you to clearly understand each account's usage and expenses, and set daily limits for each API Key.

### Supported Platform Types
- ✅ Supports Claude official accounts (requires Pro+ subscription)
- ✅ Supports any Claude Code mirror interfaces (official mirrors/GLM/Qwen, etc.)
- ✅ Supports any OpenAI API format compatible interfaces

## ✨ Core Features

**Backend Services**
- Unified multi-account pool management with intelligent load balancing
- Supports both Claude official API and Claude Console platforms
- Smart scheduling algorithm based on weight and priority
- Complete token usage statistics and cost calculation
- Layered architecture design (Controller-Service-Model)
- Complete middleware chain (Auth, CORS, rate limiting, logging, etc.)
- Automatic account disabling on request errors with scheduled recovery
- API Key daily limits and available model configuration

**Frontend Interface**
- Vue 3 + TypeScript + TDesign component library
- Real-time data statistics and visualization charts
- Complete permission management and user system
- Individual API KEY usage query (`/stats/api-key?api_key=sk-xxx`)

## 🏗 Project Architecture

**Backend Layered Structure**
```
├── controller/     # HTTP request handling, parameter validation, response formatting
├── service/        # Core business logic, account scheduling, usage statistics
├── model/          # Data model definitions, GORM operations
├── middleware/     # Authentication, rate limiting, logging, CORS
├── relay/          # Claude API relay layer
├── common/         # Utility functions, cost calculation, JWT handling
└── router/         # Route configuration
```

**Frontend Project Structure**
```
web/
├── src/
│   ├── pages/      # Business page components
│   ├── components/ # Common components
│   ├── api/        # API request wrappers
│   ├── store/      # Pinia state management
│   ├── router/     # Vue Router configuration
│   └── utils/      # Utility functions
├── package.json    # Dependencies configuration
└── vite.config.ts  # Vite configuration
```

## 🚀 Quick Start

### Environment Requirements
- Go 1.21+
- Node.js 18.18.0+ (for frontend development)
- MySQL 5.7+
- Redis

## 📋 Core API
- [Apifox Online Documentation](https://s.apifox.cn/ba2f5ebd-5a13-4e3a-9c42-628208b1d09f) covers most interfaces

## 🏗 Design Architecture

### Backend Layered Design
- **Controller Layer**: Request handling, parameter validation, response formatting
- **Service Layer**: Business logic, account scheduling, statistics calculation
- **Model Layer**: Data models, database operations, CRUD interfaces
- **Middleware Layer**: Authentication, rate limiting, CORS, logging
- **Scheduled Layer**: Account status detection, automatic recovery of abnormal accounts

### Intelligent Scheduling Algorithm
1. **Priority Sorting**: Lower numbers have higher priority
2. **Weight Selection**: Proportional selection by weight within same priority
3. **Status Filtering**: Only select accounts in normal status
4. **Failover**: Automatically skip abnormal accounts

### Technology Stack
1. **Backend**: Go 1.21+, Gin, GORM, Redis
2. **Frontend**: Vue 3.5+, TypeScript, TDesign, Vite 6+
3. **Database**: MySQL 5.7+, Redis

## 💻 Development Guide

### Backend Development Standards
- **Layering Principle**: Controller → Service → Model
- **Error Handling**: Use `common.SysLog()` and `common.SysError()`
- **User Info Retrieval**: Use `user := c.MustGet("user").(*model.User)`
- **Dependency Management**: Run `go mod tidy` after adding new dependencies

### Frontend Development Standards
- **Component Styles**: Must declare `<style scoped>`
- **API Requests**: Use unified `@/utils/request` wrapped axios
- **Code Checking**: Automatic lint checking before commits
- **Documentation Reference**: [TDesign Vue3 Documentation](https://tdesign.tencent.com/vue-next/)

### Database Standards
**Core Data Tables**
- `users` - User accounts and role permissions
- `accounts` - Claude account pool and usage statistics
- `api_keys` - API key management and usage monitoring
- `groups` - Group management and access control
- `logs` - Model request log data
- `api_logs` - API request logs and statistics

## 🐳 Deployment Guide

### Docker Deployment (Recommended)

**One-Click Full Service Deployment**
```bash
# Start MySQL + Redis + Application
docker-compose -f docker-compose-all.yml up -d

# Check service status
docker-compose -f docker-compose-all.yml ps

# Access URLs
echo "Application URL: http://localhost:10081"
echo "Default admin: admin / admin123"
```

**Using Existing Database**
```bash
# Copy and edit environment variables
cp .env.example .env

# Start application
docker-compose up -d
```

### Binary Deployment

**Build Multi-Platform Versions**
```bash
# Build using Makefile
make build

# Check build artifacts
ls out/
```

**Production Environment Startup**
```bash
# Set required environment variables
export SESSION_SECRET=$(openssl rand -base64 32)
export JWT_SECRET=$(openssl rand -base64 32)
export SALT=$(openssl rand -base64 16)
...

# Configure database
export MYSQL_HOST=your-host
export MYSQL_USER=your-user
export MYSQL_PASSWORD=your-password
...

# Start service (The frontend files for version v0.1.10 and higher have been packaged into the binary; simply open the service domain to access the homepage.)
./claude-code-relay
```

## 💐 How to Use This Service in Claude Code?

Using this service in Claude Code is very simple - just replace Claude Code's request address with this service address and key. Here are the specific steps:

```bash
# Configure API request address in Claude Code (HTTPS recommended)
export ANTHROPIC_BASE_URL=https://your-server-domain/claude-code

# Configure API key in Claude Code
export ANTHROPIC_AUTH_TOKEN="your-api-key"
```

Alternative method - create and configure Settings file: Create ~/.claude/settings.json file and configure your API key:
```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "your-api-key-here",
    "ANTHROPIC_BASE_URL": "https://your-server-domain/claude-code",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 1
  },
  "permissions": {
    "allow": [],
    "deny": []
  }
}
```

**For more detailed instructions, go to the admin panel, click "Help Documentation" in the top right corner, or directly visit the `/help/index` page for comprehensive usage tutorials.**
![help.png](docs/help.png)

## ❓ FAQ

**Q: How to reset admin password?**
A: Delete the admin user record from the database, restart the service to automatically recreate the default admin.

**Q: Which Claude models are supported?**
A: Supports all Claude models, including Claude-3.5 series. Cost calculation automatically adapts to different models.

**Q: How to view account usage statistics?**
A: View detailed usage statistics and cost analysis through the frontend admin interface or API endpoints.

**Q: Cannot access Claude official services normally?**
A: Claude and other foreign model service providers block domestic IP access. Use a proxy to resolve this (high-quality proxy IP recommended).

## 🤝 Acknowledgments
- Inspiration from: [claude-relay-service](https://github.com/Wei-Shaw/claude-relay-service)
- **90%** of this project's code was developed by [Claude Code](https://www.anthropic.com/claude-code). Thanks to Anthropic for providing powerful AI capabilities.

---

## 📄 License

[MIT License](LICENSE) - Contributions and issues are welcome!