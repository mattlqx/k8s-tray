package tray

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"fyne.io/systray"
	"github.com/k8s-tray/k8s-tray/internal/config"
	"github.com/k8s-tray/k8s-tray/internal/kubernetes"
	"github.com/k8s-tray/k8s-tray/pkg/models"
)

const osWindows = "windows"

// Pod phase constants
const (
	podPhaseRunning   = "Running"
	podPhasePending   = "Pending"
	podPhaseSucceeded = "Succeeded"
	podPhaseFailed    = "Failed"
)

// Manager handles the system tray functionality
type Manager struct {
	k8sClient *kubernetes.Client
	config    *config.Config

	// Menu items
	statusItem        *systray.MenuItem
	clusterItem       *systray.MenuItem
	namespaceItem     *systray.MenuItem
	podsItem          *systray.MenuItem
	cpuItem           *systray.MenuItem
	memoryItem        *systray.MenuItem
	podsReadyItem     *systray.MenuItem
	podsNotReadyItem  *systray.MenuItem
	podsPendingItem   *systray.MenuItem
	podsCompletedItem *systray.MenuItem
	podsFailedItem    *systray.MenuItem
	refreshItem       *systray.MenuItem
	helpItem          *systray.MenuItem
	quitItem          *systray.MenuItem

	// Namespace submenu items
	namespaceMenu      *systray.MenuItem
	namespaceItems     map[string]*systray.MenuItem
	namespaceSeparator *systray.MenuItem

	// Context submenu items
	contextMenu  *systray.MenuItem
	contextItems map[string]*systray.MenuItem

	// Settings submenu items
	settingsMenu  *systray.MenuItem
	intervalItems map[time.Duration]*systray.MenuItem

	// Pod submenu items for each state
	podsReadySubmenu     map[string]*systray.MenuItem
	podsNotReadySubmenu  map[string]*systray.MenuItem
	podsPendingSubmenu   map[string]*systray.MenuItem
	podsCompletedSubmenu map[string]*systray.MenuItem
	podsFailedSubmenu    map[string]*systray.MenuItem

	// Monitoring control
	intervalChanged chan time.Duration

	// Current state
	currentStatus *models.ClusterStatus
	currentHealth models.HealthStatus

	// Windows-specific visibility helper
	showVisibilityHint bool
}

// NewManager creates a new tray manager
func NewManager(k8sClient *kubernetes.Client, cfg *config.Config) *Manager {
	return &Manager{
		k8sClient:            k8sClient,
		config:               cfg,
		namespaceItems:       make(map[string]*systray.MenuItem),
		contextItems:         make(map[string]*systray.MenuItem),
		intervalItems:        make(map[time.Duration]*systray.MenuItem),
		podsReadySubmenu:     make(map[string]*systray.MenuItem),
		podsNotReadySubmenu:  make(map[string]*systray.MenuItem),
		podsPendingSubmenu:   make(map[string]*systray.MenuItem),
		podsCompletedSubmenu: make(map[string]*systray.MenuItem),
		podsFailedSubmenu:    make(map[string]*systray.MenuItem),
		intervalChanged:      make(chan time.Duration, 1),
		currentHealth:        models.HealthUnknown,
		showVisibilityHint:   runtime.GOOS == osWindows, // Show hint only on Windows
	}
}

// OnReady is called when the systray is ready
func (m *Manager) OnReady(ctx context.Context) {
	log.Printf("Tray manager OnReady called")

	// Set initial icon and tooltip
	m.updateIcon(models.HealthUnknown)
	systray.SetTooltip("K8s Tray - Connecting...")

	log.Printf("Set initial icon and tooltip")

	// Build menu
	m.buildMenu()

	log.Printf("Built menu")

	// Initialize namespace menu
	go m.refreshNamespaceMenu(ctx)

	log.Printf("Initialized namespace menu")

	// Initialize context menu
	go m.refreshContextMenu(ctx)

	log.Printf("Initialized context menu")

	// Initialize settings menu
	go m.refreshSettingsMenu(ctx)

	log.Printf("Initialized settings menu")

	// Start monitoring
	go m.startMonitoring(ctx)

	log.Printf("Started monitoring")

	// Handle menu actions
	go m.handleMenuActions(ctx)

	log.Printf("Started menu action handler")

	// Show Windows-specific startup hint in tooltip
	if runtime.GOOS == osWindows {
		m.showWindowsVisibilityHint()
	}
}

