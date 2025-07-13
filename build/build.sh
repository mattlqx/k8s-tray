#!/bin/bash

# Build script for k8s-tray
# This script helps with building the application for different platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="k8s-tray"
BUILD_DIR="dist"
CMD_DIR="cmd"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}"

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Using Go version: $GO_VERSION"
}

# Create build directory
create_build_dir() {
    mkdir -p "$BUILD_DIR"
    print_info "Created build directory: $BUILD_DIR"
}

# Build for current platform
build_current() {
    print_info "Building for current platform..."
    go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME" "$CMD_DIR/main.go"
    print_success "Built: $BUILD_DIR/$BINARY_NAME"
}

# Build for macOS
build_darwin() {
    print_info "Building for macOS..."

    # AMD64
    print_info "Building for macOS AMD64..."
    GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME-darwin-amd64" "$CMD_DIR/main.go"

    # ARM64
    print_info "Building for macOS ARM64..."
    GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME-darwin-arm64" "$CMD_DIR/main.go"

    # Create universal binary if lipo is available
    if command -v lipo >/dev/null 2>&1; then
        print_info "Creating universal macOS binary..."
        lipo -create -output "$BUILD_DIR/$BINARY_NAME-darwin-universal" \
            "$BUILD_DIR/$BINARY_NAME-darwin-amd64" \
            "$BUILD_DIR/$BINARY_NAME-darwin-arm64"
        print_success "Built universal macOS binary: $BUILD_DIR/$BINARY_NAME-darwin-universal"
    elif command -v x86_64-apple-darwin24-lipo >/dev/null 2>&1; then
        print_info "Creating universal macOS binary with osxcross lipo..."
        x86_64-apple-darwin24-lipo -create -output "$BUILD_DIR/$BINARY_NAME-darwin-universal" \
            "$BUILD_DIR/$BINARY_NAME-darwin-amd64" \
            "$BUILD_DIR/$BINARY_NAME-darwin-arm64"
        print_success "Built universal macOS binary: $BUILD_DIR/$BINARY_NAME-darwin-universal"
    elif command -v x86_64-apple-darwin22-lipo >/dev/null 2>&1; then
        print_info "Creating universal macOS binary with osxcross lipo..."
        x86_64-apple-darwin22-lipo -create -output "$BUILD_DIR/$BINARY_NAME-darwin-universal" \
            "$BUILD_DIR/$BINARY_NAME-darwin-amd64" \
            "$BUILD_DIR/$BINARY_NAME-darwin-arm64"
        print_success "Built universal macOS binary: $BUILD_DIR/$BINARY_NAME-darwin-universal"
    else
        print_warning "lipo not available - universal binary not created"
        print_success "Built individual macOS binaries (AMD64 and ARM64)"
    fi
}

# Build for Linux
build_linux() {
    print_info "Building for Linux..."

    # AMD64
    print_info "Building for Linux AMD64..."
    GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME-linux-amd64" "$CMD_DIR/main.go"

    # ARM64
    print_info "Building for Linux ARM64..."
    GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME-linux-arm64" "$CMD_DIR/main.go"

    print_success "Built Linux binaries"
}

# Build for Windows
build_windows() {
    print_info "Building for Windows..."

    # AMD64
    print_info "Building for Windows AMD64..."
    GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME-windows-amd64.exe" "$CMD_DIR/main.go"

    print_success "Built Windows binaries"
}

# Build all platforms
build_all() {
    print_info "Building for all platforms..."
    build_current
    build_darwin
    build_linux
    build_windows
}

# Clean build directory
clean() {
    print_info "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
    print_success "Cleaned build directory"
}

# Show build information
show_info() {
    echo -e "${BLUE}Build Information:${NC}"
    echo "  Binary Name: $BINARY_NAME"
    echo "  Version: $VERSION"
    echo "  Commit: $COMMIT"
    echo "  Build Time: $BUILD_TIME"
    echo "  Build Directory: $BUILD_DIR"
}

# List built binaries
list_binaries() {
    if [ -d "$BUILD_DIR" ]; then
        print_info "Built binaries:"
        ls -la "$BUILD_DIR"
    else
        print_warning "No build directory found"
    fi
}

# Show help
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  current     Build for current platform"
    echo "  darwin      Build for macOS"
    echo "  linux       Build for Linux"
    echo "  windows     Build for Windows"
    echo "  all         Build for all platforms"
    echo "  clean       Clean build directory"
    echo "  info        Show build information"
    echo "  list        List built binaries"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 current     # Build for current platform"
    echo "  $0 darwin      # Build for macOS"
    echo "  $0 all         # Build for all platforms"
}

# Main script logic
main() {
    case "${1:-current}" in
        current)
            check_go
            create_build_dir
            build_current
            ;;
        darwin)
            check_go
            create_build_dir
            build_darwin
            ;;
        linux)
            check_go
            create_build_dir
            build_linux
            ;;
        windows)
            check_go
            create_build_dir
            build_windows
            ;;
        all)
            check_go
            create_build_dir
            build_all
            ;;
        clean)
            clean
            ;;
        info)
            show_info
            ;;
        list)
            list_binaries
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
