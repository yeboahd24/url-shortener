# URL Shortener API 🔗

A comprehensive, production-ready URL shortener service built with Go, featuring analytics, user management, API key authentication, and interactive Swagger documentation.

## ✨ Features

- **URL Shortening**: Create short URLs with optional custom IDs
- **User Management**: User registration and API key authentication
- **Analytics**: Click tracking with geolocation data
- **Caching**: Redis-powered caching for high performance
- **Rate Limiting**: Built-in rate limiting for API protection
- **Health Monitoring**: Health checks for all services
- **Interactive Documentation**: Swagger UI for API testing
- **Production Ready**: Docker Compose setup with Nginx reverse proxy

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### 1. Clone the Repository

```bash
git clone https://github.com/yeboahd24/url-shortener
cd url-shortener
```

### 2. Environment Setup

```bash
# Copy environment template
cp .env.example .env

# Edit environment variables
nano .env
```

### 3. Start Services

```bash
# Development
docker-compose up -d

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 4. Access the Application

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

## 📖 API Documentation

### Interactive Documentation
Visit http://localhost:8080/swagger/index.html for interactive API testing.

### Quick API Overview

#### Public Endpoints
- `GET /health` - Health check
- `GET /stats` - Global statistics
- `POST /users` - Create user
- `POST /shorten` - Shorten URL (anonymous)
- `GET /{shortID}` - Redirect to original URL

#### Authenticated Endpoints (require `X-API-Key` header)
- `POST /api/keys` - Create API key
- `GET /api/keys` - List API keys
- `DELETE /api/keys` - Delete API key
- `POST /api/shorten` - Shorten URL (with custom ID)
- `GET /api/urls` - List user URLs
- `PUT /api/urls/{shortID}` - Update URL
- `DELETE /api/urls/{shortID}` - Delete URL
- `GET /api/analytics/{shortID}` - Get analytics

## 🛠️ Development

### Local Development Setup

1. **Install Dependencies**
```bash
go mod download
```

2. **Start Database Services**
```bash
docker-compose up postgres redis -d
```

3. **Run Application**
```bash
go run .
```

### Generate Swagger Documentation
```bash
swag init
```

### Run Tests
```bash
go test ./...
```

## 🏗️ Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Nginx    │────│  Go App     │────│ PostgreSQL  │
│ (Reverse    │    │ (API Server)│    │ (Database)  │
│  Proxy)     │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                   ┌─────────────┐
                   │    Redis    │
                   │   (Cache)   │
                   └─────────────┘
```

### Technology Stack

- **Backend**: Go 1.21+ with Chi router
- **Database**: PostgreSQL 15 with SQLC
- **Cache**: Redis 7
- **Documentation**: Swagger/OpenAPI 3.0
- **Reverse Proxy**: Nginx
- **Containerization**: Docker & Docker Compose

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `POSTGRES_DSN` | PostgreSQL connection string | - |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `REDIS_PASS` | Redis password | - |
| `PORT` | Application port | `8080` |
| `GEO_API_URL` | Geolocation API URL | `http://ip-api.com/json` |

### Database Schema

The application uses the following tables:
- `users` - User accounts
- `urls` - Shortened URLs
- `clicks` - Click tracking
- `api_keys` - API authentication keys

## 📊 Monitoring & Health

### Health Checks
- **Application**: `GET /health`
- **Database**: PostgreSQL health check
- **Cache**: Redis ping check

### Metrics
- Total URLs created
- Total clicks tracked
- Total users registered
- Click analytics by location

## 🚀 Deployment

### Production Deployment

1. **Prepare Environment**
```bash
cp .env.example .env
# Edit .env with production values
```

2. **Deploy with Docker Compose**
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

3. **SSL Configuration** (Optional)
- Place SSL certificates in `./ssl/` directory
- Uncomment HTTPS server block in `nginx.conf`
- Update domain configuration

### Scaling

The application supports horizontal scaling:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --scale app=3
```

## 🔒 Security

- API key authentication for protected endpoints
- Rate limiting (10 requests/minute for API, 100 for redirects)
- Input validation and sanitization
- SQL injection protection with SQLC
- Security headers via Nginx
- Non-root container execution

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check the Swagger UI at `/swagger/index.html`
- **Issues**: Create an issue on GitHub
- **Health Check**: Monitor `/health` endpoint

## 📝 Example Usage

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "email": "john@example.com"}'
```

### Shorten a URL (Anonymous)
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com"}'
```

### Shorten a URL with Custom ID (Authenticated)
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"long_url": "https://example.com", "custom_id": "my-link", "click_limit": 100}'
```

### Get Analytics
```bash
curl -X GET http://localhost:8080/api/analytics/my-link \
  -H "X-API-Key: your-api-key"
```

## 🔍 Troubleshooting

### Common Issues

**Database Connection Failed**
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# View logs
docker-compose logs postgres
```

**Redis Connection Failed**
```bash
# Check Redis status
docker-compose ps redis

# Test Redis connection
docker-compose exec redis redis-cli ping
```

**Application Won't Start**
```bash
# Check application logs
docker-compose logs app

# Verify environment variables
docker-compose config
```

### Performance Tuning

**Database Optimization**
- Ensure proper indexing (included in init.sql)
- Monitor query performance
- Consider connection pooling for high traffic

**Redis Configuration**
- Adjust memory limits based on usage
- Configure persistence settings
- Monitor cache hit rates

## 🎯 Roadmap

- [ ] Rate limiting per user
- [ ] URL expiration cleanup job
- [ ] Bulk URL operations
- [ ] Custom domains support
- [ ] Advanced analytics dashboard
- [ ] API versioning
- [ ] Webhook notifications
- [ ] QR code generation
- [ ] Link preview generation
- [ ] A/B testing for redirects
