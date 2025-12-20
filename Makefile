.PHONY: build test install clean help

# Binary name
BINARY_NAME=epgtimer
OUTPUT_DIR=bin

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/epgtimer

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/epgtimer
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/epgtimer
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/epgtimer
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/epgtimer

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	go install ./cmd/epgtimer

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(OUTPUT_DIR)
	rm -f coverage.txt coverage.html
	go clean

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  install       - Install the binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help message"
