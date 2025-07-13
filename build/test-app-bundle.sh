#!/bin/bash

# Test script to verify the macOS app bundle works correctly

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

echo "Testing macOS app bundle..."

# Check if app bundle exists
if [ ! -d "dist/K8s Tray.app" ]; then
    print_error "App bundle not found. Run 'make build-darwin-app' first."
    exit 1
fi

print_success "App bundle found"

# Check app bundle structure
if [ ! -f "dist/K8s Tray.app/Contents/Info.plist" ]; then
    print_error "Info.plist not found in app bundle"
    exit 1
fi

print_success "Info.plist found"

# Check if binary exists and is executable
if [ ! -f "dist/K8s Tray.app/Contents/MacOS/k8s-tray" ]; then
    print_error "Binary not found in app bundle"
    exit 1
fi

if [ ! -x "dist/K8s Tray.app/Contents/MacOS/k8s-tray" ]; then
    print_error "Binary is not executable"
    exit 1
fi

print_success "Binary is present and executable"

# Check LSUIElement in Info.plist
if ! grep -q "<key>LSUIElement</key>" "dist/K8s Tray.app/Contents/Info.plist"; then
    print_warning "LSUIElement not found in Info.plist"
else
    print_success "LSUIElement found in Info.plist"
fi

# Check bundle identifier
if ! grep -q "com.k8s-tray.k8s-tray" "dist/K8s Tray.app/Contents/Info.plist"; then
    print_warning "Bundle identifier not found in Info.plist"
else
    print_success "Bundle identifier found in Info.plist"
fi

# Check if running on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    print_success "Running on macOS - app bundle should work"

    # Check if we can validate the app bundle
    if command -v codesign &> /dev/null; then
        echo "Checking code signature..."
        if codesign --verify --deep --strict "dist/K8s Tray.app" 2>/dev/null; then
            print_success "App bundle is properly signed"
        else
            print_warning "App bundle is not signed (this is okay for development)"
        fi
    fi

    echo ""
    echo "To test the app bundle:"
    echo "1. Run: open 'dist/K8s Tray.app'"
    echo "2. Check the menu bar for the k8s-tray icon"
    echo "3. If it doesn't appear, check System Preferences > Security & Privacy"
    echo ""
else
    print_warning "Not running on macOS - cannot test app bundle functionality"
    echo "The app bundle will need to be tested on a macOS system"
    exit 1
fi

print_success "App bundle validation complete"