// showWindowsVisibilityHint enhances the initial tooltip for Windows users
func (m *Manager) showWindowsVisibilityHint() {
	// Set an initial helpful tooltip for Windows users
	if runtime.GOOS == osWindows {
		systray.SetTooltip("K8s Tray - Connecting...\n\n💡 Windows Tip: If you don't see this icon, check the system tray overflow area (^ arrow)\nand pin this icon for easier access. See Help menu for details.")
	} else {
		systray.SetTooltip("K8s Tray - Connecting...")
	}

	// After 15 seconds, revert to normal tooltip behavior
	go func() {
		time.Sleep(15 * time.Second)
		// This will be overridden by the normal status updates anyway
	}()
}

// OnExit is called when the systray is exiting
func (m *Manager) OnExit() {
	log.Println("Tray exiting...")
}

// buildMenu builds the system tray menu
func (m *Manager) buildMenu() {
	// Status information
	m.statusItem = systray.AddMenuItem("Status: Connecting...", "Current cluster status")
	m.statusItem.Disable()

	m.clusterItem = systray.AddMenuItem("Cluster: Unknown", "Current cluster")
	m.clusterItem.Disable()

	// Get display name for namespace
	namespaceDisplay := m.config.Namespace
	if m.config.Namespace == config.AllNamespaces {
		namespaceDisplay = "All Namespaces"
	}

	m.namespaceItem = systray.AddMenuItem("Namespace: "+namespaceDisplay, "Current namespace")
	m.namespaceItem.Disable()

	// Resource usage items (only show if metrics are enabled)
	if m.config.ShowMetrics {
		m.cpuItem = systray.AddMenuItem("CPU: Loading...", "CPU usage across all cluster nodes")
		m.cpuItem.Disable()

		m.memoryItem = systray.AddMenuItem("Memory: Loading...", "Memory usage across all cluster nodes")
		m.memoryItem.Disable()
	}

	m.podsItem = systray.AddMenuItem("Pods: Loading...", "Pod status summary")
	m.podsItem.Disable()

	// Individual pod status items with better tooltips
	m.podsReadyItem = systray.AddMenuItem("  🟢 Ready: 0", "Pods that are running and all containers are ready")
	// Keep enabled to allow submenu access on macOS

	m.podsNotReadyItem = systray.AddMenuItem("  🛑 Not Ready: 0", "Pods that are running but some containers are not ready")
	// Keep enabled to allow submenu access on macOS

	m.podsPendingItem = systray.AddMenuItem("  ⏳ Pending: 0", "Pods that are waiting to be scheduled or start")
	// Keep enabled to allow submenu access on macOS

	m.podsCompletedItem = systray.AddMenuItem("  ✅ Completed: 0", "Pods that have completed their work successfully")
	// Keep enabled to allow submenu access on macOS

	m.podsFailedItem = systray.AddMenuItem("  ❌ Failed: 0", "Pods that have failed to start or run")
	// Keep enabled to allow submenu access on macOS

	systray.AddSeparator()

	// Namespace selection
	m.namespaceMenu = systray.AddMenuItem("Switch Namespace", "Switch to different namespace")

	// Context selection
	m.contextMenu = systray.AddMenuItem("Switch Context", "Switch to different cluster context")

	systray.AddSeparator()

	// Actions
	m.refreshItem = systray.AddMenuItem("Refresh", "Refresh cluster status")
	m.settingsMenu = systray.AddMenuItem("Settings", "Application settings")

	// Add help for Windows users
	if runtime.GOOS == osWindows {
		m.helpItem = systray.AddMenuItem("Help", "Tips for using K8s Tray on Windows")
	}

	systray.AddSeparator()

	m.quitItem = systray.AddMenuItem("Quit", "Quit K8s Tray")
}

