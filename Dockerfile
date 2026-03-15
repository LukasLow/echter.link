# Use Alpine Linux for minimal size
FROM golang:1.21-alpine AS builder

# Install gcc and musl-dev for CGO
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o echte-link ./cmd/server.go

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/echte-link .

# Create data directory
RUN mkdir -p /root/data

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./echte-link"]
