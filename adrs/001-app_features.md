# ADR-001: Application Features and User Interface

## Status

Accepted - Implemented

## Date

2025-07-13

## Context

Building upon ADR-000's architectural decisions, we need to define the specific features
and user interface requirements for the k8s-tray Mac menu bar application. The application
should provide quick visual feedback about Kubernetes cluster health and allow users to
access detailed cluster information through a simple menu interface.

## Decision

We will implement a Mac menu bar application with the following feature set:

## Core Features

### 1. Health Status Indicator

The menu bar icon will serve as an active health indicator with color-coded status:

#### Status Colors

- **Red**: Critical issues detected
  - Any pods in `CrashLoopBackOff` state
  - Any pods failing readiness checks
  - Any pods in `Error` or `Failed` state
- **Yellow**: Warning conditions present
  - Any pods in `Pending` state
  - Any pods in `ContainerCreating` state
  - Any pods in `Terminating` state
  - Any pods in other non-`Running` states
- **Green**: All systems healthy
  - All pods in `Running` state
  - All readiness checks passing

#### Visual Implementation

- Square icon in the menu bar
- Color changes reflect real-time cluster status
- Smooth color transitions (200ms animation)
- High contrast colors for accessibility

### 2. Interactive Menu Interface

Clicking the menu bar icon reveals a dropdown menu with:

#### Primary Information Display

- **Cluster Name**: Current active cluster context
- **Overall Status**: Text description of current health state
- **Kubernetes Version**: Server version of the cluster
- **Pod Summary**: Count of pods by state (Running, Pending, Failed, etc.)
- **Resource Usage**: CPU and memory utilization (when available)

#### Menu Structure

```text
┌─────────────────────────────────────┐
│ Cluster: production-cluster         │
│ Status: ● Healthy (24/24 running)   │
│ Version: v1.28.3                   │
├─────────────────────────────────────┤
│ Pods:                              │
│   Running: 24                      │
│   Pending: 0                       │
│   Failed: 0                        │
├─────────────────────────────────────┤
│ Resources:                         │
│   CPU: 45% (12/27 cores)          │
│   Memory: 62% (8.3/13.5 GB)       │
├─────────────────────────────────────┤
│ ⚙️  Settings                        │
│ 🔄 Refresh Now                     │
│ ❌ Quit                            │
└─────────────────────────────────────┘
```

### 3. Application Visibility

The application will be a menu bar-only application:

- **No Dock presence**: Application will not appear in the Dock
- **Menu bar only**: Accessible solely through the menu bar icon
- **Background operation**: Runs silently in the background
- **System integration**: Proper macOS menu bar application behavior

### 4. Kubernetes Configuration Management

The application will use standard Kubernetes configuration practices:

#### Default Configuration

- **Primary location**: `~/.kube/config`
- **Standard kubeconfig format**: Full compatibility with kubectl
- **Context switching**: Support for multiple contexts within kubeconfig
- **No credential storage**: Application does not store or manage credentials

#### Alternative Configuration

- **Custom kubeconfig path**: Configurable via settings
- **Environment variable support**: `KUBECONFIG` environment variable
- **Multiple config files**: Support for merged kubeconfig files

#### Configuration Priority

1. Custom path specified in application settings
2. `KUBECONFIG` environment variable
3. Default `~/.kube/config` location

### 5. Cluster Context Selection

The application will provide intuitive cluster context management:

#### Context Display

- **Current Context**: Clearly displayed in menu header
- **Visual Indicator**: Active context highlighted with distinctive styling
- **Context Information**: Show cluster name, server URL, and namespace

#### Context Switching

- **Dropdown Menu**: Accessible context selector within main menu
- **Quick Switch**: One-click context changing
- **Immediate Update**: Status refreshes immediately after context change
- **Persistence**: Remember last selected context across application restarts

#### Enhanced Menu Structure

```text
┌─────────────────────────────────────┐
│ Context: production-cluster ▼       │
│ ├─ dev-cluster                     │
│ ├─ staging-cluster                 │
│ └─ ● production-cluster (active)   │
│ Status: ● Healthy (24/24 running)   │
│ Version: v1.28.3                   │
├─────────────────────────────────────┤
│ Pods:                              │
│   Running: 24                      │
│   Pending: 0                       │
│   Failed: 0                        │
├─────────────────────────────────────┤
│ Resources:                         │
│   CPU: 45% (12/27 cores)          │
│   Memory: 62% (8.3/13.5 GB)       │
├─────────────────────────────────────┤
│ ⚙️  Settings                        │
│ 🔄 Refresh Now                     │
│ ❌ Quit                            │
└─────────────────────────────────────┘
```

#### Future Multi-Cluster Support

- **Multiple Menu Bar Icons**: Option to display separate icons for different clusters
- **Cluster Grouping**: Organize contexts by environment (dev, staging, prod)
- **Simultaneous Monitoring**: Monitor multiple clusters concurrently
- **Aggregated Status**: Combined health status across selected clusters

## Technical Implementation Details

### 1. Health Status Monitoring

