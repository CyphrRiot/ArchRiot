# ArchRiot Makefile
# Build system for the ArchRiot installer

.PHONY: build clean install test help

# Variables
BINARY_NAME = archriot
SOURCE_DIR = source
INSTALL_DIR = install
BUILD_DIR = $(SOURCE_DIR)

# Default target
all: build

# Build the installer (with optimizations by default)
build:
	@echo "🔨 Building ArchRiot installer..."
	@cd $(SOURCE_DIR) && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $(BINARY_NAME) .
	@mv $(SOURCE_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Build complete: $(INSTALL_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@rm -f $(SOURCE_DIR)/$(BINARY_NAME)
	@echo "✅ Clean complete"

# Install dependencies
deps:
	@echo "📦 Installing Go dependencies..."
	@cd $(SOURCE_DIR) && go mod tidy
	@echo "✅ Dependencies updated"

# Run tests
test:
	@echo "🧪 Running tests..."
	@cd $(SOURCE_DIR) && go test ./...

# Verify the installer works
verify: build
	@echo "🔍 Verifying installer..."
	@$(INSTALL_DIR)/$(BINARY_NAME) --version
	@echo "✅ Installer verified"

# Development build (faster, no optimizations)
dev:
	@echo "🚀 Building development version..."
	@cd $(SOURCE_DIR) && go build -o $(BINARY_NAME) .
	@mv $(SOURCE_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Development build complete"

# Ultra-optimized build (maximum compression)
ultra:
	@echo "🎯 Building ultra-optimized version..."
	@cd $(SOURCE_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-s -w -extldflags '-static'" \
		-trimpath \
		-o $(BINARY_NAME) .
	@mv $(SOURCE_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "📦 Compressing with UPX (if available)..."
	@which upx > /dev/null 2>&1 && upx --best --lzma $(INSTALL_DIR)/$(BINARY_NAME) || echo "⚠️  UPX not found, skipping compression"
	@echo "✅ Ultra-optimized build complete"

# Release build (optimized)
release:
	@echo "🎯 Building release version..."
	@cd $(SOURCE_DIR) && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $(BINARY_NAME) .
	@mv $(SOURCE_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✅ Release build complete"

# Help target
help:
	@echo "ArchRiot Build System"
	@echo "====================="
	@echo ""
	@echo "Available targets:"
	@echo "  build     - Build the installer (default)"
	@echo "  clean     - Remove build artifacts"
	@echo "  deps      - Install/update Go dependencies"
	@echo "  test      - Run tests"
	@echo "  verify    - Build and verify installer works"
	@echo "  dev       - Fast development build"
	@echo "  release   - Optimized release build"
	@echo "  ultra     - Maximum optimization with UPX compression"
	@echo "  help      - Show this help message"
	@echo ""
	@echo "Example usage:"
	@echo "  make build    # Build the installer"
	@echo "  make clean    # Clean up"
	@echo "  make release  # Build optimized version"
