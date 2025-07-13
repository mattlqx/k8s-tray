#!/bin/bash

# Cross-compilation script for k8s-tray
# This script sets up the environment for cross-compilation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_tools() {
    print_status "Checking cross-compilation tools..."

    # Check for ARM64 cross-compiler
    if ! command -v aarch64-linux-gnu-gcc &> /dev/null; then
        print_error "ARM64 cross-compiler not found. Install with: sudo apt install gcc-aarch64-linux-gnu"
        exit 1
    fi

    # Check for MinGW cross-compiler
    if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
        print_error "Windows cross-compiler not found. Install with: sudo apt install gcc-mingw-w64"
        exit 1
    fi

    print_status "All required tools are available"
}

# Build for Linux ARM64
build_linux_arm64() {
    print_status "Building for Linux ARM64..."

    export CGO_ENABLED=1
    export GOOS=linux
    export GOARCH=arm64
    export CC=aarch64-linux-gnu-gcc
    export CXX=aarch64-linux-gnu-g++
    export AR=aarch64-linux-gnu-ar
    export STRIP=aarch64-linux-gnu-strip
    export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig

    if go build -ldflags="-w -s" -o dist/k8s-tray-linux-arm64 cmd/main.go; then
        print_status "Linux ARM64 build successful"
        return 0
    else
        print_error "Linux ARM64 build failed"
        return 1
    fi
}

# Build for Windows AMD64
build_windows_amd64() {
    print_status "Building for Windows AMD64..."

    export CGO_ENABLED=1
    export GOOS=windows
    export GOARCH=amd64
    export CC=x86_64-w64-mingw32-gcc
    export CXX=x86_64-w64-mingw32-g++
    export AR=x86_64-w64-mingw32-ar
    export STRIP=x86_64-w64-mingw32-strip

    if go build -ldflags="-w -s" -o dist/k8s-tray-windows-amd64.exe cmd/main.go; then
        print_status "Windows AMD64 build successful"
        return 0
    else
        print_error "Windows AMD64 build failed"
        return 1
    fi
}

# Build for Windows ARM64 (requires special setup)
build_windows_arm64() {
    print_status "Building for Windows ARM64..."
    print_warning "Windows ARM64 cross-compilation may require additional setup"

    export CGO_ENABLED=1
    export GOOS=windows
    export GOARCH=arm64
    # Note: ARM64 Windows cross-compilation is complex and may not work out of the box

    if go build -ldflags="-w -s" -o dist/k8s-tray-windows-arm64.exe cmd/main.go; then
        print_status "Windows ARM64 build successful"
        return 0
    else
        print_warning "Windows ARM64 build failed (this is expected without specialized toolchain)"
        return 1
    fi
}

