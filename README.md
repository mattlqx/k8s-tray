# k8s-tray

A Mac menu bar application for monitoring Kubernetes cluster health and status.

## Features

- **Health Status Indicator**: Color-coded menu bar icon showing cluster health
  - 🟢 Green: All pods running and healthy
  - 🟡 Yellow: Some pods in warning states (pending, terminating)
  - 🔴 Red: Critical issues (failed pods, crashes)
- **Cluster Information**: Quick access to cluster name, version, and context
- **Pod Status**: Real-time pod counts and status overview
- **Namespace Switching**: Easy namespace selection from dropdown menu
- **Auto-refresh**: Configurable polling interval for status updates
- **Multi-cluster Support**: Switch between different Kubernetes contexts

## Installation

### Prerequisites

- macOS 10.15 or later
- kubectl configured with access to your Kubernetes cluster
- Valid kubeconfig file

### Download

1. Download the latest release from the [releases page](https://github.com/k8s-tray/k8s-tray/releases)
2. Extract the archive
3. Move `k8s-tray` to your Applications folder or `/usr/local/bin`

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mattlqx/k8s-tray.git
cd k8s-tray

# Install dependencies
make deps

# Build the application for current platform
make build

# Build for macOS (creates proper app bundle)
make build-darwin-app

# Run the application
./dist/k8s-tray
```

### macOS App Bundle

For macOS, it's recommended to use the app bundle instead of running the raw binary:

```bash
# Build the macOS app bundle
make build-darwin-app

# Run the app bundle
open "dist/K8s Tray.app"

# Or install to Applications folder
cp -r "dist/K8s Tray.app" /Applications/
```

**Important**: On macOS, the system tray functionality requires the application to be packaged as a
proper `.app` bundle. Running the raw binary directly may not show the menu bar icon. Always use the
app bundle for macOS deployment.

### Troubleshooting macOS

If the menu bar icon doesn't appear:

1. Make sure you're using the `.app` bundle, not the raw binary
2. Check macOS Security & Privacy settings if the app is blocked
3. See [macOS Troubleshooting Guide](docs/MACOS_TROUBLESHOOTING.md) for detailed help

## Configuration

k8s-tray uses a YAML configuration file located at `~/.k8s-tray.yaml`. The configuration
file is created automatically with default values on first run.

### Configuration Options

```yaml
# Kubernetes configuration
kubeconfig: ~/.kube/config      # Path to kubeconfig file
context: ""                     # Kubernetes context (empty = current context)
namespace: "default"            # Default namespace to monitor

# Polling configuration
poll_interval: 5s               # How often to refresh cluster status

# UI configuration
show_notifications: true        # Show desktop notifications
theme: "auto"                   # Theme: auto, light, dark

# Feature flags
show_metrics: true              # Show resource metrics (if available)
show_logs: false                # Show log viewer (future feature)
show_events: true               # Show recent events (future feature)
```

## Usage

### Running the Application

```bash
# Run directly
./k8s-tray

# Run with custom kubeconfig
KUBECONFIG=/path/to/config ./k8s-tray

# Run with specific context
./k8s-tray --context=my-cluster
```

### Menu Options

- **Status**: Shows current cluster health status
- **Cluster**: Displays cluster name and version
- **Namespace**: Shows current namespace
- **Pods**: Pod count summary
- **Switch Namespace**: Dropdown to select different namespace
- **Refresh**: Manually refresh cluster status
- **Settings**: Open configuration (future feature)
- **Quit**: Exit the application

### Status Indicators

| Color | Status | Description |
|-------|--------|-------------|
| 🟢 Green | Healthy | All pods running, no issues detected |
| 🟡 Yellow | Warning | Some pods pending, creating, or terminating |
| 🔴 Red | Critical | Failed pods, crashes, or other critical issues |
| ⚫ Gray | Unknown | Unable to connect or determine status |

## Development

### Development Prerequisites

- Go 1.21 or later
- pre-commit (for development)
- golangci-lint

### Setup Development Environment

```bash
# Install development dependencies
make setup-dev

# Install pre-commit hooks
make pre-commit-install

# Run tests
make test

# Run linter
make lint

# Format code
make format
```

### Project Structure

```text
k8s-tray/
├── cmd/                    # Application entry point
│   └── main.go
├── internal/               # Internal application code
│   ├── config/            # Configuration management
│   ├── kubernetes/        # Kubernetes client and operations
│   ├── tray/              # System tray management
│   └── ui/                # UI components
├── pkg/                   # Shared packages
│   └── models/            # Data models
├── assets/                # Icons and resources
├── build/                 # Build scripts
├── dist/                  # Build output
├── adrs/                  # Architecture Decision Records
└── .github/               # GitHub Actions workflows
```

### Architecture

k8s-tray follows a modular architecture:

- **Configuration Layer**: Handles application settings and kubeconfig management
- **Kubernetes Layer**: Manages cluster connections and API interactions
- **Tray Layer**: Handles system tray integration and menu management
- **UI Layer**: Future extensibility for settings dialogs and detailed views

### Adding Features

1. Create an ADR (Architecture Decision Record) in the `adrs/` directory
2. Implement the feature following the existing patterns
3. Add tests for new functionality
4. Update documentation
5. Submit a pull request

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -race -cover ./...

# Run specific package tests
go test ./internal/kubernetes/
```

### Code Quality

This project uses pre-commit hooks to ensure code quality:

- **Go formatting**: gofmt and goimports
- **Linting**: golangci-lint with comprehensive rules
- **Security scanning**: gosec
- **Commit messages**: Conventional commit format
- **Documentation**: Markdown linting and spell checking

### Cross-compilation

```bash
# Build for macOS (both architectures)
make build-darwin

# Build for all platforms
make build-all
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

Please ensure your commits follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Report bugs: [GitHub Issues](https://github.com/k8s-tray/k8s-tray/issues)
- Feature requests: [GitHub Discussions](https://github.com/k8s-tray/k8s-tray/discussions)
- Documentation: [Wiki](https://github.com/k8s-tray/k8s-tray/wiki)

## Acknowledgments

- [Fyne](https://fyne.io/) for the cross-platform GUI framework
- [client-go](https://github.com/kubernetes/client-go) for Kubernetes API client
- [systray](https://github.com/getlantern/systray) for system tray functionality

## Roadmap

- [ ] Settings dialog UI
- [ ] Pod log viewer
- [ ] Event viewer
- [ ] Resource metrics display
- [ ] Multi-cluster management
- [ ] Desktop notifications
- [ ] Plugin system
- [ ] Windows and Linux support
