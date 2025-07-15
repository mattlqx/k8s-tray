# Makefile for k8s-tray
.PHONY: setup-dev lint test format pre-commit-install build build-darwin build-darwin-universal build-darwin-app build-app-from-binary test-app-bundle deploy-macos test-ssh build-linux build-windows build-native-only build-all build-all-with-app cross-compile cross-info setup-osxcross clean run help

# Variables
BINARY_NAME=k8s-tray
BUILD_DIR=dist
CMD_DIR=cmd
MAIN_FILE=$(CMD_DIR)/main.go
LDFLAGS=-w -s
BUILD_FLAGS=-ldflags="$(LDFLAGS)"

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup-dev: pre-commit-install ## Set up development environment
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/github-action-add-sarif@latest || true

pre-commit-install: ## Install pre-commit hooks
	pip install pre-commit || echo "Please install pre-commit manually"
	pre-commit install
	pre-commit install --hook-type commit-msg

lint: ## Run linters
	golangci-lint run ./...

test: ## Run tests
	go test -race -cover ./...

format: ## Format code
	gofmt -w .
	goimports -w .

pre-commit-run: ## Run pre-commit hooks on all files
	pre-commit run --all-files

build: ## Build the application
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

build-darwin: ## Build for macOS (requires macOS or osxcross)
	@mkdir -p $(BUILD_DIR)
	@echo "Building for macOS AMD64..."
	@if command -v x86_64-apple-darwin24-clang >/dev/null 2>&1; then \
		echo "Using osxcross for macOS cross-compilation"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=x86_64-apple-darwin24-clang CXX=x86_64-apple-darwin24-clang++ go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE); \
	elif command -v x86_64-apple-darwin22-clang >/dev/null 2>&1; then \
		echo "Using osxcross (darwin22) for macOS cross-compilation"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=x86_64-apple-darwin22-clang CXX=x86_64-apple-darwin22-clang++ go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE); \
	elif command -v o64-clang >/dev/null 2>&1; then \
		echo "Using osxcross o64-clang for macOS cross-compilation"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE); \
	else \
		echo "Attempting build without osxcross (may fail)"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE) || echo "Failed to build macOS AMD64 - install osxcross or run on macOS"; \
	fi
	@echo "Building for macOS ARM64..."
	@if command -v arm64-apple-darwin24-clang >/dev/null 2>&1; then \
		echo "Using osxcross for macOS ARM64 cross-compilation"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC=arm64-apple-darwin24-clang CXX=arm64-apple-darwin24-clang++ go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE); \
	elif command -v arm64-apple-darwin22-clang >/dev/null 2>&1; then \
		echo "Using osxcross (darwin22) for macOS ARM64 cross-compilation"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC=arm64-apple-darwin22-clang CXX=arm64-apple-darwin22-clang++ go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE); \
	elif command -v oa64-clang >/dev/null 2>&1; then \
		echo "Using osxcross oa64-clang for macOS ARM64 cross-compilation"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC=oa64-clang CXX=oa64-clang++ go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE); \
	else \
		echo "Attempting build without osxcross (may fail)"; \
		CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE) || echo "Failed to build macOS ARM64 - install osxcross or run on macOS"; \
	fi
	@echo "Creating universal macOS binary..."
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64" ] && [ -f "$(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64" ]; then \
		if command -v lipo >/dev/null 2>&1; then \
			echo "Using lipo to create universal binary"; \
			lipo -create -output $(BUILD_DIR)/$(BINARY_NAME)-darwin-universal $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64; \
			echo "Universal binary created: $(BUILD_DIR)/$(BINARY_NAME)-darwin-universal"; \
		elif command -v x86_64-apple-darwin24-lipo >/dev/null 2>&1; then \
			echo "Using osxcross lipo to create universal binary"; \
			x86_64-apple-darwin24-lipo -create -output $(BUILD_DIR)/$(BINARY_NAME)-darwin-universal $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64; \
			echo "Universal binary created: $(BUILD_DIR)/$(BINARY_NAME)-darwin-universal"; \
		elif command -v x86_64-apple-darwin22-lipo >/dev/null 2>&1; then \
			echo "Using osxcross lipo (darwin22) to create universal binary"; \
			x86_64-apple-darwin22-lipo -create -output $(BUILD_DIR)/$(BINARY_NAME)-darwin-universal $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64; \
			echo "Universal binary created: $(BUILD_DIR)/$(BINARY_NAME)-darwin-universal"; \
		else \
			echo "lipo not available - skipping universal binary creation"; \
			echo "Individual binaries available: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64, $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64"; \
		fi \
	else \
		echo "Cannot create universal binary - one or both architecture builds failed"; \
	fi
	chmod 755 $(BUILD_DIR)/$(BINARY_NAME)*

build-darwin-universal: build-darwin ## Build universal macOS binary (combines AMD64 and ARM64)
	@echo "Universal macOS binary target completed - check dist/ for k8s-tray-darwin-universal"

build-darwin-app: build-darwin ## Build macOS .app bundle
	@echo "Creating macOS app bundle..."
	@./build/create-app-bundle.sh darwin-universal
	@echo "macOS app bundle created: dist/K8s Tray.app"

build-app-from-binary: ## Create macOS .app bundle from existing binary (usage: make build-app-from-binary ARCH=darwin-amd64)
	@if [ -z "$(ARCH)" ]; then \
		echo "Error: ARCH parameter required. Usage: make build-app-from-binary ARCH=darwin-amd64"; \
		exit 1; \
	fi
	@echo "Creating macOS app bundle from existing binary..."
	@./build/create-app-bundle.sh $(ARCH)
	@echo "macOS app bundle created: dist/K8s Tray.app"

test-app-bundle: ## Test the macOS app bundle
	@./build/test-app-bundle.sh

deploy-macos: build-darwin-app ## Deploy and test on macOS system (usage: make deploy-macos HOST=mac-mini-m4.local)
	@./build/deploy-and-test-macos.sh $(HOST) $(REMOTE_PATH)

test-ssh: ## Test SSH connectivity to macOS system (usage: make test-ssh HOST=mac-mini-m4.local)
	@./build/test-ssh-connectivity.sh $(HOST) $(REMOTE_PATH)

build-linux: ## Build for Linux
	@mkdir -p $(BUILD_DIR)
	@echo "Building for Linux AMD64..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@echo "Building for Linux ARM64..."
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ AR=aarch64-linux-gnu-ar go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE) || echo "Failed to build Linux ARM64 - cross-compilation tools required"

