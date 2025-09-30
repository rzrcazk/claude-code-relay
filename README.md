# Claude Code Relay

A proxy service for Claude API requests with features like account pooling, load balancing, usage statistics, and cost calculation.

## Features

- **Multi-Platform Support**: Proxy requests to Claude official API, Claude Console, and OpenAI-compatible APIs
- **Account Pooling**: Manage multiple Claude accounts with priority-based load balancing
- **API Key System**: Secure access control with daily usage limits and model restrictions
- **Rate Limiting**: Both at the service level and automatic handling of upstream rate limits
- **Usage Statistics**: Detailed token usage tracking and cost calculation for each account and API key
- **Cost Calculation**: Real-time cost tracking with support for Claude's cache tokens pricing
- **Account Health Monitoring**: Automatic recovery of abnormal accounts and rate limit management
- **Group Management**: Organize accounts and API keys into logical groups
- **Dashboard Analytics**: Comprehensive usage statistics and trends visualization

## Architecture

The service is built with:
- **Backend**: Go + Gin web framework
- **Frontend**: Vue 3 + TypeScript with TDesign components
- **Database**: MySQL with GORM ORM
- **Caching**: Redis for session storage and caching
- **Deployment**: Single binary with embedded frontend assets

## Getting Started

1. Copy `.env.example` to `.env` and configure your settings
2. Run `go run main.go` to start the development server
3. Visit `http://localhost:8080` to access the web interface

For production deployment, build the project with `make build` which will create binaries for multiple platforms.

## Documentation

See [CLAUDE.md](CLAUDE.md) for detailed project documentation and development guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.