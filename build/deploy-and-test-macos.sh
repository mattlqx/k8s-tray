#!/bin/bash

# Script to copy and test the k8s-tray app bundle on macOS system
# Usage: ./deploy-and-test-macos.sh [host] [remote-path]

set -e

# Configuration
HOST=${1:-"mac-mini-m4.local"}
REMOTE_PATH=${2:-"~/src/k8s-tray"}
APP_BUNDLE="dist/K8s Tray.app"
BINARY_NAME="k8s-tray"

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

# Check if app bundle exists
check_app_bundle() {
    if [ ! -d "$APP_BUNDLE" ]; then
        print_error "App bundle not found: $APP_BUNDLE"
        print_info "Run 'make build-darwin-app' first"
        exit 1
    fi
    print_success "App bundle found: $APP_BUNDLE"
}

# Test SSH connection
test_ssh() {
    print_info "Testing SSH connection to $HOST..."
    if ! ssh -o ConnectTimeout=5 "$HOST" "echo 'SSH connection successful'" >/dev/null 2>&1; then
        print_error "Cannot connect to $HOST via SSH"
        print_info "Make sure SSH is enabled on the Mac and you have access"
        exit 1
    fi
    print_success "SSH connection to $HOST successful"
}

# Create remote directory
create_remote_directory() {
    print_info "Creating remote directory: $REMOTE_PATH"
    ssh "$HOST" "mkdir -p $REMOTE_PATH/dist"
    print_success "Remote directory created"
}

# Copy app bundle to macOS system
copy_app_bundle() {
    print_info "Copying app bundle to $HOST:$REMOTE_PATH..."

    # Use rsync for efficient copying
    rsync -avz --progress --delete "$APP_BUNDLE" "$HOST:$REMOTE_PATH/dist/"

    print_success "App bundle copied successfully"
}

# Copy test scripts
copy_test_scripts() {
    print_info "Copying test scripts..."

    # Copy the test script
    rsync -avz build/test-app-bundle.sh "$HOST:$REMOTE_PATH/build/"

    # Make it executable
    ssh "$HOST" "chmod +x $REMOTE_PATH/build/test-app-bundle.sh"

    print_success "Test scripts copied"
}

# Run validation test on remote system
run_validation_test() {
    print_info "Running app bundle validation on $HOST..."

    ssh "$HOST" "cd $REMOTE_PATH && ./build/test-app-bundle.sh"

    print_success "Validation test completed"
}

