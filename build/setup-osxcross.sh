#!/bin/bash

# Setup script for osxcross to enable macOS cross-compilation on Linux
# This requires a macOS SDK which must be obtained legally

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

# Configuration
OSXCROSS_DIR="$HOME/osxcross"
XCODE_VERSION="16.4"
SDK_VERSION="15.5"

print_info "Setting up osxcross for macOS cross-compilation..."

# Check if osxcross is already installed
if [ -d "$OSXCROSS_DIR" ] && [ -f "$OSXCROSS_DIR/target/bin/x86_64-apple-darwin24-clang" ]; then
    print_success "osxcross is already installed at $OSXCROSS_DIR"
    export PATH="$OSXCROSS_DIR/target/bin:$PATH"
    print_info "osxcross tools are available at: $OSXCROSS_DIR/target/bin"
    exit 0
fi

# Install prerequisites
print_info "Installing prerequisites..."
sudo apt update
sudo apt install -y \
    clang \
    llvm-dev \
    libxml2-dev \
    uuid-dev \
    libssl-dev \
    bash \
    patch \
    make \
    tar \
    xz-utils \
    bzip2 \
    gzip \
    sed \
    cpio \
    libbz2-dev \
    cmake \
    libz-dev \
    liblzma-dev \
    python3 \
    python3-pip

# Clone osxcross
print_info "Cloning osxcross..."
if [ ! -d "$OSXCROSS_DIR" ]; then
    git clone https://github.com/tpoechtrager/osxcross.git "$OSXCROSS_DIR"
fi

cd "$OSXCROSS_DIR"

# Check for SDK
print_warning "IMPORTANT: You need a macOS SDK to continue."
print_warning "The SDK must be obtained legally from Apple."
print_warning "You can get it from:"
print_warning "  1. Xcode.app (if you have access to a Mac)"
print_warning "  2. Command Line Tools for Xcode"
print_warning "  3. Download from Apple Developer portal"

print_info "Looking for SDK files..."
SDK_FOUND=false

# Check common SDK locations with more flexible naming
for sdk_file in \
    "$OSXCROSS_DIR/tarballs/MacOSX${SDK_VERSION}.sdk.tar.xz" \
    "$OSXCROSS_DIR/tarballs/MacOSX${SDK_VERSION}.sdk.tar.bz2" \
    "$OSXCROSS_DIR/tarballs/MacOSX${SDK_VERSION}.sdk.tar.gz" \
    "$HOME/Downloads/MacOSX${SDK_VERSION}.sdk.tar.xz" \
    "$HOME/Downloads/MacOSX${SDK_VERSION}.sdk.tar.bz2" \
    "$HOME/Downloads/MacOSX${SDK_VERSION}.sdk.tar.gz" \
    "/tmp/MacOSX${SDK_VERSION}.sdk.tar.xz" \
    "/tmp/MacOSX${SDK_VERSION}.sdk.tar.bz2" \
    "/tmp/MacOSX${SDK_VERSION}.sdk.tar.gz"; do

    if [ -f "$sdk_file" ]; then
        print_success "Found SDK: $sdk_file"
        mkdir -p "$OSXCROSS_DIR/tarballs"
        cp "$sdk_file" "$OSXCROSS_DIR/tarballs/" || true
        SDK_FOUND=true
        SDK_FILE="$sdk_file"
        break
    fi
done

# Also check for any SDK files with similar patterns
if [ "$SDK_FOUND" = false ]; then
    print_info "Checking for any macOS SDK files..."

    for search_dir in "$HOME/Downloads" "/tmp" "$OSXCROSS_DIR/tarballs"; do
        if [ -d "$search_dir" ]; then
            print_info "Searching in: $search_dir"
            find "$search_dir" -name "*OSX*.sdk*" -o -name "*acOS*.sdk*" 2>/dev/null | while read -r found_sdk; do
                print_info "Found potential SDK: $found_sdk"
                if [ -f "$found_sdk" ]; then
                    print_success "Using SDK: $found_sdk"
                    mkdir -p "$OSXCROSS_DIR/tarballs"
                    cp "$found_sdk" "$OSXCROSS_DIR/tarballs/"
                    SDK_FOUND=true
                    SDK_FILE="$found_sdk"
                    break
                fi
            done
        fi
    done
