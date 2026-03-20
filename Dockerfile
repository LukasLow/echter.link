# Stage 1: Builder
FROM golang:1.21-alpine AS builder

# Install GCC für CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statisches Binary mit CGO für SQLite
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o echte-link ./cmd/server.go

# Stage 2: Minimal Runtime
FROM alpine:latest

# CA-Certs und SQLite Support für CGO
RUN apk --no-cache add ca-certificates sqlite-dev

# Binary kopieren
COPY --from=builder /app/echte-link /echte-link

# Datenverzeichnis erstellen
RUN mkdir -p /root/data

WORKDIR /root

EXPOSE 8080

CMD ["/echte-link"]
