# echter.link - Modern URL Shortener!

A lightweight URL shortener built with Go, featuring Dark Mode UI, interactive mouse effects, and snow animation. Optimized for small servers (0.5 vCPU, 500MB RAM).

## ✨ Features

- **🚀 URL Shortening**: Create short URLs with random or custom codes
- **🌙 Dark Mode**: Modern dark theme with gradient accents
- **✨ Mouse Glow Effect**: Interactive light source following cursor
- **❄️ Snow Animation**: Falling snow with melt interaction
- **📱 Responsive Design**: Mobile-friendly interface
- **🔗 Input Group**: Visual `https://` prefix for better UX
- **📊 Statistics**: Demo statistics in footer
- **⏰ Expiration**: Optional URL expiration dates
- **✅ URL Validation**: Anti-phishing protection
- **🐳 Docker Ready**: Optimized container deployment
- **💾 SQLite**: Lightweight, file-based database

## 🎯 Live Demo

Visit `http://localhost:8080` to see the interactive dark mode interface with mouse glow and snow effects.

## 🚀 Quick Start

### Docker (Recommended)

```bash
# Clone and run
git clone https://github.com/LukasLow/echter.link.git
cd echter.link
docker compose up -d

# Or pull from GitHub Container Registry
docker run -d -p 8080:8080 ghcr.io/lukaslow/echter.link:latest
```

### Local Development

```bash
# Clone repository
git clone https://github.com/LukasLow/echter.link.git
cd echter.link

# Install dependencies
go mod tidy

# Run server
go run cmd/server.go
```

## 🐳 Docker Registry

The Docker image is automatically built and pushed to GitHub Container Registry:

```bash
# Pull the latest image
docker pull ghcr.io/lukaslow/echter.link:latest

# Run with volume for persistence
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/root/data \
  ghcr.io/lukaslow/echter.link:latest
```

## 📁 Project Structure

```
echter.link/
├── cmd/
│   └── server.go              # Main application entry point
├── internal/
│   ├── database/
│   │   └── database.go        # Database initialization and configuration
│   ├── handlers/
│   │   ├── url.go            # URL shortening API handlers
│   │   └── web.go            # Web UI handler with dark mode
│   └── models/
│       └── shorturl.go        # Data models
├── .github/workflows/
│   └── docker.yml            # GitHub Actions for CI/CD
├── Dockerfile                 # Multi-stage Docker build
├── docker-compose.yml         # Local development setup
├── go.mod                    # Go module definition
└── README.md                 # This file
```

## 🔧 Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `GIN_MODE` | `release` | Gin framework mode |
| `DB_PATH` | `./echter.link.sqlite` | SQLite database file path |
| `DOMAIN` | `http://localhost:8080` | Base domain for short URLs |

## 🎨 UI Features

### Dark Mode Design
- **Background**: Deep black (`#0a0a0a`)
- **Containers**: Dark gray (`#1a1a1a`)
- **Inputs**: Medium gray (`#2a2a2a`)
- **Accents**: Blue gradients with hover effects

### Interactive Effects
- **Mouse Glow**: Smooth light source following cursor movement
- **Snow Animation**: Continuous falling snowflakes with random properties
- **Snow Melt**: Interactive melting when mouse touches snowflakes
- **Sparkle Effects**: Visual feedback during snow melt

### Footer
- **Version Display**: Current version badge
- **Navigation**: Statistics, About, GitHub links
- **Credits**: Build information and copyright

## 📊 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Web interface |
| `POST` | `/api/shorten` | Create short URL |
| `GET` | `/:code` | Redirect to original URL |
| `GET` | `/#:code` | Hash-based redirect |

### Create Short URL

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "original_url": "https://example.com",
    "custom_code": "my-link",
    "expires_in": 24
  }'
```

Response:
```json
{
  "short_url": "http://localhost:8080/#my-link",
  "short_code": "my-link",
  "expires_at": null
}
```

## 🗄️ Database Schema

SQLite database with one table:

```sql
CREATE TABLE short_urls (
    id TEXT PRIMARY KEY,
    short_code TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    clicks INTEGER DEFAULT 0,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🔒 Security

- **URL Validation**: Only allows `http://` and `https://` URLs
- **Input Sanitization**: Proper validation and sanitization
- **SQLite Security**: File-based database with proper permissions
- **Docker Security**: Multi-stage builds, minimal base image

## 📈 Performance

- **Memory Usage**: ~250MB base, 500MB limit
- **CPU Usage**: Optimized for 0.5 vCPU
- **Response Time**: <100ms for API calls
- **Container Size**: ~20MB compressed

## 🔄 CI/CD

GitHub Actions automatically:
- Builds Docker images on push to main
- Pushes to GitHub Container Registry
- Supports multi-architecture (amd64/arm64)
- Generates SBOM for security scanning

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is open source. See the LICENSE file for details.

## 🙏 Acknowledgments

- [Gin Framework](https://gin-gonic.com/) - HTTP web framework
- [Modern SQLite](https://modernc.org/sqlite) - SQLite driver
- [Docker](https://www.docker.com/) - Container platform
- [GitHub Actions](https://github.com/features/actions) - CI/CD

---

**Built with ❤️ using Go, Gin, & SQLite**