```go
type ClusterHealth struct {
    Status      HealthStatus
    PodCounts   PodStatusCounts
    Version     string
    LastUpdated time.Time
}

type HealthStatus int

const (
    HealthStatusGreen HealthStatus = iota
    HealthStatusYellow
    HealthStatusRed
)

type PodStatusCounts struct {
    Running           int
    Pending           int
    Failed            int
    CrashLoopBackOff  int
    ContainerCreating int
    Terminating       int
}
```

### 2. Menu Bar Integration

- **Framework**: Fyne's system tray with custom menu implementation
- **Update frequency**: 30-second interval (configurable)
- **Error handling**: Graceful degradation when cluster is unreachable
- **Caching**: Local cache of last known state for offline scenarios

### 3. Application Lifecycle

```go
type App struct {
    tray          *systray.App
    k8sClient     kubernetes.Interface
    healthMonitor *HealthMonitor
    settings      *Settings
}

// App configuration
type Settings struct {
    KubeconfigPath   string
    RefreshInterval  time.Duration
    ShowResourceInfo bool
    LastContext      string
    MultiClusterMode bool
}

// Context management
type ContextManager struct {
    availableContexts []KubeContext
    activeContext     string
    kubeConfig        *rest.Config
}

type KubeContext struct {
    Name      string
    Cluster   string
    Server    string
    Namespace string
    User      string
}
```

### 4. macOS Integration

- **LSUIElement**: Set to `true` to hide from Dock
- **NSApplication**: Menu bar application type
- **Accessibility**: VoiceOver support for menu items
- **Dark mode**: Icon adapts to system appearance

## User Interaction Flows

### 1. First Launch

1. Check for kubeconfig at default location
2. If not found, prompt for configuration
3. Test connection and display initial status
4. Start background monitoring

### 2. Normal Operation

1. Monitor cluster health every 30 seconds
2. Update menu bar icon color based on status
3. Cache results for quick menu display
4. Handle network interruptions gracefully

### 3. Configuration Changes

1. Settings accessible via menu
2. Kubeconfig path changes trigger reconnection
3. Context switching updates monitoring target
4. Validation of configuration before applying

### 4. Context Switching

1. User clicks on context dropdown in menu
2. Available contexts loaded from kubeconfig
3. User selects new context
4. Application switches to new context immediately
5. Health monitoring restarts for new cluster
6. Menu updates to reflect new cluster status
7. Selected context persisted for next application start

## Error Handling and Edge Cases

### 1. Network Connectivity

- **Offline mode**: Show last known status with timestamp
- **Connection errors**: Display error state in menu
- **Timeout handling**: 10-second timeout for API calls

### 2. Configuration Issues

- **Invalid kubeconfig**: Clear error messages in menu
- **Authentication failures**: Prompt for credential refresh
- **Permission errors**: Informative error descriptions

### 3. Performance Considerations

- **Background threads**: Non-blocking UI updates
- **Memory usage**: Efficient caching and cleanup
- **CPU usage**: Minimal impact on system performance

## Accessibility and Usability

### 1. Visual Accessibility

- **High contrast colors**: WCAG AA compliant color choices
- **Color blind support**: Additional visual indicators beyond color
- **Font sizing**: Respect system font size preferences

### 2. Keyboard Navigation

- **Menu navigation**: Full keyboard support
- **Shortcuts**: Common actions accessible via keyboard
- **Focus management**: Proper tab order and focus indicators

## Security Considerations

### 1. Credential Handling

- **No credential storage**: Rely on system keychain integration
- **Secure transmission**: TLS for all Kubernetes API communications
- **Permission model**: Minimal required permissions

### 2. Data Privacy

- **No telemetry**: No data collection or transmission
- **Local processing**: All data processing happens locally
- **Secure defaults**: Conservative security settings

## Future Extensibility

### 1. Customization Options

- **Refresh intervals**: User-configurable update frequency
- **Display options**: Toggleable information elements
- **Notification preferences**: Optional desktop notifications

### 2. Advanced Features

- **Multiple cluster icons**: Separate menu bar icons for simultaneous cluster monitoring
- **Context grouping**: Organize contexts by environment or team
- **Custom health checks**: User-defined health criteria per cluster
- **Export capabilities**: Status reports and logs
- **Cluster comparison**: Side-by-side cluster health comparison
- **Alert notifications**: Desktop notifications for cluster state changes

## Testing Strategy

### 1. Unit Testing

- Health status calculation logic
- Configuration parsing and validation
- Error handling scenarios

### 2. Integration Testing

- Kubernetes API integration
- Menu bar behavior testing
- Configuration file handling

### 3. User Acceptance Testing

- Usability testing with real Kubernetes clusters
- Performance testing under various load conditions
- Accessibility testing with assistive technologies

## Success Metrics

### 1. Performance Metrics

- **Startup time**: < 2 seconds to first status display
- **Update latency**: < 5 seconds for status changes
- **Memory usage**: < 50MB resident memory
- **CPU usage**: < 1% average CPU utilization

### 2. User Experience Metrics

- **Ease of setup**: Single-click configuration for standard setups
- **Status accuracy**: 99.9% correlation with actual cluster state
- **Reliability**: No crashes during normal operation

## Decision Outcome

This ADR defines a focused, user-friendly Mac menu bar application that provides essential
Kubernetes cluster health information through an intuitive visual interface while maintaining
security best practices and system integration standards.
