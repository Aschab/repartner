# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/configs /app/configs

EXPOSE 8080

CMD ["/app/api"]
