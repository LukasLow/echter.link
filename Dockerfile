# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statisches Binary ohne CGO (modernc.org/sqlite ist pure Go)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o echte-link ./cmd/server.go

# Stage 2: Minimal Runtime
FROM alpine:latest

# CA-Certs für HTTPS
RUN apk --no-cache add ca-certificates

# Binary kopieren
COPY --from=builder /app/echte-link /echte-link

# Datenverzeichnis erstellen
RUN mkdir -p /root/data

WORKDIR /root

EXPOSE 8080

CMD ["/echte-link"]
