# Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Final image
FROM alpine:3.20
WORKDIR /app

# Copy backend binary and config
COPY --from=backend-builder /app/bin/api /app/api
COPY --from=backend-builder /app/configs /app/configs

# Copy frontend build to static folder
COPY --from=frontend-builder /app/dist /app/static

EXPOSE 8080

ENV PORT=8080
ENV PACKS_CONFIG_PATH=/app/configs/packs.json

CMD ["/app/api"]