# Create and copy macOS-specific test script
create_macos_test_script() {
    print_info "Creating macOS-specific test script..."

    # Create a temporary test script
    cat > /tmp/test-macos-tray.sh << 'EOF'
#!/bin/bash

# macOS-specific test script for k8s-tray

set -e

APP_BUNDLE="dist/K8s Tray.app"
BINARY_PATH="$APP_BUNDLE/Contents/MacOS/k8s-tray"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo "=== macOS k8s-tray Test Script ==="
echo "Testing app bundle: $APP_BUNDLE"
echo ""

# Check macOS version
print_info "macOS version: $(sw_vers -productVersion)"
print_info "Architecture: $(uname -m)"
echo ""

# Check if app bundle exists
if [ ! -d "$APP_BUNDLE" ]; then
    print_error "App bundle not found: $APP_BUNDLE"
    exit 1
fi

print_success "App bundle found"

# Check if binary is executable
if [ ! -x "$BINARY_PATH" ]; then
    print_error "Binary is not executable: $BINARY_PATH"
    exit 1
fi

print_success "Binary is executable"

# Check for required frameworks
print_info "Checking for required frameworks..."
if ! otool -L "$BINARY_PATH" | grep -q "CoreFoundation"; then
    print_error "CoreFoundation framework not found"
else
    print_success "CoreFoundation framework found"
fi

# Test running the binary briefly
print_info "Testing binary execution (will timeout after 10 seconds)..."

# Check for timeout command
if command -v /opt/homebrew/bin/timeout >/dev/null 2>&1; then
    /opt/homebrew/bin/timeout 10s "$BINARY_PATH" &
    BINARY_PID=$!
elif command -v timeout >/dev/null 2>&1; then
    timeout 10s "$BINARY_PATH" &
    BINARY_PID=$!
else
    # Fallback: start the process and kill it after 10 seconds
    "$BINARY_PATH" &
    BINARY_PID=$!

    # Sleep for 10 seconds then kill
    (sleep 10; kill $BINARY_PID 2>/dev/null) &
    KILLER_PID=$!
fi

sleep 2

# Check if process is still running
if kill -0 $BINARY_PID 2>/dev/null; then
    print_success "Binary started successfully"

    # Check if it appears in Activity Monitor
    if pgrep -f "k8s-tray" >/dev/null; then
        print_success "Process found in Activity Monitor"
    else
        print_warning "Process not found in Activity Monitor"
    fi

    # Kill the process
    kill $BINARY_PID 2>/dev/null || true
    wait $BINARY_PID 2>/dev/null || true

    # Kill the killer process if it exists
    if [ -n "$KILLER_PID" ]; then
        kill $KILLER_PID 2>/dev/null || true
    fi
else
    print_error "Binary failed to start or crashed"
    exit 1
fi

echo ""
print_info "Manual test steps:"
echo "1. Run: open '$APP_BUNDLE'"
echo "2. Check the menu bar (top-right area) for the k8s-tray icon"
echo "3. If no icon appears, check Console.app for error messages"
echo "4. To kill the app: pkill -f k8s-tray"
echo ""

print_success "Basic tests completed successfully"
print_info "App bundle appears to be working correctly"
EOF

    # Copy to remote system
    scp /tmp/test-macos-tray.sh "$HOST:$REMOTE_PATH/test-macos-tray.sh"
    ssh "$HOST" "chmod +x $REMOTE_PATH/test-macos-tray.sh"

    # Clean up
    rm /tmp/test-macos-tray.sh

    print_success "macOS test script created and copied"
}

# Run the macOS-specific test
run_macos_test() {
    print_info "Running macOS-specific tests..."

    ssh "$HOST" "cd $REMOTE_PATH && ./test-macos-tray.sh"

    print_success "macOS tests completed"
}

# Provide manual testing instructions
show_manual_instructions() {
    echo ""
    echo "======================================"
    echo "Manual Testing Instructions"
    echo "======================================"
    echo ""
    echo "1. SSH to the macOS system:"
    echo "   ssh $HOST"
    echo ""
    echo "2. Navigate to the project directory:"
    echo "   cd $REMOTE_PATH"
    echo ""
    echo "3. Run the app bundle:"
    echo "   open 'dist/K8s Tray.app'"
    echo ""
    echo "4. Check the menu bar for the k8s-tray icon"
    echo ""
    echo "5. If the icon doesn't appear, check logs:"
    echo "   - Open Console.app"
    echo "   - Search for 'k8s-tray' messages"
    echo "   - Or run: log show --predicate 'process == \"k8s-tray\"' --last 5m"
    echo ""
    echo "6. To kill the app:"
    echo "   pkill -f k8s-tray"
    echo ""
    echo "7. To run with console output:"
    echo "   ./dist/K8s\ Tray.app/Contents/MacOS/k8s-tray"
    echo ""
    echo "======================================"
}

# Main execution
main() {
    echo "=== k8s-tray macOS Deployment and Test ==="
    echo "Host: $HOST"
    echo "Remote path: $REMOTE_PATH"
    echo ""

    check_app_bundle
    test_ssh
    create_remote_directory
    copy_app_bundle
    copy_test_scripts
    run_validation_test
    create_macos_test_script
    run_macos_test
    show_manual_instructions

    print_success "Deployment and testing completed successfully"
    print_info "The app bundle is ready for manual testing on $HOST"
}

main "$@"
