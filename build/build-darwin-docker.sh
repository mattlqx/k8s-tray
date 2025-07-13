#!/bin/bash

# Alternative Darwin build script using Docker
# This approach uses a Docker container with macOS cross-compilation tools

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker to use this build method."
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

print_info "Building Darwin binaries using Docker..."

# Create a temporary Dockerfile
cat > /tmp/Dockerfile.darwin << 'EOF'
FROM golang:1.21-alpine

# Install cross-compilation tools
RUN apk add --no-cache \
    git \
    build-base \
    clang \
    musl-dev

# Set up Go for cross-compilation
ENV CGO_ENABLED=0
ENV GOOS=darwin

WORKDIR /app

# Copy source code
COPY . .

# Build AMD64
RUN GOARCH=amd64 go build -ldflags="-w -s" -o dist/k8s-tray-darwin-amd64 cmd/main.go

# Build ARM64
RUN GOARCH=arm64 go build -ldflags="-w -s" -o dist/k8s-tray-darwin-arm64 cmd/main.go

CMD ["ls", "-la", "dist/"]
EOF

# Build the Docker image
print_info "Building Docker image for Darwin cross-compilation..."
docker build -f /tmp/Dockerfile.darwin -t k8s-tray-darwin-builder .

# Run the container to build binaries
print_info "Running Docker container to build Darwin binaries..."
docker run --rm -v "$(pwd)/dist:/app/dist" k8s-tray-darwin-builder

# Clean up
rm -f /tmp/Dockerfile.darwin

print_success "Darwin binaries built successfully using Docker!"
print_info "Note: These binaries are built without CGO, so systray functionality may be limited."
print_info "For full functionality, consider setting up osxcross instead."
