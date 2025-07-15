#!/bin/bash

# macOS App Bundle Creation Script
# This script creates a proper macOS .app bundle from the binary

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
BINARY_NAME="k8s-tray"
BUILD_DIR="dist"
APP_NAME="K8s Tray"
BUNDLE_ID="net.lqx.k8s-tray"

# Function to create app bundle
create_app_bundle() {
    local binary_path="$1"
    local app_bundle_path="$2"

    if [ ! -f "$binary_path" ]; then
        print_error "Binary not found: $binary_path"
        return 1
    fi

    print_info "Creating app bundle: $app_bundle_path"

    # Remove existing bundle if it exists
    if [ -d "$app_bundle_path" ]; then
        rm -rf "$app_bundle_path"
    fi

    # Create bundle structure
    mkdir -p "$app_bundle_path/Contents/MacOS"
    mkdir -p "$app_bundle_path/Contents/Resources"

    # Copy binary
    cp "$binary_path" "$app_bundle_path/Contents/MacOS/$BINARY_NAME"
    chmod +x "$app_bundle_path/Contents/MacOS/$BINARY_NAME"

    # Copy Info.plist
    if [ -f "assets/Info.plist" ]; then
        cp "assets/Info.plist" "$app_bundle_path/Contents/Info.plist"
    else
        print_warning "Info.plist not found, creating basic one"
        create_info_plist "$app_bundle_path/Contents/Info.plist"
    fi

    # Copy icons if they exist
    if [ -d "assets/icons" ]; then
        cp -r assets/icons/* "$app_bundle_path/Contents/Resources/"
    fi

    print_success "App bundle created: $app_bundle_path"
}

# Function to create basic Info.plist if it doesn't exist
create_info_plist() {
    local plist_path="$1"

    cat > "$plist_path" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$BINARY_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>$BUNDLE_ID</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSAppTransportSecurity</key>
    <dict>
        <key>NSAllowsArbitraryLoads</key>
        <true/>
    </dict>
    <key>NSHumanReadableCopyright</key>
    <string>© 2025 K8s Tray. All rights reserved.</string>
</dict>
</plist>
EOF
}

# Function to sign the app bundle (optional)
sign_app_bundle() {
    local app_bundle_path="$1"
    local identity="$2"

    if [ -z "$identity" ]; then
        print_info "No signing identity provided, skipping code signing"
        return 0
    fi

    print_info "Signing app bundle with identity: $identity"

    # Sign the binary first
    codesign --force --sign "$identity" "$app_bundle_path/Contents/MacOS/$BINARY_NAME"

    # Sign the app bundle
    codesign --force --sign "$identity" "$app_bundle_path"

    print_success "App bundle signed successfully"
}

# Main execution
main() {
    local binary_arch="$1"
    local signing_identity="$2"

    if [ -z "$binary_arch" ]; then
        print_error "Usage: $0 <binary-architecture> [signing-identity]"
        print_error "Example: $0 darwin-amd64"
        print_error "Example: $0 darwin-universal \"Developer ID Application: Your Name\""
        exit 1
    fi

    local binary_path="$BUILD_DIR/$BINARY_NAME-$binary_arch"
    local app_bundle_path="$BUILD_DIR/$APP_NAME.app"

    # Create app bundle
    create_app_bundle "$binary_path" "$app_bundle_path"

    # Sign if identity provided
    if [ -n "$signing_identity" ]; then
        sign_app_bundle "$app_bundle_path" "$signing_identity"
    fi

    print_success "macOS app bundle ready: $app_bundle_path"
    print_info "To install: drag $app_bundle_path to /Applications"
}

# Run main function with all arguments
main "$@"
