# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statisches Binary ohne CGO
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o echte-link ./cmd/server.go

# Stage 2: Minimal Runtime
FROM scratch

# CA-Certs für HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Binary kopieren
COPY --from=builder /app/echte-link /echte-link

WORKDIR /root
RUN mkdir -p /root/data

EXPOSE 8080

CMD ["/echte-link"]
