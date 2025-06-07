# URL Shortener API Endpoints

## Public Endpoints

### Health Check
```bash
GET /health
```
Returns the health status of the application and its dependencies.

### Global Statistics
```bash
GET /stats
```
Returns global statistics (total URLs, clicks, users).

### Create User
```bash
POST /users
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com"
}
```

### Shorten URL
```bash
POST /shorten
Content-Type: application/json

{
  "long_url": "https://example.com",
  "custom_id": "optional-custom-id",
  "expires_at": "2024-12-31T23:59:59Z",
  "click_limit": 100
}
```

### Redirect
```bash
GET /{shortID}
```
Redirects to the original URL.

## Authenticated Endpoints
All authenticated endpoints require the `X-API-Key` header.

### API Key Management

#### Create API Key
```bash
POST /api-keys
X-API-Key: your-existing-api-key
```

#### List API Keys
```bash
GET /api-keys
X-API-Key: your-api-key
```

#### Delete API Key
```bash
DELETE /api-keys
X-API-Key: your-api-key
Content-Type: application/json

{
  "api_key": "api-key-to-delete"
}
```

### URL Management

#### List User URLs
```bash
GET /urls
X-API-Key: your-api-key
```

#### Update URL
```bash
PUT /urls/{shortID}
X-API-Key: your-api-key
Content-Type: application/json

{
  "long_url": "https://new-url.com",
  "expires_at": "2024-12-31T23:59:59Z",
  "click_limit": 200
}
```

#### Delete URL
```bash
DELETE /urls/{shortID}
X-API-Key: your-api-key
```

#### Get Analytics
```bash
GET /analytics/{shortID}
X-API-Key: your-api-key
```

## Example Usage Flow

1. **Create a user:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com"}'
```

2. **Create an API key** (you'll need to insert this manually in the database first, or create a user and API key via SQL):

3. **Create a shortened URL:**
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://google.com"}'
```

4. **List your URLs:**
```bash
curl -X GET http://localhost:8080/urls \
  -H "X-API-Key: your-api-key"
```

5. **Check health:**
```bash
curl -X GET http://localhost:8080/health
```

6. **Get stats:**
```bash
curl -X GET http://localhost:8080/stats
```
