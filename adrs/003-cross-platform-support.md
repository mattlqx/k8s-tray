# ADR-003: Cross-Platform Support Strategy

## Status

Accepted - Implemented

## Date

2025-07-14

## Context

The k8s-tray application was initially designed for macOS as a menu bar application. However,
to maximize adoption and provide value to the broader Kubernetes community, we need to extend
support to Windows and Linux platforms. This requires careful consideration of platform-specific
system tray implementations, icon formats, user interface paradigms, and build/deployment
strategies while maintaining feature parity across all supported platforms.

## Decision

We will implement comprehensive cross-platform support targeting Windows, macOS, and Linux
operating systems for both ARM64 and AMD64 architectures. The application will maintain
consistent functionality across all platforms within the constraints of the underlying
system tray libraries and platform-specific requirements.

## Supported Platforms

### Operating Systems

1. **Windows**
   - Windows 10 (version 1909 and later)
   - Windows 11 (all versions)
   - Windows Server 2019 and later

2. **macOS**
   - macOS 10.15 (Catalina) and later
   - macOS 11 (Big Sure) and later
   - macOS 12 (Monterey) and later
   - macOS 13 (Ventura) and later
   - macOS 14 (Sonoma) and later

3. **Linux**
   - Ubuntu 20.04 LTS and later
   - Debian 11 and later
   - CentOS 8 and later / RHEL 8 and later
   - Fedora 35 and later
   - openSUSE Leap 15.3 and later
   - Arch Linux (rolling release)

### Architectures

- **AMD64/x86_64**: Primary architecture for all platforms
- **ARM64/aarch64**: Secondary architecture for all platforms
  - Apple Silicon (M1, M2, M3) on macOS
  - Windows on ARM
  - ARM64 Linux distributions

## Platform-Specific Implementation Details

### System Tray Integration

#### Library Choice: fyne.io/systray

We have chosen `fyne.io/systray` as our cross-platform system tray library because:

- **Cross-platform compatibility**: Native support for Windows, macOS, and Linux
- **Active maintenance**: Regular updates and bug fixes
- **Go-native**: Pure Go implementation with minimal CGO dependencies
- **Feature completeness**: Supports icons, tooltips, menus, and notifications
- **Community adoption**: Widely used in the Go ecosystem

#### Platform-Specific Behaviors

**Windows System Tray:**

- Uses Windows Shell API for system tray integration
- Requires ICO format icons for proper display and taskbar settings integration
- Supports rich tooltips with multi-line text
- Integrates with Windows notification area settings
- Respects Windows theme settings (light/dark mode)

**macOS Menu Bar:**

- Uses Cocoa/AppKit frameworks for menu bar integration
- Supports both regular and template icons for dark/light mode adaptation
- Requires proper app bundle structure for full functionality
- Integrates with macOS accessibility features
- Supports macOS-specific UI conventions

**Linux System Tray:**

- Uses DBus and StatusNotifierItem/AppIndicator specifications
- Compatibility varies by desktop environment:
  - GNOME: Requires extensions for system tray support
  - KDE Plasma: Native system tray support
  - XFCE: Built-in system tray support
  - i3/sway: Requires compatible status bar
- May require additional packages (e.g., `snixembed`) for older environments

### Icon Format Strategy

#### Platform-Specific Icon Formats

**Windows:**

```go
func createICOIcon(r, g, b uint8) []byte {
    // Generate ICO format icon with:
    // - 16x16 pixel size for system tray
    // - 32-bit BGRA color format
    // - Proper ICO header structure
    // - AND mask for transparency
}
```

**macOS:**

```go
func createPNGIcon(r, g, b uint8) []byte {
    // Generate PNG format icon with:
    // - 16x16 pixel size for menu bar
    // - RGBA color format
    // - Transparency support
    // - Template icon support for dark mode
}
```

**Linux:**