build-windows: ## Build for Windows (requires Windows or CGO cross-compilation setup)
	@mkdir -p $(BUILD_DIR)
	@echo "Building for Windows AMD64..."
	@echo "Note: Using MinGW cross-compiler for Windows builds"
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ AR=x86_64-w64-mingw32-ar go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE) || echo "Failed to build Windows AMD64 - MinGW cross-compiler required"
	@echo "Building for Windows ARM64..."
	CGO_ENABLED=1 GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_FILE) || echo "Failed to build Windows ARM64 - specialized toolchain required"

build-native-only: build build-linux ## Build only for current platform and Linux (most reliable)

build-all: build build-darwin build-linux build-windows ## Build for all platforms and architectures (may fail without proper CGO setup)

build-all-with-app: build-all build-darwin-app ## Build for all platforms including macOS app bundle

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)

run: ## Run the application
	go run $(MAIN_FILE)

deps: ## Download dependencies
	go mod download
	go mod tidy

update-deps: ## Update dependencies
	go get -u ./...
	go mod tidy

vendor: ## Vendor dependencies
	go mod vendor

# Development targets
dev-setup: setup-dev ## Alias for setup-dev
dev-run: run ## Alias for run
dev-test: test lint ## Run tests and linting
dev-clean: clean ## Clean and reset development environment

# CI/CD targets
ci-test: test lint ## Run CI tests
ci-build: build-all ## Build for CI

# Release targets
release-build: clean build-all ## Build release binaries
	@echo "Built binaries:"
	@ls -la $(BUILD_DIR)/

cross-compile: ## Use cross-compilation script for better cross-platform builds
	@echo "Running cross-compilation script..."
	@./build/cross-compile.sh

cross-info: ## Show cross-compilation setup and status
	@./build/cross-info.sh

setup-osxcross: ## Setup osxcross for macOS cross-compilation
	@echo "Setting up osxcross for macOS cross-compilation..."
	@./build/setup-osxcross.sh
