#!/bin/bash

# Cross-compilation information script
# Shows available toolchains and their status

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Cross-compilation setup for k8s-tray${NC}"
echo -e "${BLUE}=====================================${NC}"
echo

# Check Go installation
echo -e "${BLUE}Go Environment:${NC}"
echo "Go version: $(go version)"
echo "GOOS: $(go env GOOS)"
echo "GOARCH: $(go env GOARCH)"
echo

# Check cross-compilation toolchains
echo -e "${BLUE}Cross-compilation Toolchains:${NC}"

# Linux ARM64
if command -v aarch64-linux-gnu-gcc &> /dev/null; then
    echo -e "${GREEN}✓ Linux ARM64:${NC} $(aarch64-linux-gnu-gcc --version | head -1)"
else
    echo -e "${RED}✗ Linux ARM64:${NC} Not available (install: sudo apt install gcc-aarch64-linux-gnu)"
fi

# Windows AMD64
if command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo -e "${GREEN}✓ Windows AMD64:${NC} $(x86_64-w64-mingw32-gcc --version | head -1)"
else
    echo -e "${RED}✗ Windows AMD64:${NC} Not available (install: sudo apt install gcc-mingw-w64)"
fi

# Windows ARM64 (note: limited support)
echo -e "${YELLOW}⚠ Windows ARM64:${NC} Limited support (Go can build but may lack proper CGO toolchain)"

# macOS support
if command -v x86_64-apple-darwin22-clang &> /dev/null; then
    echo -e "${GREEN}✓ macOS:${NC} osxcross installed ($(x86_64-apple-darwin22-clang --version | head -1))"
elif command -v o64-clang &> /dev/null; then
    echo -e "${GREEN}✓ macOS:${NC} osxcross installed ($(o64-clang --version | head -1))"
else
    echo -e "${YELLOW}⚠ macOS:${NC} osxcross not found - run ./build/setup-osxcross.sh to install"
fi

echo

# Show build targets
echo -e "${BLUE}Available Build Targets:${NC}"
echo "  make build-linux     - Build for Linux (amd64 + arm64)"
echo "  make build-windows   - Build for Windows (amd64 + arm64)"
echo "  make build-darwin    - Build for macOS (requires osxcross)"
echo "  make build-all       - Build for all platforms"
echo "  make cross-compile   - Advanced cross-compilation script"
echo

# Show current project dependencies
echo -e "${BLUE}Project Dependencies:${NC}"
if [ -f go.mod ]; then
    echo "Dependencies requiring CGO:"
    grep -E "(systray|cgo)" go.mod || echo "  None explicitly marked (but systray likely requires CGO)"
else
    echo "  No go.mod found in current directory"
fi

echo

# Show build status
echo -e "${BLUE}Build Status:${NC}"
if [ -d "dist" ]; then
    echo "Built binaries in dist/:"
    ls -la dist/ | grep -E "(k8s-tray|\.exe)" | sed 's/^/  /'
else
    echo "  No built binaries found (run 'make build-all' to build)"
fi

echo

echo -e "${BLUE}Quick Start:${NC}"
echo "  1. Run 'make build-linux' to build for Linux"
echo "  2. Run 'make build-windows' to build for Windows"
echo "  3. Run 'make cross-compile' for advanced cross-compilation"
echo "  4. Run 'make build-all' to build for all supported platforms"
