#!/bin/bash

# Setup script for k8s-tray development environment
# This script sets up the development environment for k8s-tray

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

# Check if running on macOS
check_macos() {
    if [[ "$OSTYPE" != "darwin"* ]]; then
        print_warning "This application is designed for macOS. You can still develop on other platforms."
    fi
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        print_info "Visit: https://golang.org/dl/"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go version: $GO_VERSION"
}

# Check if Python is installed (for pre-commit)
check_python() {
    if ! command -v python3 &> /dev/null; then
        print_warning "Python3 is not installed. Some development tools may not work."
        print_info "Visit: https://www.python.org/downloads/"
    else
        PYTHON_VERSION=$(python3 --version)
        print_success "Python version: $PYTHON_VERSION"
    fi
}

# Install Go dependencies
install_go_deps() {
    print_info "Installing Go dependencies..."
    go mod download
    go mod tidy
    print_success "Go dependencies installed"
}

# Install development tools
install_dev_tools() {
    print_info "Installing development tools..."

    # Install golangci-lint
    if ! command -v golangci-lint &> /dev/null; then
        print_info "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        print_success "golangci-lint installed"
    else
        print_info "golangci-lint already installed"
    fi

    # Install pre-commit
    if command -v python3 &> /dev/null; then
        if ! command -v pre-commit &> /dev/null; then
            print_info "Installing pre-commit..."
            pip3 install pre-commit --break-system-packages
            print_success "pre-commit installed"
        else
            print_info "pre-commit already installed"
        fi
    fi
}

# Setup pre-commit hooks
setup_pre_commit() {
    if command -v pre-commit &> /dev/null; then
        print_info "Setting up pre-commit hooks..."
        pre-commit install
        pre-commit install --hook-type commit-msg
        print_success "pre-commit hooks installed"
    else
        print_warning "pre-commit not available, skipping hook setup"
    fi
}

# Create example configuration
create_example_config() {
    CONFIG_FILE="$HOME/.k8s-tray.yaml"
    if [ ! -f "$CONFIG_FILE" ]; then
        print_info "Creating example configuration..."
        cp config.example.yaml "$CONFIG_FILE"
        print_success "Example configuration created at $CONFIG_FILE"
        print_info "Edit $CONFIG_FILE to customize your settings"
    else
        print_info "Configuration file already exists at $CONFIG_FILE"
    fi
}

# Run tests
run_tests() {
    print_info "Running tests..."
    go test ./...
    print_success "Tests passed"
}

# Build the application
build_app() {
    print_info "Building application..."
    make build
    print_success "Application built successfully"
}

# Show next steps
show_next_steps() {
    echo ""
    echo -e "${GREEN}🎉 Setup completed successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Edit ~/.k8s-tray.yaml to configure your settings"
    echo "  2. Make sure you have a valid kubeconfig file"
    echo "  3. Run the application: ./dist/k8s-tray"
    echo ""
    echo "Development commands:"
    echo "  make test          # Run tests"
    echo "  make lint          # Run linter"
    echo "  make format        # Format code"
    echo "  make build         # Build application"
    echo "  make pre-commit-run # Run pre-commit hooks"
    echo ""
    echo "For more information, see the README.md file"
}

# Main setup function
main() {
    echo -e "${BLUE}Setting up k8s-tray development environment...${NC}"
    echo ""

    check_macos
    check_go
    check_python
    install_go_deps
    install_dev_tools
    setup_pre_commit
    create_example_config
    run_tests
    build_app
    show_next_steps
}

# Run main function
main "$@"