fi

if [ "$SDK_FOUND" = false ]; then
    print_error "No macOS SDK found!"
    print_info "To continue, you need to:"
    print_info "1. Download the macOS SDK from Apple (legally)"
    print_info "2. Place it in one of these locations:"
    print_info "   - $OSXCROSS_DIR/tarballs/MacOSX${SDK_VERSION}.sdk.tar.xz"
    print_info "   - $HOME/Downloads/MacOSX${SDK_VERSION}.sdk.tar.xz"
    print_info "   - /tmp/MacOSX${SDK_VERSION}.sdk.tar.xz"
    print_info "3. Run this script again"

    print_info "You can also create the SDK from Xcode using:"
    print_info "  ./tools/gen_sdk_package.sh (on macOS with Xcode installed)"

    exit 1
fi

# Build osxcross
print_info "Building osxcross..."
print_info "This may take several minutes..."

# Enable debug mode for troubleshooting
export OCDEBUG=1

# Set a more specific SDK version if we know what we have
if [ -n "$SDK_FILE" ]; then
    print_info "Using SDK file: $SDK_FILE"
    SDK_BASENAME=$(basename "$SDK_FILE")
    print_info "SDK basename: $SDK_BASENAME"

    # Extract version from filename if possible
    if [[ "$SDK_BASENAME" =~ OSX([0-9]+\.[0-9]+) ]]; then
        DETECTED_VERSION="${BASH_REMATCH[1]}"
        print_info "Detected SDK version: $DETECTED_VERSION"
        export OSXCROSS_SDK_VERSION="$DETECTED_VERSION"
    fi
fi

# Try to build osxcross
if ! UNATTENDED=1 ./build.sh; then
    print_error "osxcross build failed!"
    print_info "Checking what went wrong..."

    print_info "Contents of tarballs directory:"
    ls -la "$OSXCROSS_DIR/tarballs/" || echo "No tarballs directory"

    print_info "You can try manual debugging with:"
    print_info "  cd $OSXCROSS_DIR"
    print_info "  OCDEBUG=1 ./build.sh"

    exit 1
fi

# Verify installation
if [ -f "$OSXCROSS_DIR/target/bin/x86_64-apple-darwin24-clang" ] || [ -f "$OSXCROSS_DIR/target/bin/o64-clang" ]; then
    print_success "osxcross installed successfully!"

    # Add to PATH suggestion
    print_info "To use osxcross, add this to your ~/.bashrc or ~/.zshrc:"
    echo "export PATH=\"$OSXCROSS_DIR/target/bin:\$PATH\""

    # Test the installation
    print_info "Testing installation..."
    export PATH="$OSXCROSS_DIR/target/bin:$PATH"

    # Try different compiler names
    COMPILER_FOUND=false
    for compiler in x86_64-apple-darwin24-clang o64-clang x86_64-apple-darwin23-clang x86_64-apple-darwin22-clang; do
        if command -v "$compiler" > /dev/null 2>&1; then
            if $compiler --version > /dev/null 2>&1; then
                print_success "osxcross is working correctly! ($compiler)"
                COMPILER_FOUND=true
                break
            fi
        fi
    done

    if [ "$COMPILER_FOUND" = false ]; then
        print_warning "osxcross installed but may have issues"
        print_info "Available compilers:"
        ls -la "$OSXCROSS_DIR/target/bin/"*clang* 2>/dev/null || echo "No clang compilers found"
    fi

else
    print_error "osxcross build failed!"
    print_info "Expected compiler not found. Checking what was built..."
    if [ -d "$OSXCROSS_DIR/target/bin" ]; then
        print_info "Contents of target/bin:"
        ls -la "$OSXCROSS_DIR/target/bin/"
    fi
    exit 1
fi

print_success "osxcross setup complete!"
print_info "You can now build macOS binaries using the updated Makefile targets."
