# Stage 1: Build
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum (if it exists)
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o billing-service ./cmd/api

# Stage 2: Final
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/billing-service .

# Copy migrations
COPY --from=builder /app/internal/infrastructure/database/migrations ./internal/infrastructure/database/migrations
COPY --from=builder /app/docs ./docs

# Expose the application port
EXPOSE 8022

# Run the application
CMD ["./billing-service"]
