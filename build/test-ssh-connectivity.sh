#!/bin/bash

# Simple SSH connectivity test for macOS system

HOST=${1:-"mac-mini-m4.local"}
REMOTE_PATH=${2:-"~/src/k8s-tray"}

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

echo "=== SSH Connectivity Test ==="
echo "Host: $HOST"
echo "Remote path: $REMOTE_PATH"
echo ""

# Test basic SSH connection
print_info "Testing SSH connection..."
if ssh -o ConnectTimeout=10 "$HOST" "echo 'SSH connection successful'" 2>/dev/null; then
    print_success "SSH connection successful"
else
    print_error "SSH connection failed"
    echo ""
    echo "Troubleshooting steps:"
    echo "1. Make sure SSH is enabled on $HOST"
    echo "2. Check if you can ping the host: ping $HOST"
    echo "3. Try connecting manually: ssh $HOST"
    echo "4. Make sure SSH keys are set up or you have the password"
    exit 1
fi

# Test creating directory
print_info "Testing directory creation..."
if ssh "$HOST" "mkdir -p $REMOTE_PATH && echo 'Directory created successfully'" 2>/dev/null; then
    print_success "Directory creation successful"
else
    print_error "Directory creation failed"
    exit 1
fi

# Test file operations
print_info "Testing file operations..."
if ssh "$HOST" "echo 'test file' > $REMOTE_PATH/test.txt && rm $REMOTE_PATH/test.txt && echo 'File operations successful'" 2>/dev/null; then
    print_success "File operations successful"
else
    print_error "File operations failed"
    exit 1
fi

# Get system info
print_info "Getting system information..."
echo ""
echo "macOS System Info:"
ssh "$HOST" "sw_vers && echo 'Architecture:' \$(uname -m) && echo 'Hostname:' \$(hostname)"

print_success "SSH connectivity test completed successfully"
print_info "System is ready for deployment"
