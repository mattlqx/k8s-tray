# ADR-000: Mac Menu Bar Application Architecture

## Status
Proposed

## Date
2025-07-13

## Context
We need to develop a Mac menu bar application (k8s-tray) that provides quick access to Kubernetes cluster information and operations. The application must be compilable on non-Mac systems to support CI/CD pipelines and development workflows on Linux/Windows machines.

## Decision
We will build a Mac menu bar application using **Go** with the following architecture and toolchain:

### Primary Technology Stack
- **Language**: Go 1.21+
- **GUI Framework**: [Fyne](https://fyne.io/) with system tray support
- **Alternative**: [Wails v2](https://wails.io/) for web-based UI
- **Cross-compilation**: Go's native cross-compilation capabilities
- **Build System**: GitHub Actions with macOS runners for final packaging

### Architecture Components

#### 1. Core Application Structure
```
k8s-tray/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── kubernetes/          # K8s client and operations
│   ├── tray/                # System tray management
│   └── ui/                  # UI components
├── pkg/
│   └── models/              # Shared data models
├── assets/                  # Icons and resources
├── build/                   # Build scripts and configurations
└── dist/                    # Distribution artifacts
```

#### 2. System Tray Implementation
- **Primary**: Fyne's `systray` package for native macOS integration
- **Fallback**: Custom CGO bindings if needed for advanced features
- **Icons**: SVG-based with PNG fallbacks for different DPI settings

#### 3. Kubernetes Integration
- **Client**: `k8s.io/client-go` for cluster communication
- **Configuration**: Support for multiple kubeconfig files
- **Authentication**: OIDC, token-based, and certificate authentication

## Technical Decisions

### Cross-Platform Compilation Strategy
1. **Development Environment**: Any OS (Linux, Windows, macOS)
2. **Build Process**:
   - Local development builds on any platform
   - CI/CD builds macOS binaries for both amd64 and arm64 architectures
   - Universal binary creation using `lipo` tool for macOS distribution
   - Code signing and notarization in CI/CD pipeline

### Build Tools and Dependencies
- **Go Modules**: Dependency management
- **Make**: Build automation
- **goreleaser**: Release automation and cross-platform builds
- **GitHub Actions**: CI/CD pipeline
- **macOS Code Signing**: Apple Developer certificates in CI

### Code Quality Standards
- **Linting**: golangci-lint with strict configuration
- **Testing**: Unit tests with >80% coverage requirement
- **Documentation**: GoDoc comments for all public APIs
- **Git Hooks**: Pre-commit hooks for formatting and linting

### Universal Binary Strategy
- **Architecture Support**: Intel (amd64) and Apple Silicon (arm64) Macs
- **Build Process**: Separate compilation for each architecture followed by `lipo` merge
- **Distribution**: Single universal binary for simplified deployment
- **Performance**: Native performance on both Intel and Apple Silicon hardware

### Project Structure Standards
```go
// Example application structure
type App struct {
    config     *config.Config
    k8sClient  kubernetes.Interface
    tray       *systray.App
    ui         *ui.Manager
}

// Interfaces for testability
type KubernetesClient interface {
    GetPods(namespace string) ([]v1.Pod, error)
    GetServices(namespace string) ([]v1.Service, error)
}
```

## Implementation Plan

### Phase 1: Core Infrastructure
1. Set up Go module with proper project structure
2. Implement basic system tray with static menu
3. Configure cross-compilation build system
4. Set up CI/CD pipeline with macOS runners

### Phase 2: Kubernetes Integration
1. Implement kubeconfig loading and cluster connection
2. Add basic cluster information display
3. Implement pod/service listing functionality
4. Add real-time status updates

### Phase 3: Advanced Features
1. Multiple cluster support
2. Context switching
3. Custom action shortcuts
4. Configuration persistence

## Cross-Platform Build Configuration

### Makefile targets:
```makefile
.PHONY: build-darwin build-darwin-universal build-linux build-windows
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o dist/k8s-tray-darwin-amd64 ./cmd

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o dist/k8s-tray-darwin-arm64 ./cmd

build-darwin-universal: build-darwin-amd64 build-darwin-arm64
	lipo -create -output dist/k8s-tray-darwin-universal \
		dist/k8s-tray-darwin-amd64 \
		dist/k8s-tray-darwin-arm64

build-local:
	go build -o dist/k8s-tray ./cmd

test:
	go test ./...

lint:
	golangci-lint run
```

### GitHub Actions CI/CD:
- **Linux/Windows**: Cross-compile and test on ubuntu-latest
- **macOS**: Build binaries for both amd64 and arm64 architectures on macos-latest
- **Universal Binary**: Create fat binary using `lipo` for macOS distribution
- **Release**: Use goreleaser for multi-platform releases with universal macOS binary

## Alternatives Considered

### 1. Electron + TypeScript
- **Pros**: Web technologies, rich UI capabilities
- **Cons**: Large bundle size, memory overhead, complex cross-compilation

### 2. Swift (macOS native)
- **Pros**: Native performance, system integration
- **Cons**: macOS-only development, no cross-compilation

### 3. Rust + Tauri
- **Pros**: Small bundle size, good performance
- **Cons**: Learning curve, less mature ecosystem for system tray

### 4. Python + PyQt/Tkinter
- **Pros**: Rapid development, good libraries
- **Cons**: Distribution complexity, runtime dependencies

## Consequences

### Positive
- **Cross-platform development**: Developers can work on any OS
- **CI/CD friendly**: Automated builds without macOS requirement for most development
- **Performance**: Native Go performance with small memory footprint
- **Maintainability**: Simple deployment, single binary distribution
- **Testing**: Easy unit testing and mocking

### Negative
- **UI limitations**: System tray UI is inherently limited compared to full applications
- **macOS-specific features**: Some advanced macOS integrations may require CGO
- **Build complexity**: Final packaging still requires macOS environment

## Compliance and Security

### Code Signing
- Apple Developer certificate required for distribution
- Notarization process for macOS Catalina+ compatibility
- Universal binary signing for both Intel and Apple Silicon Macs
- Automated signing in CI/CD pipeline

### Security Considerations
- Secure storage of kubeconfig credentials
- Network security for cluster communications
- Regular dependency updates and vulnerability scanning

## Future Considerations
- **Multi-platform support**: Potential expansion to Linux/Windows system trays
- **Plugin system**: Extensibility for custom Kubernetes operations
- **Configuration UI**: Settings panel for advanced configuration
- **Metrics integration**: Integration with monitoring systems

## References
- [Fyne Documentation](https://developer.fyne.io/)
- [Go Cross Compilation](https://golang.org/doc/install/goos)
- [Kubernetes Client-Go](https://github.com/kubernetes/client-go)
- [macOS Code Signing Guide](https://developer.apple.com/documentation/xcode/notarizing_macos_software_before_distribution)

## Decision Outcome
This ADR provides a foundation for building a maintainable, cross-platform developable Mac menu bar application that follows Go best practices while enabling development on any operating system.