.PHONY: all build run clean swagger dev help

# Default target
all: swagger build run

# Build the application
build:
	@echo "Building application..."
	@go build -o app .

# Run the application
run: swagger
	@echo "Running application..."
	@go run .

# Run with live-reload during development
dev:
	@echo "Starting development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air is not installed. Running without live reload..."; \
		go run .; \
	fi

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init; \
	else \
		echo "Swag is not installed. Please install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f app
	@rm -rf docs

# Show help
help:
	@echo "Available commands:"
	@echo "  make          - Generate Swagger docs, build and run the application"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application without building"
	@echo "  make swagger  - Generate Swagger documentation"
	@echo "  make dev      - Run with live-reload if Air is installed"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make help     - Show this help message"