// handleMenuActions handles menu item clicks
func (m *Manager) handleMenuActions(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.refreshItem.ClickedCh:
			go m.refreshStatus(ctx)
		case <-m.quitItem.ClickedCh:
			systray.Quit()
			return
		case <-m.namespaceMenu.ClickedCh:
			go m.refreshNamespaceMenu(ctx)
		case <-m.contextMenu.ClickedCh:
			go m.refreshContextMenu(ctx)
		case <-m.settingsMenu.ClickedCh:
			go m.refreshSettingsMenu(ctx)
		case <-m.podsReadyItem.ClickedCh:
			// Pod status items are now clickable but we don't need to do anything
			// The submenus will be handled automatically by the systray library
		case <-m.podsNotReadyItem.ClickedCh:
			// Pod status items are now clickable but we don't need to do anything
		case <-m.podsPendingItem.ClickedCh:
			// Pod status items are now clickable but we don't need to do anything
		case <-m.podsCompletedItem.ClickedCh:
			// Pod status items are now clickable but we don't need to do anything
		case <-m.podsFailedItem.ClickedCh:
			// Pod status items are now clickable but we don't need to do anything
		}

		// Handle Windows help menu if it exists
		if m.helpItem != nil {
			select {
			case <-m.helpItem.ClickedCh:
				go m.showWindowsHelp()
			default:
			}
		}
	}
}

// startMonitoring starts the periodic monitoring of cluster status
func (m *Manager) startMonitoring(ctx context.Context) {
	// Initial refresh
	m.refreshStatus(ctx)

	// Set up periodic refresh with dynamic interval changes
	ticker := time.NewTicker(m.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.refreshStatus(ctx)
		case newInterval := <-m.intervalChanged:
			// Reset ticker with new interval
			ticker.Stop()
			ticker = time.NewTicker(newInterval)
			log.Printf("Updated monitoring interval to: %s", newInterval)
		}
	}
}

// refreshStatus refreshes the cluster status
func (m *Manager) refreshStatus(ctx context.Context) {
	status, err := m.k8sClient.GetClusterStatus(ctx)
	if err != nil {
		log.Printf("Failed to get cluster status: %v", err)
		m.updateError(err)
		return
	}

	log.Printf("Refreshed cluster status... %+v", status.PodStatus)

	m.currentStatus = status
	m.updateDisplay(status)
}

