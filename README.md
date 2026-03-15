# echte.link - URL Shortener & Linktree Service

A lightweight URL shortener and Linktree-style service built with Go, optimized for small servers (0.5 vCPU, 500MB RAM).

## Features

- **UUID-based Authentication**: Mullvad-style user management with no email recovery
- **URL Shortening**: Create short URLs with `/#code` format
- **Linktree Pages**: Personal link pages with `/@username` format
- **SQLite Database**: Lightweight, file-based database
- **Docker Ready**: Optimized for small container deployments
- **Resource Efficient**: Designed for minimal resource usage

## URL Structure

- `/` - Homepage with URL shortener and account creation
- `/@username` - User's Linktree-style profile page
- `/#code` - Short URL redirect (public or user-owned)
- `/@username/#code` - User-owned short URL redirect

## API Endpoints

- `POST /api/user/create` - Create new user account
- `POST /api/user/auth` - Authenticate user
- `POST /api/links` - Create profile link
- `POST /api/shorten` - Create short URL
- `GET /api/profile/:username` - Get user profile

## Deployment

### Docker (Recommended)

```bash
# Build and run
docker-compose up -d

# Or build manually
docker build -t echte-link .
docker run -d -p 8080:8080 -v $(pwd)/data:/root/data echte-link
```

### Local Development

```bash
# Install dependencies
go mod tidy

# Run server
go run main.go
```

## Database Schema

SQLite database (`echter.link.sqlite`) with three tables:

- `users`: UUID-based user accounts
- `links`: User profile links (Linktree functionality)
- `short_urls`: Shortened URLs with click tracking

## Security

- Passwords hashed with bcrypt
- UUID-based user IDs (no personal information)
- No email storage or recovery options
- SQLite file permissions should be secured

## Resource Usage

- **CPU**: Optimized for 0.5 vCPU
- **Memory**: ~250MB base usage, 500MB limit
- **Storage**: Minimal, grows with user data
- **Network**: HTTP/HTTPS only

## Domain Configuration

Configure your reverse proxy (nginx/caddy) to forward `echter.link` to port 8080.

Example nginx config:
```nginx
server {
    server_name echter.link;
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
    }
}
```
