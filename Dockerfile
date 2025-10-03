# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o kidspos cmd/server/main.go

# Runtime stage - Ultra lightweight
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite

# Create non-root user
RUN addgroup -g 1000 kidspos && \
    adduser -D -u 1000 -G kidspos kidspos

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/kidspos .
COPY --from=builder /build/web ./web

# Create data directory
RUN mkdir -p /app/data && chown -R kidspos:kidspos /app

# Switch to non-root user
USER kidspos

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080 || exit 1

# Set environment variables
ENV DATABASE_PATH=/app/data/kidspos.db
ENV PORT=8080

# Run the application
CMD ["./kidspos"]