// updateDisplay updates the tray display with current status
func (m *Manager) updateDisplay(status *models.ClusterStatus) {
	// Update icon if health status changed
	if status.HealthStatus != m.currentHealth {
		m.updateIcon(status.HealthStatus)
		m.currentHealth = status.HealthStatus
	}

	// Get display name for namespace
	namespaceDisplay := m.config.Namespace
	if m.config.Namespace == config.AllNamespaces {
		namespaceDisplay = "All Namespaces"
	}

	// Update tooltip with Windows-specific guidance if applicable
	tooltip := fmt.Sprintf("K8s Tray - %s\nCluster: %s (%s)\nNamespace: %s\nPods: %d total",
		status.HealthStatus.String(),
		status.ClusterName,
		status.ServerVersion,
		namespaceDisplay,
		status.PodStatus.Total)

	// Add resource stats to tooltip if available
	if m.config.ShowMetrics && status.Resources != nil {
		if status.Resources.CPU != nil {
			tooltip += fmt.Sprintf("\nCPU: %.1f/%.1f cores (%.1f%%)",
				status.Resources.CPU.Used,
				status.Resources.CPU.Available,
				status.Resources.CPU.Percentage)
		}
		if status.Resources.Memory != nil {
			tooltip += fmt.Sprintf("\nMemory: %.1f/%.1f GB (%.1f%%)",
				status.Resources.Memory.Used,
				status.Resources.Memory.Available,
				status.Resources.Memory.Percentage)
		}
	}

	// Add Windows-specific visibility hint if needed
	if runtime.GOOS == osWindows && m.showVisibilityHint {
		tooltip += "\n\n💡 Tip: Pin this icon to the visible tray area for easier access"
		// Only show this hint for the first few status updates
		m.showVisibilityHint = false
	}

	systray.SetTooltip(tooltip)

	// Update menu items
	m.statusItem.SetTitle(fmt.Sprintf("Status: %s", status.HealthStatus.String()))
	m.clusterItem.SetTitle(fmt.Sprintf("Cluster: %s (%s)", status.ClusterName, status.ServerVersion))
	m.namespaceItem.SetTitle(fmt.Sprintf("Namespace: %s", namespaceDisplay))

	// Update resource stats if enabled and available
	if m.config.ShowMetrics && status.Resources != nil {
		if status.Resources.CPU != nil {
			m.cpuItem.SetTitle(fmt.Sprintf("CPU: %.1f/%.1f cores (%.1f%%)",
				status.Resources.CPU.Used,
				status.Resources.CPU.Available,
				status.Resources.CPU.Percentage))
		}
		if status.Resources.Memory != nil {
			m.memoryItem.SetTitle(fmt.Sprintf("Memory: %.1f/%.1f GB (%.1f%%)",
				status.Resources.Memory.Used,
				status.Resources.Memory.Available,
				status.Resources.Memory.Percentage))
		}
	}

	m.podsItem.SetTitle(fmt.Sprintf("Pods: %d total", status.PodStatus.Total))

	// Update individual pod status items with visual indicators
	m.podsReadyItem.SetTitle(fmt.Sprintf("  🟢 Ready: %d", status.PodStatus.RunningReady))
	m.podsNotReadyItem.SetTitle(fmt.Sprintf("  🛑 Not Ready: %d", status.PodStatus.RunningNotReady))
	m.podsPendingItem.SetTitle(fmt.Sprintf("  ⏳ Pending: %d", status.PodStatus.Pending))
	m.podsCompletedItem.SetTitle(fmt.Sprintf("  ✅ Completed: %d", status.PodStatus.Completed))
	m.podsFailedItem.SetTitle(fmt.Sprintf("  ❌ Failed: %d", status.PodStatus.Failed))

	// Update pod submenus with individual pod names
	m.updatePodSubmenus(status.PodStatus)

	// Show/hide items based on count (optional - keeps menu clean)
	if status.PodStatus.RunningReady == 0 {
		m.podsReadyItem.Hide()
	} else {
		m.podsReadyItem.Show()
	}
	if status.PodStatus.RunningNotReady == 0 {
		m.podsNotReadyItem.Hide()
	} else {
		m.podsNotReadyItem.Show()
	}
	if status.PodStatus.Pending == 0 {
		m.podsPendingItem.Hide()
	} else {
		m.podsPendingItem.Show()
	}
	if status.PodStatus.Completed == 0 {
		m.podsCompletedItem.Hide()
	} else {
		m.podsCompletedItem.Show()
	}
	if status.PodStatus.Failed == 0 {
		m.podsFailedItem.Hide()
	} else {
		m.podsFailedItem.Show()
	}
}

// updateError updates the display when an error occurs
func (m *Manager) updateError(err error) {
	m.updateIcon(models.HealthCritical)
	systray.SetTooltip(fmt.Sprintf("K8s Tray - Error: %v", err))
	m.statusItem.SetTitle(fmt.Sprintf("Status: Error - %v", err))
}

// updateIcon updates the tray icon based on health status
func (m *Manager) updateIcon(health models.HealthStatus) {
	var iconData []byte

	switch health {
	case models.HealthHealthy:
		iconData = getGreenIcon()
	case models.HealthWarning:
		iconData = getYellowIcon()
	case models.HealthCritical:
		iconData = getRedIcon()
	default:
		iconData = getGrayIcon()
	}

	log.Printf("Setting tray icon for health status: %s", health)
	systray.SetIcon(iconData)
}

// refreshNamespaceMenu refreshes the namespace submenu
func (m *Manager) refreshNamespaceMenu(ctx context.Context) {
	namespaces, err := m.k8sClient.GetAllNamespaces(ctx)
	if err != nil {
		log.Printf("Failed to get namespaces: %v", err)
		return
	}

	// Clear existing items
	for _, item := range m.namespaceItems {
		item.Hide()
	}
	m.namespaceItems = make(map[string]*systray.MenuItem)

	// Hide existing separator if it exists
	if m.namespaceSeparator != nil {
		m.namespaceSeparator.Hide()
	}

	// Add "All Namespaces" option first
	allItem := m.namespaceMenu.AddSubMenuItem("All Namespaces", "View pods from all namespaces")
	m.namespaceItems[config.AllNamespaces] = allItem

	// Mark current selection
	if m.config.Namespace == config.AllNamespaces {
		allItem.Check()
	}

	// Handle clicks for all namespaces
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-allItem.ClickedCh:
				m.switchNamespace(ctx, config.AllNamespaces)
			}
		}
	}()

	// Add separator
	m.namespaceSeparator = m.namespaceMenu.AddSubMenuItem("─────────────", "")
	m.namespaceSeparator.Disable()

	// Add namespace items
	for _, ns := range namespaces {
		item := m.namespaceMenu.AddSubMenuItem(ns, fmt.Sprintf("Switch to namespace %s", ns))
		m.namespaceItems[ns] = item

		// Mark current selection
		if m.config.Namespace == ns {
			item.Check()
		}

		// Handle clicks
		go func(namespace string, menuItem *systray.MenuItem) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-menuItem.ClickedCh:
					m.switchNamespace(ctx, namespace)
				}
			}
		}(ns, item)
	}
}