```go
func createPNGIcon(r, g, b uint8) []byte {
    // Generate PNG format icon with:
    // - 16x16 pixel size for system tray
    // - RGBA color format
    // - Transparency support
    // - Desktop environment compatibility
}
```

#### Runtime Icon Selection

```go
func createSimpleIcon(r, g, b uint8) []byte {
    if runtime.GOOS == "windows" {
        return createICOIcon(r, g, b)
    }
    return createPNGIcon(r, g, b)
}
```

### User Interface Adaptations

#### Menu Structure

**Common Elements (All Platforms):**

- Cluster health status indicator
- Cluster name and version information
- Namespace switching menu
- Context switching menu
- Pod status details with counts
- Refresh interval settings
- Manual refresh action
- Quit/Exit option

**Platform-Specific Elements:**

**Windows:**

- Help menu with Windows-specific tray icon instructions
- Administrative privilege detection
- Windows Explorer restart guidance
- System tray troubleshooting tips

**macOS:**

- App bundle requirement notices
- macOS security and privacy guidance
- Menu bar space management tips
- Accessibility integration

**Linux:**

- Desktop environment specific instructions
- System tray package recommendations
- DBus troubleshooting guidance
- Alternative system tray solutions

#### Tooltips and Help Text

**Windows:**

```go
func (m *Manager) showWindowsVisibilityHint() {
    tooltip := "K8s Tray - Connecting...\n\n💡 Windows Tip: If you don't see this icon, " +
              "check the system tray overflow area (^ arrow)\nand pin this icon for easier access."
    systray.SetTooltip(tooltip)
}
```

**macOS:**

```go
func (m *Manager) showMacOSAppBundleHint() {
    tooltip := "K8s Tray - Connecting...\n\n💡 macOS Tip: Use the .app bundle for " +
              "proper menu bar integration."
    systray.SetTooltip(tooltip)
}
```

**Linux:**

```go
func (m *Manager) showLinuxTrayHint() {
    tooltip := "K8s Tray - Connecting...\n\n💡 Linux Tip: Some desktop environments " +
              "may require additional packages for system tray support."
    systray.SetTooltip(tooltip)
}
```

### Build and Deployment Strategy

#### Cross-Compilation Setup

**Makefile Targets:**

```makefile
# Cross-compilation for all platforms
build-all: build-windows build-darwin build-linux

# Windows builds
build-windows:
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build ...
    CGO_ENABLED=1 GOOS=windows GOARCH=arm64 go build ...

# macOS builds
build-darwin:
    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build ...
    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build ...

# Linux builds
build-linux:
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ...
    CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build ...
```

#### Platform-Specific Build Requirements

**Windows:**

- MinGW-w64 cross-compiler for CGO
- Windows application manifest embedding
- Code signing certificates (optional)
- Windows Installer (MSI) generation

**macOS:**

- Xcode command line tools
- macOS SDK for cross-compilation
- App bundle generation with proper Info.plist
- Code signing and notarization (for distribution)

**Linux:**

- GCC cross-compilers for ARM64
- Static linking considerations
- AppImage/Flatpak packaging (optional)
- Distribution-specific package formats

### Platform-Specific Dependencies

#### CGO Dependencies

**Windows:**

- Windows API libraries (automatically linked)
- MinGW runtime libraries
- Windows manifest resources

**macOS:**

- Cocoa framework (automatically linked)
- AppKit framework (automatically linked)
- Core Foundation (automatically linked)

**Linux:**

- GTK+ libraries (for some system tray implementations)
- X11 libraries (for window system integration)
- DBus libraries (for system tray communication)

#### Runtime Dependencies

**Windows:**

- Windows 10/11 system tray infrastructure
- Windows Explorer shell integration
- No additional runtime dependencies

**macOS:**

- macOS system frameworks (built-in)
- Proper app bundle structure
- No additional runtime dependencies

**Linux:**

- System tray implementation (varies by desktop environment)
- DBus daemon (usually present)
- Optional: `snixembed` for legacy environments

