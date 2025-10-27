# URL Shortener

A simple URL shortening service built in Go.

## Features

- Shorten long URLs
- Redirect short codes to original URLs
- RESTful API

## Running

```bash
go run cmd/server/main.go
```

## API Endpoints

### POST /shorten

Shortens a long URL.

**Request:**

```json
{
  "url": "https://example.com/very/long/url"
}
```

**Response (201 Created):**

```json
{
  "short_code": "abc123",
  "short_url": "http://localhost:8080/abc123"
}
```

### GET /:short_code

Redirects to the original URL.

**Response (301 Moved Permanently):** Redirects to the original URL

**Error Response (404 Not Found):** If short code doesn't exist

### GET /health

Health check endpoint.

**Response (200 OK):**

```json
{
  "status": "healthy"
}
```