// refreshContextMenu refreshes the context submenu
func (m *Manager) refreshContextMenu(ctx context.Context) {
	contexts, err := m.k8sClient.GetAllContexts()
	if err != nil {
		log.Printf("Failed to get contexts: %v", err)
		return
	}

	// Clear existing items
	for _, item := range m.contextItems {
		item.Hide()
	}
	m.contextItems = make(map[string]*systray.MenuItem)

	// Get current context
	currentContext, err := m.k8sClient.GetCurrentContext()
	if err != nil {
		log.Printf("Failed to get current context: %v", err)
		currentContext = ""
	}

	// Add context items
	for _, contextName := range contexts {
		item := m.contextMenu.AddSubMenuItem(contextName, fmt.Sprintf("Switch to context %s", contextName))
		m.contextItems[contextName] = item

		// Mark current selection
		if m.config.Context == contextName || (m.config.Context == "" && contextName == currentContext) {
			item.Check()
		}

		// Handle clicks
		go func(context string, menuItem *systray.MenuItem) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-menuItem.ClickedCh:
					m.switchContext(ctx, context)
				}
			}
		}(contextName, item)
	}
}

// refreshSettingsMenu refreshes the settings submenu
func (m *Manager) refreshSettingsMenu(ctx context.Context) {
	// Clear existing items
	for _, item := range m.intervalItems {
		item.Hide()
	}
	m.intervalItems = make(map[time.Duration]*systray.MenuItem)

	// Define available refresh intervals
	intervals := []struct {
		duration time.Duration
		label    string
	}{
		{5 * time.Second, "5 seconds"},
		{10 * time.Second, "10 seconds"},
		{15 * time.Second, "15 seconds"},
		{30 * time.Second, "30 seconds"},
		{1 * time.Minute, "1 minute"},
		{2 * time.Minute, "2 minutes"},
		{5 * time.Minute, "5 minutes"},
	}

	// Add refresh interval section
	m.settingsMenu.AddSubMenuItem("Refresh Interval:", "Current refresh interval setting").Disable()

	// Add interval items
	for _, interval := range intervals {
		item := m.settingsMenu.AddSubMenuItem(fmt.Sprintf("  %s", interval.label), fmt.Sprintf("Set refresh interval to %s", interval.label))
		m.intervalItems[interval.duration] = item

		// Mark current selection
		if m.config.PollInterval == interval.duration {
			item.Check()
		}

		// Handle clicks
		go func(duration time.Duration, menuItem *systray.MenuItem) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-menuItem.ClickedCh:
					m.setRefreshInterval(ctx, duration)
				}
			}
		}(interval.duration, item)
	}
}

// setRefreshInterval changes the refresh interval
func (m *Manager) setRefreshInterval(_ context.Context, interval time.Duration) {
	// Uncheck previous selection
	if prevItem, exists := m.intervalItems[m.config.PollInterval]; exists {
		prevItem.Uncheck()
	}

	// Update configuration
	m.config.PollInterval = interval

	// Check new selection
	if newItem, exists := m.intervalItems[interval]; exists {
		newItem.Check()
	}

	// Save configuration
	if err := m.config.Save(); err != nil {
		log.Printf("Failed to save config: %v", err)
	}

	// Notify monitoring goroutine about the interval change
	select {
	case m.intervalChanged <- interval:
	default:
		// Channel is full, but that's okay - we'll use the latest value
	}

	log.Printf("Changed refresh interval to: %s", interval)
}

