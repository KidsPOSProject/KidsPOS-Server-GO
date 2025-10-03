# KidsPOS Go Server Makefile

.PHONY: help build run clean test deps dev build-pi deploy-pi docker-build docker-run

# Default target
help:
	@echo "Available commands:"
	@echo "  make deps        - Download dependencies"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make dev         - Run in development mode with hot reload"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make build-pi    - Build for Raspberry Pi (ARM)"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Run Docker container"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	go build -ldflags="-s -w" -o bin/kidspos cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Development mode with hot reload (requires air)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf dist/
	rm -f kidspos.db

# Build for Raspberry Pi (ARM architectures)
# Using CGO_ENABLED=0 for static binaries (Pure Go with modernc.org/sqlite)
build-pi:
	@echo "Building for Raspberry Pi 4/5 (ARM64)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/kidspos-arm64 cmd/server/main.go
	@echo "Building for Raspberry Pi 3/Zero 2W (ARMv7)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o dist/kidspos-armv7 cmd/server/main.go
	@echo "Building for Raspberry Pi Zero W (ARMv6)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o dist/kidspos-armv6 cmd/server/main.go
	@echo "Builds complete! Files in dist/"
	@echo ""
	@echo "Binary sizes:"
	@ls -lh dist/kidspos-arm*

# Deploy to Raspberry Pi (requires SSH access)
deploy-pi:
	@read -p "Enter Raspberry Pi IP address: " PI_IP && \
	scp dist/kidspos-arm64 pi@$$PI_IP:/home/pi/kidspos/ && \
	ssh pi@$$PI_IP "chmod +x /home/pi/kidspos/kidspos-arm64 && sudo systemctl restart kidspos"

# Build Docker image
docker-build:
	docker build -t kidspos-go:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 -v $(PWD)/kidspos.db:/app/kidspos.db kidspos-go:latest

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/kidspos-linux-amd64 cmd/server/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/kidspos-linux-arm64 cmd/server/main.go
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/kidspos-darwin-amd64 cmd/server/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/kidspos-darwin-arm64 cmd/server/main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/kidspos-windows-amd64.exe cmd/server/main.go
	@echo "All builds complete!"

# Initialize database
init-db:
	@echo "Initializing database..."
	@rm -f kidspos.db
	@go run cmd/server/main.go &
	@sleep 2
	@pkill -f "cmd/server/main.go"
	@echo "Database initialized!"

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Generate code coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"