# Build for macOS (requires osxcross or similar)
build_darwin() {
    print_status "Building for macOS..."

    # Check for osxcross
    if command -v x86_64-apple-darwin22-clang &> /dev/null; then
        print_status "Using osxcross for macOS cross-compilation"

        # macOS AMD64
        export CGO_ENABLED=1
        export GOOS=darwin
        export GOARCH=amd64
        export CC=x86_64-apple-darwin22-clang
        export CXX=x86_64-apple-darwin22-clang++

        if go build -ldflags="-w -s" -o dist/k8s-tray-darwin-amd64 cmd/main.go; then
            print_status "macOS AMD64 build successful"
        else
            print_error "macOS AMD64 build failed"
            return 1
        fi

        # macOS ARM64
        export GOARCH=arm64
        export CC=arm64-apple-darwin22-clang
        export CXX=arm64-apple-darwin22-clang++

        if go build -ldflags="-w -s" -o dist/k8s-tray-darwin-arm64 cmd/main.go; then
            print_status "macOS ARM64 build successful"
        else
            print_error "macOS ARM64 build failed"
            return 1
        fi

        return 0

    elif command -v o64-clang &> /dev/null; then
        print_status "Using osxcross o64-clang for macOS cross-compilation"

        # macOS AMD64
        export CGO_ENABLED=1
        export GOOS=darwin
        export GOARCH=amd64
        export CC=o64-clang
        export CXX=o64-clang++

        if go build -ldflags="-w -s" -o dist/k8s-tray-darwin-amd64 cmd/main.go; then
            print_status "macOS AMD64 build successful"
        else
            print_error "macOS AMD64 build failed"
            return 1
        fi

        # macOS ARM64
        export GOARCH=arm64
        export CC=oa64-clang
        export CXX=oa64-clang++

        if go build -ldflags="-w -s" -o dist/k8s-tray-darwin-arm64 cmd/main.go; then
            print_status "macOS ARM64 build successful"
        else
            print_error "macOS ARM64 build failed"
            return 1
        fi

        return 0

    else
        print_warning "osxcross not found - attempting build without cross-compilation"
        print_warning "This will likely fail without proper macOS SDK setup"

        # This would require osxcross setup, which is complex
        # For now, we'll just attempt the build and let it fail gracefully

        export CGO_ENABLED=1
        export GOOS=darwin
        export GOARCH=amd64

        if go build -ldflags="-w -s" -o dist/k8s-tray-darwin-amd64 cmd/main.go; then
            print_status "macOS AMD64 build successful"
        else
            print_warning "macOS AMD64 build failed (expected without osxcross)"
        fi

        export GOARCH=arm64

        if go build -ldflags="-w -s" -o dist/k8s-tray-darwin-arm64 cmd/main.go; then
            print_status "macOS ARM64 build successful"
        else
            print_warning "macOS ARM64 build failed (expected without osxcross)"
        fi

        return 1
    fi

    # Create universal binary if both architectures built successfully
    if [ -f "dist/k8s-tray-darwin-amd64" ] && [ -f "dist/k8s-tray-darwin-arm64" ]; then
        print_status "Creating universal macOS binary..."

        if command -v lipo &> /dev/null; then
            lipo -create -output dist/k8s-tray-darwin-universal dist/k8s-tray-darwin-amd64 dist/k8s-tray-darwin-arm64
            print_status "Universal binary created: dist/k8s-tray-darwin-universal"
        elif command -v x86_64-apple-darwin24-lipo &> /dev/null; then
            x86_64-apple-darwin24-lipo -create -output dist/k8s-tray-darwin-universal dist/k8s-tray-darwin-amd64 dist/k8s-tray-darwin-arm64
            print_status "Universal binary created with osxcross lipo: dist/k8s-tray-darwin-universal"
        elif command -v x86_64-apple-darwin22-lipo &> /dev/null; then
            x86_64-apple-darwin22-lipo -create -output dist/k8s-tray-darwin-universal dist/k8s-tray-darwin-amd64 dist/k8s-tray-darwin-arm64
            print_status "Universal binary created with osxcross lipo: dist/k8s-tray-darwin-universal"
        else
            print_warning "lipo not available - universal binary not created"
            print_status "Individual binaries available: dist/k8s-tray-darwin-amd64, dist/k8s-tray-darwin-arm64"
        fi
    else
        print_warning "Cannot create universal binary - one or both architecture builds failed"
    fi
}

# Main function
main() {
    print_status "Starting cross-compilation setup..."

    # Create dist directory
    mkdir -p dist

    # Check tools
    check_tools

    # Track success/failure
    success_count=0
    total_count=0

    # Build for each platform
    platforms=("linux_arm64" "windows_amd64" "windows_arm64" "darwin")

    for platform in "${platforms[@]}"; do
        total_count=$((total_count + 1))
        echo
        print_status "Building for $platform..."

        case $platform in
            "linux_arm64")
                if build_linux_arm64; then
                    success_count=$((success_count + 1))
                fi
                ;;
            "windows_amd64")
                if build_windows_amd64; then
                    success_count=$((success_count + 1))
                fi
                ;;
            "windows_arm64")
                if build_windows_arm64; then
                    success_count=$((success_count + 1))
                fi
                ;;
            "darwin")
                if build_darwin; then
                    success_count=$((success_count + 1))
                fi
                ;;
        esac
    done

    echo
    print_status "Cross-compilation completed: $success_count/$total_count builds successful"

    # List built binaries
    echo
    print_status "Built binaries:"
    ls -la dist/ | grep -E "(k8s-tray|\.exe)" || echo "No binaries found"
}

# Run main function
main "$@"