### Testing Strategy

#### Platform-Specific Testing

**Automated Testing:**

```yaml
# GitHub Actions matrix testing
strategy:
  matrix:
    os: [windows-latest, macos-latest, ubuntu-latest]
    arch: [amd64, arm64]
    go-version: ['1.21']
```

**Manual Testing Checklist:**

**Windows:**

- [ ] Icon appears in system tray
- [ ] Icon visible in taskbar settings
- [ ] Tooltip displays correctly
- [ ] Menu opens and functions
- [ ] Administrative privileges work
- [ ] Windows Explorer restart behavior

**macOS:**

- [ ] Icon appears in menu bar
- [ ] App bundle functions properly
- [ ] Light/dark mode adaptation
- [ ] Menu bar space management
- [ ] Security permissions work
- [ ] Accessibility compliance

**Linux:**

- [ ] Icon appears in system tray (per DE)
- [ ] GNOME Shell extension compatibility
- [ ] KDE Plasma integration
- [ ] XFCE system tray support
- [ ] i3/sway status bar integration
- [ ] Legacy environment support

### Documentation Strategy

#### Platform-Specific Documentation

**README.md Sections:**

- Installation instructions per platform
- Platform-specific troubleshooting
- Build requirements per platform
- Distribution-specific notes

**Platform-Specific Guides:**

- `docs/WINDOWS.md`: Windows-specific setup and troubleshooting
- `docs/MACOS.md`: macOS-specific setup and app bundle requirements
- `docs/LINUX.md`: Linux desktop environment specific instructions

#### User Support Materials

**Windows:**

- System tray configuration guide
- Troubleshooting invisible icons
- Administrative privilege requirements
- Windows version compatibility

**macOS:**

- App bundle importance explanation
- Security and privacy settings
- Menu bar management tips
- Development vs. distribution differences

**Linux:**

- Desktop environment compatibility matrix
- System tray package requirements
- DBus troubleshooting guide
- Alternative system tray solutions

### Performance Considerations

#### Platform-Specific Optimizations

**Windows:**

- Efficient ICO icon generation
- Minimal Windows API calls
- Proper resource cleanup
- Memory usage monitoring

**macOS:**

- Efficient PNG icon generation
- Template icon usage for themes
- Minimal Cocoa framework usage
- App bundle optimization

**Linux:**

- Efficient PNG icon generation
- DBus message optimization
- Desktop environment detection
- Fallback mechanism implementation

#### Resource Usage

**Memory Usage:**

- Windows: ~15-20MB (including system tray integration)
- macOS: ~12-18MB (including menu bar integration)
- Linux: ~15-25MB (varies by desktop environment)

**CPU Usage:**

- Polling interval optimization (5s-5min configurable)
- Efficient Kubernetes API calls
- Minimal background processing
- Event-driven updates where possible

### Distribution Strategy

#### Release Artifacts

**Windows:**

- `k8s-tray-windows-amd64.exe`
- `k8s-tray-windows-arm64.exe`
- Windows Installer (MSI) packages
- Chocolatey package (future)

**macOS:**

- `k8s-tray-darwin-amd64`
- `k8s-tray-darwin-arm64`
- `k8s-tray-darwin-universal` (fat binary)
- `K8s Tray.app` (app bundle)
- Homebrew formula (future)

**Linux:**

- `k8s-tray-linux-amd64`
- `k8s-tray-linux-arm64`
- AppImage packages (future)
- Distribution-specific packages (future)

#### Distribution Channels

**GitHub Releases:**

- Automated release builds
- Platform-specific binaries
- Checksums and signatures
- Release notes with platform-specific changes

**Package Managers:**

- Homebrew (macOS)
- Chocolatey (Windows)
- Snap Store (Linux)
- APT/YUM repositories (Linux)

### Maintenance and Support

#### Platform-Specific Maintenance

**Windows:**

