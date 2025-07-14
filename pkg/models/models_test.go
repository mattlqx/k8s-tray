package models

import (
	"testing"
	"time"
)

func TestHealthStatus_String(t *testing.T) {
	tests := []struct {
		status   HealthStatus
		expected string
	}{
		{HealthHealthy, "Healthy"},
		{HealthWarning, "Warning"},
		{HealthCritical, "Critical"},
		{HealthUnknown, "Unknown"},
	}

	for _, test := range tests {
		result := test.status.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestMenuAction_String(t *testing.T) {
	tests := []struct {
		action   MenuAction
		expected string
	}{
		{ActionRefresh, "Refresh"},
		{ActionSwitchNamespace, "Switch Namespace"},
		{ActionSwitchContext, "Switch Context"},
		{ActionViewLogs, "View Logs"},
		{ActionViewEvents, "View Events"},
		{ActionSettings, "Settings"},
		{ActionQuit, "Quit"},
	}

	for _, test := range tests {
		result := test.action.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestClusterStatus(t *testing.T) {
	podStatus := &PodStatus{
		Total:   5,
		Running: 3,
		Pending: 1,
		Failed:  1,
		Unknown: 0,
	}

	status := &ClusterStatus{
		ClusterName:   "test-cluster",
		ServerVersion: "v1.28.0",
		PodStatus:     podStatus,
		LastUpdated:   time.Now(),
		HealthStatus:  HealthWarning,
	}

	if status.ClusterName != "test-cluster" {
		t.Errorf("Expected cluster name 'test-cluster', got %s", status.ClusterName)
	}

	if status.PodStatus.Total != 5 {
		t.Errorf("Expected total pods 5, got %d", status.PodStatus.Total)
	}

	if status.HealthStatus != HealthWarning {
		t.Errorf("Expected health status Warning, got %s", status.HealthStatus.String())
	}
}

func TestPodDetail(t *testing.T) {
	detail := PodDetail{
		Name:      "test-pod",
		Namespace: "default",
		Phase:     "Running",
		Ready:     true,
		Restarts:  0,
		Age:       5 * time.Minute,
	}

	if detail.Name != "test-pod" {
		t.Errorf("Expected pod name 'test-pod', got %s", detail.Name)
	}

	if !detail.Ready {
		t.Error("Expected pod to be ready")
	}

	if detail.Restarts != 0 {
		t.Errorf("Expected 0 restarts, got %d", detail.Restarts)
	}
}

func TestEvent(t *testing.T) {
	event := Event{
		Type:      "Normal",
		Reason:    "Started",
		Message:   "Container started successfully",
		Object:    "test-pod",
		Timestamp: time.Now(),
	}

	if event.Type != "Normal" {
		t.Errorf("Expected event type 'Normal', got %s", event.Type)
	}

	if event.Reason != "Started" {
		t.Errorf("Expected event reason 'Started', got %s", event.Reason)
	}

	if event.Object != "test-pod" {
		t.Errorf("Expected event object 'test-pod', got %s", event.Object)
	}
}

func TestResourceStats(t *testing.T) {
	cpuStat := &ResourceStat{
		Used:       2.5,
		Available:  4.0,
		Percentage: 62.5,
	}

	memoryStat := &ResourceStat{
		Used:       8.0,
		Available:  16.0,
		Percentage: 50.0,
	}

	resourceStats := &ResourceStats{
		CPU:    cpuStat,
		Memory: memoryStat,
	}

	if resourceStats.CPU.Used != 2.5 {
		t.Errorf("Expected CPU used 2.5, got %f", resourceStats.CPU.Used)
	}

	if resourceStats.CPU.Available != 4.0 {
		t.Errorf("Expected CPU available 4.0, got %f", resourceStats.CPU.Available)
	}

	if resourceStats.CPU.Percentage != 62.5 {
		t.Errorf("Expected CPU percentage 62.5, got %f", resourceStats.CPU.Percentage)
	}

	if resourceStats.Memory.Used != 8.0 {
		t.Errorf("Expected Memory used 8.0, got %f", resourceStats.Memory.Used)
	}

	if resourceStats.Memory.Available != 16.0 {
		t.Errorf("Expected Memory available 16.0, got %f", resourceStats.Memory.Available)
	}

	if resourceStats.Memory.Percentage != 50.0 {
		t.Errorf("Expected Memory percentage 50.0, got %f", resourceStats.Memory.Percentage)
	}
}

func TestClusterStatusWithResources(t *testing.T) {
	podStatus := &PodStatus{
		Total:   3,
		Running: 2,
		Pending: 1,
	}

	resourceStats := &ResourceStats{
		CPU: &ResourceStat{
			Used:       1.5,
			Available:  4.0,
			Percentage: 37.5,
		},
		Memory: &ResourceStat{
			Used:       4.0,
			Available:  8.0,
			Percentage: 50.0,
		},
	}

	status := &ClusterStatus{
		ClusterName:   "test-cluster",
		ServerVersion: "v1.28.0",
		PodStatus:     podStatus,
		Resources:     resourceStats,
		LastUpdated:   time.Now(),
		HealthStatus:  HealthHealthy,
	}

	if status.Resources == nil {
		t.Error("Expected resources to be set")
	}

	if status.Resources.CPU == nil {
		t.Error("Expected CPU stats to be set")
	}

	if status.Resources.Memory == nil {
		t.Error("Expected Memory stats to be set")
	}

	if status.Resources.CPU.Percentage != 37.5 {
		t.Errorf("Expected CPU percentage 37.5, got %f", status.Resources.CPU.Percentage)
	}
}