// switchNamespace switches to a different namespace
func (m *Manager) switchNamespace(ctx context.Context, namespace string) {
	// Uncheck previous selection
	if prevItem, exists := m.namespaceItems[m.config.Namespace]; exists {
		prevItem.Uncheck()
	}

	// Update configuration
	m.config.Namespace = namespace

	// Check new selection
	if newItem, exists := m.namespaceItems[namespace]; exists {
		newItem.Check()
	}

	// Save configuration
	if err := m.config.Save(); err != nil {
		log.Printf("Failed to save config: %v", err)
	}

	// Clear pod submenus to avoid showing stale pod data from the old namespace
	m.clearPodSubmenus()

	// Refresh status
	m.refreshStatus(ctx)

	log.Printf("Switched to namespace: %s", namespace)
}

// switchContext switches to a different context
func (m *Manager) switchContext(ctx context.Context, contextName string) {
	// Uncheck previous selection
	currentContext, _ := m.k8sClient.GetCurrentContext()
	if m.config.Context == "" {
		// If no context is set in config, use the current context from kubeconfig
		if prevItem, exists := m.contextItems[currentContext]; exists {
			prevItem.Uncheck()
		}
	} else {
		if prevItem, exists := m.contextItems[m.config.Context]; exists {
			prevItem.Uncheck()
		}
	}

	// Update configuration
	m.config.Context = contextName

	// Check new selection
	if newItem, exists := m.contextItems[contextName]; exists {
		newItem.Check()
	}

	// Save configuration
	if err := m.config.Save(); err != nil {
		log.Printf("Failed to save config: %v", err)
	}

	// Need to recreate the Kubernetes client with the new context
	newClient, err := kubernetes.NewClient(m.config)
	if err != nil {
		log.Printf("Failed to create new client with context %s: %v", contextName, err)
		return
	}

	// Update the client
	m.k8sClient = newClient

	// Reset all menu items to prevent showing stale data from the old context
	m.resetMenuState()

	// Refresh status
	m.refreshStatus(ctx)

	// Refresh namespace menu since we switched clusters
	go m.refreshNamespaceMenu(ctx)

	log.Printf("Switched to context: %s", contextName)
}

// resetMenuState resets all menu items to their initial/loading state
func (m *Manager) resetMenuState() {
	// Reset main status items to loading state
	m.statusItem.SetTitle("Status: Connecting...")
	m.clusterItem.SetTitle("Cluster: Unknown")
	m.namespaceItem.SetTitle("Namespace: Loading...")

	// Reset resource items if they exist
	if m.cpuItem != nil {
		m.cpuItem.SetTitle("CPU: Loading...")
	}
	if m.memoryItem != nil {
		m.memoryItem.SetTitle("Memory: Loading...")
	}

	// Reset pod status items
	m.podsItem.SetTitle("Pods: Loading...")
	m.podsReadyItem.SetTitle("  🟢 Ready: 0")
	m.podsNotReadyItem.SetTitle("  🛑 Not Ready: 0")
	m.podsPendingItem.SetTitle("  ⏳ Pending: 0")
	m.podsCompletedItem.SetTitle("  ✅ Completed: 0")
	m.podsFailedItem.SetTitle("  ❌ Failed: 0")

	// Hide all pod status items initially
	m.podsReadyItem.Hide()
	m.podsNotReadyItem.Hide()
	m.podsPendingItem.Hide()
	m.podsCompletedItem.Hide()
	m.podsFailedItem.Hide()

	// Clear all pod submenus
	m.clearPodSubmenus()

	// Reset tooltip
	systray.SetTooltip("K8s Tray - Connecting...")

	// Reset icon to unknown state
	m.updateIcon(models.HealthUnknown)
	m.currentHealth = models.HealthUnknown

	// Clear current status
	m.currentStatus = nil
}

