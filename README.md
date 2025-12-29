# URL Shortener - Go

A high-performance URL shortening service built with Go, Redis, and Docker. This service provides a RESTful API for creating shortened URLs with custom aliases, expiration times, and built-in rate limiting.

## 🚀 Features

- **URL Shortening**: Convert long URLs into short, shareable links
- **Custom Short URLs**: Create custom short URL aliases
- **Expiration Time**: Set custom expiration times for shortened URLs (default: 24 hours)
- **Rate Limiting**: Built-in IP-based rate limiting to prevent abuse (10 requests per 30 minutes)
- **Redis Storage**: Fast and efficient data storage using Redis
- **URL Validation**: Validates URLs before shortening
- **HTTP Enforcement**: Automatically enforces HTTP/HTTPS protocol
- **Domain Protection**: Prevents shortening of the service's own domain
- **Analytics Counter**: Tracks redirect counts for analytics
- **Dockerized**: Fully containerized for easy deployment

## 📋 Table of Contents

- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Usage Examples](#usage-examples)
- [Rate Limiting](#rate-limiting)
- [Error Handling](#error-handling)

## 🏗️ Architecture

The application follows a clean architecture pattern with the following components:

1. **API Layer**: Built with Fiber (Go web framework) for handling HTTP requests
2. **Database Layer**: Redis for storing URL mappings and rate limit data
3. **Helper Functions**: Utility functions for URL validation and formatting
4. **Route Handlers**: Separate handlers for URL shortening and resolution

### Data Flow

```
Client Request → API Gateway (Fiber) → Rate Limiter → Validator → 
Redis Storage → Response with Short URL
```

## 🛠️ Tech Stack

- **Language**: Go 1.22.0
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber) - Express-inspired web framework
- **Database**: [Redis](https://redis.io/) - In-memory data store
- **Containerization**: Docker & Docker Compose
- **Key Libraries**:
  - `go-redis/redis/v8` - Redis client
  - `govalidator` - URL validation
  - `google/uuid` - UUID generation for short URLs
  - `godotenv` - Environment variable management

## 📁 Project Structure

```
URL-Shortner-GOLANG/
├── api/
│   ├── main.go              # Application entry point
│   ├── go.mod               # Go module dependencies
│   ├── Dockerfile           # API service Dockerfile
│   ├── .env                 # Environment variables
│   ├── database/
│   │   └── database.go      # Redis client configuration
│   ├── helpers/
│   │   └── helpers.go       # Utility functions
│   └── routes/
│       ├── shorten.go       # URL shortening handler
│       └── resolve.go       # URL resolution handler
├── db/
│   └── Dockerfile           # Redis Dockerfile
├── docker-compose.yml       # Docker Compose configuration
└── README.md                # Project documentation
```

## 🔌 API Endpoints

### 1. Shorten URL
Create a shortened URL with optional custom alias and expiration.

**Endpoint**: `POST /api/v1`

**Request Body**:
```json
{
  "url": "https://www.example.com/very/long/url",
  "short": "custom-alias",  // Optional
  "expiry": 48              // Optional (hours, default: 24)
}
```

**Response**:
```json
{
  "url": "https://www.example.com/very/long/url",
  "short": "yourdomain.com/custom-alias",
  "expiry": 48,
  "rate_limit": 9,
  "rate_limit_reset": 30
}
```

### 2. Resolve URL
Redirect to the original URL using the short code.

**Endpoint**: `GET /:url`

**Example**: `GET /abc123`

**Response**: HTTP 301 redirect to the original URL

## 🚦 Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation & Running

1. **Clone the repository**:
   ```bash
   git clone https://github.com/kajal-jotwani/URL-Shortner-GOLANG.git
   cd URL-Shortner-GOLANG
   ```

2. **Create environment file**:
   Create a `.env` file in the `api/` directory with the following variables:
   ```env
   APP_PORT=:3000
   DOMAIN=localhost:3000
   DB_ADDR=db:6379
   DB_PASS=
   API_QUOTA=10
   ```

3. **Build and run with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

4. **Access the service**:
   - API: `http://localhost:3000`
   - Redis: `localhost:6379`

5. **Stop the service**:
   ```bash
   docker-compose down
   ```

## 🔧 Environment Variables

| Variable    | Description                          | Default       |
|-------------|--------------------------------------|---------------|
| `APP_PORT`  | Port where the API server runs       | :3000         |
| `DOMAIN`    | Your domain name                     | localhost:3000|
| `DB_ADDR`   | Redis database address               | db:6379       |
| `DB_PASS`   | Redis password (if any)              | (empty)       |
| `API_QUOTA` | Rate limit quota per IP              | 10            |

## 💡 Usage Examples

### Using cURL

**Shorten a URL**:
```bash
curl -X POST http://localhost:3000/api/v1 \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.github.com/kajal-jotwani",
    "expiry": 24
  }'
```

**Shorten with custom alias**:
```bash
curl -X POST http://localhost:3000/api/v1 \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.github.com/kajal-jotwani",
    "short": "myprofile",
    "expiry": 72
  }'
```

**Access shortened URL**:
```bash
curl -L http://localhost:3000/myprofile
```

### Using JavaScript (Fetch API)

```javascript
// Shorten URL
const response = await fetch('http://localhost:3000/api/v1', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    url: 'https://www.example.com',
    short: 'example',
    expiry: 48
  })
});

const data = await response.json();
console.log(data.short); // Output: localhost:3000/example
```

## ⏱️ Rate Limiting

The service implements IP-based rate limiting to prevent abuse:

- **Quota**: 10 requests per IP address (configurable via `API_QUOTA`)
- **Reset Time**: 30 minutes
- **Storage**: Rate limit data stored in Redis DB #1
- **Response Headers**: 
  - `rate_limit`: Remaining requests
  - `rate_limit_reset`: Time until reset (in minutes)

When rate limit is exceeded, you'll receive:
```json
{
  "error": "Rate limit exceeded",
  "rate_limit_reset": 1800
}
```

## 🛡️ Error Handling

The API handles various error scenarios:

| Error | Status Code | Description |
|-------|-------------|-------------|
| Invalid JSON | 400 | Request body is not valid JSON |
| Invalid URL | 400 | URL format is invalid |
| Custom URL in use | 403 | Custom short URL already exists |
| Rate limit exceeded | 503 | Too many requests from IP |
| Domain error | 503 | Attempting to shorten own domain |
| URL not found | 404 | Short URL doesn't exist |
| Database error | 500 | Cannot connect to Redis |

## 📊 Database Schema

The application uses Redis with two database instances:

**DB 0**: URL Mappings
- Key: Short code (e.g., "abc123")
- Value: Original URL
- TTL: Custom expiry time (default: 24 hours)

**DB 1**: Rate Limiting & Analytics
- Key: IP address → Rate limit counter
- Key: "counter" → Total redirect analytics counter

## 🔐 Security Features

1. **URL Validation**: Validates URLs before shortening
2. **Domain Protection**: Prevents self-referencing URLs
3. **HTTP Enforcement**: Ensures URLs have proper protocol
4. **Rate Limiting**: Prevents API abuse
5. **Custom URL Validation**: Checks for existing custom URLs

## 🐳 Docker Configuration

### Multi-stage Build
The API uses a multi-stage Docker build:
1. **Stage 1**: Build the Go binary
2. **Stage 2**: Create minimal Alpine-based image with the binary

### Services
- **api**: Go application running on port 3000
- **db**: Redis database running on port 6379

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Kajal Jotwani**
- GitHub: [@kajal-jotwani](https://github.com/kajal-jotwani)

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [Redis](https://redis.io/) - Data storage
- [Go](https://golang.org/) - Programming language