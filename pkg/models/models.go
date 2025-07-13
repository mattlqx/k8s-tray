package models

import (
	"time"
)

// HealthStatus represents the health status of the cluster
type HealthStatus int

const (
	HealthUnknown HealthStatus = iota
	HealthHealthy
	HealthWarning
	HealthCritical
)

// String returns the string representation of the health status
func (h HealthStatus) String() string {
	switch h {
	case HealthHealthy:
		return "Healthy"
	case HealthWarning:
		return "Warning"
	case HealthCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// ClusterStatus represents the overall cluster status
type ClusterStatus struct {
	ClusterName   string       `json:"cluster_name"`
	ServerVersion string       `json:"server_version"`
	PodStatus     *PodStatus   `json:"pod_status"`
	LastUpdated   time.Time    `json:"last_updated"`
	HealthStatus  HealthStatus `json:"health_status"`
}

// PodStatus represents the status of pods in a namespace
type PodStatus struct {
	Total           int         `json:"total"`
	Running         int         `json:"running"`
	RunningReady    int         `json:"running_ready"`
	RunningNotReady int         `json:"running_not_ready"`
	Pending         int         `json:"pending"`
	Failed          int         `json:"failed"`
	Unknown         int         `json:"unknown"`
	Completed       int         `json:"completed"`
	Details         []PodDetail `json:"details"`
}

// PodDetail represents detailed information about a pod
type PodDetail struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Phase     string        `json:"phase"`
	Ready     bool          `json:"ready"`
	Restarts  int32         `json:"restarts"`
	Age       time.Duration `json:"age"`
}

// Event represents a Kubernetes event
type Event struct {
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	Message   string    `json:"message"`
	Object    string    `json:"object"`
	Timestamp time.Time `json:"timestamp"`
}

// TrayState represents the current state of the system tray
type TrayState struct {
	Status     HealthStatus `json:"status"`
	LastUpdate time.Time    `json:"last_update"`
	Error      string       `json:"error,omitempty"`
}

// MenuAction represents an action that can be performed from the menu
type MenuAction int

const (
	ActionRefresh MenuAction = iota
	ActionSwitchNamespace
	ActionSwitchContext
	ActionViewLogs
	ActionViewEvents
	ActionSettings
	ActionQuit
)

// String returns the string representation of the menu action
func (a MenuAction) String() string {
	switch a {
	case ActionRefresh:
		return "Refresh"
	case ActionSwitchNamespace:
		return "Switch Namespace"
	case ActionSwitchContext:
		return "Switch Context"
	case ActionViewLogs:
		return "View Logs"
	case ActionViewEvents:
		return "View Events"
	case ActionSettings:
		return "Settings"
	case ActionQuit:
		return "Quit"
	default:
		return "Unknown"
	}
}