// showWindowsHelp displays Windows-specific help information in the log/console
func (m *Manager) showWindowsHelp() {
	log.Println("=== K8s Tray for Windows ===")
	log.Println("If the K8s Tray icon is not visible in your system tray:")
	log.Println("1. Look for the ^ arrow icon in your system tray")
	log.Println("2. Click it to see hidden icons")
	log.Println("3. Drag the K8s Tray icon from the hidden area to the visible tray")
	log.Println("4. Right-click on an empty area of the taskbar")
	log.Println("5. Select 'Taskbar settings' > 'Select which icons appear on the taskbar'")
	log.Println("6. Find 'K8s Tray' and turn it 'On'")
	log.Println("Visit: https://support.microsoft.com/en-us/windows/how-to-customize-the-taskbar-notification-area")
}

// updatePodSubmenus updates the submenu items for each pod state category
func (m *Manager) updatePodSubmenus(podStatus *models.PodStatus) {
	// Clear existing submenu items
	m.clearPodSubmenus()

	// Group pods by state
	var readyPods, notReadyPods, pendingPods, completedPods, failedPods []models.PodDetail

	for _, pod := range podStatus.Details {
		switch pod.Phase {
		case podPhaseRunning:
			if pod.Ready {
				readyPods = append(readyPods, pod)
			} else {
				notReadyPods = append(notReadyPods, pod)
			}
		case podPhasePending:
			pendingPods = append(pendingPods, pod)
		case podPhaseSucceeded:
			completedPods = append(completedPods, pod)
		case podPhaseFailed:
			failedPods = append(failedPods, pod)
		}
	}

	// Add submenu items for each category
	m.addPodSubmenuItems(m.podsReadyItem, readyPods, m.podsReadySubmenu)
	m.addPodSubmenuItems(m.podsNotReadyItem, notReadyPods, m.podsNotReadySubmenu)
	m.addPodSubmenuItems(m.podsPendingItem, pendingPods, m.podsPendingSubmenu)
	m.addPodSubmenuItems(m.podsCompletedItem, completedPods, m.podsCompletedSubmenu)
	m.addPodSubmenuItems(m.podsFailedItem, failedPods, m.podsFailedSubmenu)
}

// clearPodSubmenus clears all existing pod submenu items
func (m *Manager) clearPodSubmenus() {
	// Clear ready pods submenu
	for _, item := range m.podsReadySubmenu {
		item.Hide()
	}
	m.podsReadySubmenu = make(map[string]*systray.MenuItem)

	// Clear not ready pods submenu
	for _, item := range m.podsNotReadySubmenu {
		item.Hide()
	}
	m.podsNotReadySubmenu = make(map[string]*systray.MenuItem)

	// Clear pending pods submenu
	for _, item := range m.podsPendingSubmenu {
		item.Hide()
	}
	m.podsPendingSubmenu = make(map[string]*systray.MenuItem)

	// Clear completed pods submenu
	for _, item := range m.podsCompletedSubmenu {
		item.Hide()
	}
	m.podsCompletedSubmenu = make(map[string]*systray.MenuItem)

	// Clear failed pods submenu
	for _, item := range m.podsFailedSubmenu {
		item.Hide()
	}
	m.podsFailedSubmenu = make(map[string]*systray.MenuItem)
}

// addPodSubmenuItems adds submenu items for pods in a specific state
func (m *Manager) addPodSubmenuItems(parentItem *systray.MenuItem, pods []models.PodDetail, submenuMap map[string]*systray.MenuItem) {
	if len(pods) == 0 {
		return
	}

	for _, pod := range pods {
		// Create display name with namespace if not "all namespaces" view
		displayName := pod.Name
		if m.config.Namespace == config.AllNamespaces {
			displayName = fmt.Sprintf("%s (%s)", pod.Name, pod.Namespace)
		}

		// Create tooltip with additional pod information
		tooltip := fmt.Sprintf("Pod: %s\nNamespace: %s\nPhase: %s\nReady: %t",
			pod.Name, pod.Namespace, pod.Phase, pod.Ready)
		if pod.Restarts > 0 {
			tooltip += fmt.Sprintf("\nRestarts: %d", pod.Restarts)
		}
		tooltip += fmt.Sprintf("\nAge: %s", pod.Age.Truncate(time.Second))

		// Add submenu item
		item := parentItem.AddSubMenuItem(displayName, tooltip)
		item.Disable() // Make it non-clickable for now, just informational

		// Store in the submenu map using a unique key
		key := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
		submenuMap[key] = item
	}
}
