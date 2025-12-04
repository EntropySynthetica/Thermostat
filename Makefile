.PHONY: all cli webserver clean run-web help

# Default target
all: cli webserver

# Build CLI application
cli:
	@echo "Building CLI application..."
	@go build -o bin/thermostat thermostat.go
	@echo "✓ CLI built: bin/thermostat"

# Build web server
webserver:
	@echo "Building web server..."
	@go build -o bin/webserver ./cmd/webserver
	@echo "✓ Web server built: bin/webserver"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bin/thermostat bin/webserver
	@echo "✓ Clean complete"

# Run web server (default port 8080)
run-web: webserver
	@./bin/webserver

# Show help
help:
	@echo "Available targets:"
	@echo "  make all       - Build both CLI and web server (default)"
	@echo "  make cli       - Build only the CLI application"
	@echo "  make webserver - Build only the web server"
	@echo "  make run-web   - Build and run the web server"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make help      - Show this help message"