- Windows Update compatibility testing
- Windows Defender exclusion guidance
- System tray API changes monitoring
- Enterprise deployment support

**macOS:**

- macOS version compatibility testing
- App bundle requirement maintenance
- Code signing certificate renewal
- macOS security policy changes

**Linux:**

- Desktop environment compatibility updates
- Distribution-specific testing
- Package dependency management
- System tray specification changes

#### Long-term Support Strategy

**Version Support:**

- Windows: Support current and previous major version
- macOS: Support current and two previous major versions
- Linux: Support LTS versions plus recent releases

**Deprecation Policy:**

- 12-month advance notice for platform deprecation
- Migration guidance for unsupported platforms
- Community support for legacy platforms

## Implementation Phases

### Phase 1: Core Cross-Platform Support (Completed)

- [x] Windows system tray integration with ICO icons
- [x] macOS menu bar integration with PNG icons
- [x] Linux system tray integration with PNG icons
- [x] Platform-specific icon generation
- [x] Cross-compilation build system
- [x] Platform-specific user guidance

### Phase 2: Platform Optimization (Current)

- [ ] Windows application manifest integration
- [ ] macOS app bundle automation
- [ ] Linux desktop environment detection
- [ ] Platform-specific testing automation
- [ ] Distribution package generation

### Phase 3: Enhanced Platform Integration (Future)

- [ ] Windows notification integration
- [ ] macOS notification center integration
- [ ] Linux desktop notification support
- [ ] Platform-specific keyboard shortcuts
- [ ] System theme integration

### Phase 4: Distribution and Packaging (Future)

- [ ] Windows MSI installer
- [ ] macOS DMG distribution
- [ ] Linux package repositories
- [ ] Automated package manager submissions
- [ ] Code signing and notarization

## Success Metrics

### Platform Adoption

- **Windows**: 40% of user base (target)
- **macOS**: 35% of user base (target)
- **Linux**: 25% of user base (target)

### Architecture Distribution

- **AMD64**: 75% of installations (expected)
- **ARM64**: 25% of installations (expected)

### Platform-Specific Quality Metrics

- **Windows**: <5% system tray visibility issues
- **macOS**: <2% app bundle related issues
- **Linux**: <10% desktop environment compatibility issues

### User Experience Metrics

- **Cross-platform feature parity**: 100%
- **Platform-specific help effectiveness**: >90% issue resolution
- **Cross-compilation reliability**: >98% successful builds

## Risks and Mitigation

### Technical Risks

**System Tray Library Changes:**

- Risk: Breaking changes in fyne.io/systray
- Mitigation: Pin to stable versions, maintain compatibility layer

**Platform API Changes:**

- Risk: OS updates breaking system tray integration
- Mitigation: Regular compatibility testing, community feedback

**CGO Compilation Issues:**

- Risk: Cross-compilation failures
- Mitigation: Docker-based build environments, comprehensive CI/CD

### Operational Risks

**Support Burden:**

- Risk: Increased support complexity across platforms
- Mitigation: Platform-specific documentation, community support

**Distribution Complexity:**

- Risk: Multiple distribution channels to maintain
- Mitigation: Automated release processes, phased rollout

**Testing Coverage:**

- Risk: Platform-specific bugs in production
- Mitigation: Comprehensive testing matrix, beta testing program

## Decision Outcome

This ADR establishes k8s-tray as a truly cross-platform application supporting Windows,
macOS, and Linux on both AMD64 and ARM64 architectures. The implementation provides
consistent functionality across all platforms while respecting platform-specific user
interface conventions and technical requirements. The approach balances feature parity
with platform optimization, ensuring a native experience for users regardless of their
operating system choice.

The cross-platform support strategy positions k8s-tray to serve the entire Kubernetes
community, from Windows-based enterprise developers to Linux-running DevOps engineers
and macOS-using application developers, all while maintaining the high-quality user
experience established in the initial macOS implementation.
