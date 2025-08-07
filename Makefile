# Makefile for SecretShare

# Build the application
build:
	go build -o secret_share cmd/secret_share/main.go

# Run tests
test:
	go test -v ./core

# Run tests with coverage
test-coverage:
	go test -cover ./core

# Install the application
install:
	go install cmd/secret_share/main.go

# Clean build artifacts
clean:
	rm -f secret_share

# Help
help:
	@echo "SecretShare Makefile"
	@echo "==================="
	@echo "build        - Build the application"
	@echo "test         - Run tests"
	@echo "test-coverage - Run tests with coverage"
	@echo "install      - Install the application"
	@echo "clean        - Clean build artifacts"
	@echo "help         - Show this help message"

.PHONY: build test test-coverage install clean help
