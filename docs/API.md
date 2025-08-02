# LWNRA Devotional API Documentation

A professional REST API for managing daily devotionals from Facebook posts with automated scheduling.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Facebook Access Token (optional, for sync functionality)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd lwnra-devo-api
```

2. **Install dependencies**
```bash
go mod tidy
make install-deps  # Optional: installs development tools
```

3. **Set environment variables**
```bash
export FB_ACCESS_TOKEN="your_facebook_access_token"
export PORT=8080
export ENVIRONMENT=development
```

4. **Start the server**
```bash
make run
# OR for development with hot reload:
make dev
```

**🤖 Automatic Sync**: When FB_ACCESS_TOKEN is set, devotionals will automatically sync from Facebook daily at 4:45 AM Philippine time (UTC+8).

## 📚 API Endpoints

### Base URL
```
http://localhost:8080
```

### Endpoints

#### 1. **Health Check**
```
GET /health
```
**Response:**
```json
{
  "status": "ok",
  "message": "API is running"
}
```

#### 2. **API Information**
```
GET /
```
**Response:**
```json
{
  "name": "LWNRA Devotional API",
  "version": "1.0.0",
  "description": "REST API for managing daily devotionals",
  "endpoints": {
    "GET /api/devotionals": "Get all devotionals (with optional ?limit=N)",
    "GET /api/devotionals/{date}": "Get devotional by date (YYYY-MM-DD format)",
    "POST /api/devotionals/sync": "Sync devotionals from Facebook",
    "POST /api/devotionals/parse": "Parse devotional text"
  }
}
```

#### 3. **Get All Devotionals**
```
GET /api/devotionals
GET /api/devotionals?limit=5
```
**Query Parameters:**
- `limit` (optional): Number of devotionals to return (default: 10)

**Response:**
```json
{
  "success": true,
  "message": "Devotionals retrieved successfully",
  "data": [
    {
      "date": "August 2, 2025",
      "reading": "Matthew 6:16-18",
      "version": "NIV",
      "passage": "16 When you fast, do not look somber...",
      "reflection_qs": [
        "Jesus warns against fasting to be seen (v. 16). What spiritual habit do you do partly for others to notice?"
      ],
      "title": "FASTING IN SECRET",
      "author": "John Smith",
      "body": "Today's devotional focuses on...",
      "prayer": "Lord, help us to seek You..."
    }
  ]
}
```

#### 4. **Get Devotional by Date**
```
GET /api/devotionals/2025-08-02
```
**Response:**
```json
{
  "success": true,
  "message": "Devotional retrieved successfully",
  "data": {
    "date": "August 2, 2025",
    "reading": "Matthew 6:16-18",
    "version": "NIV",
    "passage": "16 When you fast, do not look somber...",
    "reflection_qs": [
      "Jesus warns against fasting to be seen (v. 16). What spiritual habit do you do partly for others to notice?"
    ],
    "title": "FASTING IN SECRET",
    "author": "John Smith",
    "body": "Today's devotional focuses on...",
    "prayer": "Lord, help us to seek You..."
  }
}
```

#### 5. **Sync Devotionals from Facebook**
```
POST /api/devotionals/sync
```
**Response:**
```json
{
  "success": true,
  "message": "Sync completed",
  "data": {
    "synced_count": 2,
    "total_posts": 3,
    "errors": []
  }
}
```

#### 6. **Parse Devotional Text**
```
POST /api/devotionals/parse
```
**Request Body:**
```json
{
  "message": "DAILY DEVOTIONAL\nRead Matthew 6:16-18\nAugust 2, 2025\nMatthew 6:16-18 NIV\n16 When you fast..."
}
```
**Response:**
```json
{
  "success": true,
  "message": "Devotional parsed successfully",
  "data": {
    "date": "August 2, 2025",
    "reading": "Matthew 6:16-18",
    "version": "NIV",
    "passage": "16 When you fast, do not look somber...",
    "reflection_qs": [],
    "title": "",
    "author": "",
    "body": "",
    "prayer": ""
  }
}
```

#### 7. **Get Scheduler Status**
```
GET /api/scheduler/status
```
**Response:**
```json
{
  "success": true,
  "data": {
    "is_running": true,
    "next_run": "2025-08-03T04:45:00+08:00",
    "timezone": "Asia/Manila"
  }
}
```

**Description:** Returns the current status of the automated sync scheduler, including when the next sync is scheduled to run.

## 🤖 Automated Scheduling

The API includes built-in scheduling that automatically syncs devotionals from Facebook:

- **Primary Sync**: Daily at 4:45 AM Philippine Time (UTC+8)
- **Backup Sync**: Daily at 5:15 AM Philippine Time (UTC+8) 
- **Timezone**: Asia/Manila
- **Requires**: FB_ACCESS_TOKEN environment variable

### Scheduler Features

- Graceful startup and shutdown
- Philippine timezone support
- Automatic retry with backup sync
- Status monitoring via API endpoint
- Detailed logging of sync operations

## 🛠️ Development

### Available Commands
```bash
make help          # Show all available commands
make build         # Build the application
make run           # Build and run
make dev           # Development mode with hot reload
make test          # Run tests
make fmt           # Format code
make lint          # Lint code
make clean         # Clean build artifacts
```

### Environment Variables
- `PORT`: Server port (default: 8080)
- `DB_PATH`: Database file path (default: devotionals.db)
- `FB_ACCESS_TOKEN`: Facebook access token for syncing
- `ENVIRONMENT`: Environment (development/production)

### Architecture

```
lwnra-devo-api/
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── handlers/            # HTTP request handlers
├── routes/              # HTTP routing
├── middleware/          # HTTP middleware
├── models/              # Data models
├── database/            # Database operations
├── facebook/            # Facebook API client
├── parser/              # Content parsing
└── Makefile            # Build and development commands
```

## 🔒 Error Handling

All endpoints return consistent error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

### HTTP Status Codes
- `200`: Success
- `400`: Bad Request
- `404`: Not Found
- `500`: Internal Server Error

## 🚀 Production Deployment

### Environment Variables for Production
```bash
export ENVIRONMENT=production
export PORT=8080
export DB_PATH=/app/data/devotionals.db
export FB_ACCESS_TOKEN="your_token"
```

### Build for Production
```bash
make build
./bin/lwnra-devo-api
```

## 📝 Examples

### Using curl

**Get all devotionals:**
```bash
curl http://localhost:8080/api/devotionals
```

**Get devotional by date:**
```bash
curl http://localhost:8080/api/devotionals/2025-08-02
```

**Sync from Facebook:**
```bash
curl -X POST http://localhost:8080/api/devotionals/sync
```

**Parse devotional text:**
```bash
curl -X POST http://localhost:8080/api/devotionals/parse \
  -H "Content-Type: application/json" \
  -d '{"message":"DAILY DEVOTIONAL\nRead Matthew 6:16-18..."}'
```

### Using JavaScript/Frontend

```javascript
// Get all devotionals
const response = await fetch('http://localhost:8080/api/devotionals');
const data = await response.json();
console.log(data.data); // Array of devotionals

// Get devotional by date
const devotional = await fetch('http://localhost:8080/api/devotionals/2025-08-02');
const result = await devotional.json();
console.log(result.data); // Single devotional object
```
