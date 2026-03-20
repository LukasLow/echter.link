# Stage 1: Builder
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies (cacheable)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statisches Binary ohne CGO → kleiner + schneller
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o echte-link ./cmd/server.go

# Stage 2: Minimal Runtime
FROM scratch

# Copy CA certs for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/echte-link /echte-link

# Create data directory
WORKDIR /root
RUN mkdir -p /root/data

# Expose port
EXPOSE 8080

# Run the binary
CMD ["/echte-link"]
