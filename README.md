# LWNRA Devotional API

A professional REST API built with Go for managing daily devotionals from Facebook posts.

## 🚀 Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd lwnra-devo-api
go mod tidy

# Set environment variables
export FB_ACCESS_TOKEN="your_facebook_token"
export PORT=8082

# Run the API server
make run
```

The API will be available at `http://localhost:8082`

## 📚 API Endpoints

- `GET /api/devotionals` - Get all devotionals
- `GET /api/devotionals/{date}` - Get devotional by date
- `POST /api/devotionals/sync` - Sync from Facebook
- `POST /api/devotionals/parse` - Parse devotional text
- `GET /api/scheduler/status` - Get scheduler status and next run time
- `GET /health` - Health check

**🤖 Automated Sync**: Devotionals sync automatically daily at 4:45 AM Philippine time!

## 📖 Full Documentation

See [API Documentation](docs/API.md) for complete endpoint details, examples, and usage.

## 🏗️ Project Structure

```
lwnra-devo-api/
├── cmd/server/          # Application entry point
├── handlers/            # HTTP request handlers
├── routes/              # HTTP routing
├── middleware/          # HTTP middleware
├── config/              # Configuration management
├── database/            # Database operations
├── facebook/            # Facebook API client
├── parser/              # Content parsing
├── models/              # Data models
├── docs/                # API documentation
└── Makefile            # Build commands
```

## 🌟 Features

- **Professional API Design**: RESTful endpoints with consistent responses
- **Facebook Integration**: Sync devotionals from Facebook posts
- **Smart Parsing**: Extract structured data from devotional text
- **Bible Version Support**: Handles multiple Bible translations (NIV, ESV, etc.)
- **Date Parsing**: Flexible date extraction from various formats
- **Database Storage**: SQLite with proper schema and relationships
- **CORS Support**: Ready for frontend integration
- **Error Handling**: Comprehensive error responses
- **Logging**: Request/response logging
- **Configuration**: Environment-based settings

## 🛠️ Development

```bash
make help          # Show all commands
make dev           # Development mode with hot reload
make test          # Run tests
make fmt           # Format code
make build         # Build for production
```

## 🚀 Production Ready

- Environment-based configuration
- Structured logging
- Error recovery middleware
- CORS handling
- Health check endpoints
- Clean shutdown handling

## Architecture Overview

Built with professional Go patterns:

- **Clean Architecture**: Separated concerns with handlers, services, and repositories
- **REST API**: Standard HTTP methods and JSON responses
- **Middleware**: Logging, CORS, recovery, and error handling
- **Configuration**: Environment-based configuration management
- **Database**: SQLite with proper schema and migrations
- **Testing**: Unit tests for critical components
- **Development Tools**: Hot reload, formatting, linting
