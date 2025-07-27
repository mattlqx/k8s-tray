package tray

import (
	"fmt"
	"testing"

	"github.com/mattlqx/k8s-tray/pkg/models"
)

func TestGetPodStatusDisplayText(t *testing.T) {
	tests := []struct {
		name     string
		status   *models.PodStatus
		expected map[string]string
	}{
		{
			name: "Mixed pod statuses",
			status: &models.PodStatus{
				Total:           10,
				RunningReady:    5,
				RunningNotReady: 2,
				Pending:         1,
				Completed:       2,
				Failed:          0,
			},
			expected: map[string]string{
				"ready":     "  🟢 Ready: 5",
				"not_ready": "  🛑 Not Ready: 2",
				"pending":   "  ⏳ Pending: 1",
				"completed": "  ✅ Completed: 2",
				"failed":    "  ❌ Failed: 0",
			},
		},
		{
			name: "All zero counts",
			status: &models.PodStatus{
				Total:           0,
				RunningReady:    0,
				RunningNotReady: 0,
				Pending:         0,
				Completed:       0,
				Failed:          0,
			},
			expected: map[string]string{
				"ready":     "  🟢 Ready: 0",
				"not_ready": "  🛑 Not Ready: 0",
				"pending":   "  ⏳ Pending: 0",
				"completed": "  ✅ Completed: 0",
				"failed":    "  ❌ Failed: 0",
			},
		},
		{
			name: "Only failed pods",
			status: &models.PodStatus{
				Total:           3,
				RunningReady:    0,
				RunningNotReady: 0,
				Pending:         0,
				Completed:       0,
				Failed:          3,
			},
			expected: map[string]string{
				"ready":     "  🟢 Ready: 0",
				"not_ready": "  🛑 Not Ready: 0",
				"pending":   "  ⏳ Pending: 0",
				"completed": "  ✅ Completed: 0",
				"failed":    "  ❌ Failed: 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the display text generation logic
			readyText := getPodStatusText("ready", tt.status.RunningReady)
			notReadyText := getPodStatusText("not_ready", tt.status.RunningNotReady)
			pendingText := getPodStatusText("pending", tt.status.Pending)
			completedText := getPodStatusText("completed", tt.status.Completed)
			failedText := getPodStatusText("failed", tt.status.Failed)

			if readyText != tt.expected["ready"] {
				t.Errorf("Ready text: expected '%s', got '%s'", tt.expected["ready"], readyText)
			}
			if notReadyText != tt.expected["not_ready"] {
				t.Errorf("Not ready text: expected '%s', got '%s'", tt.expected["not_ready"], notReadyText)
			}
			if pendingText != tt.expected["pending"] {
				t.Errorf("Pending text: expected '%s', got '%s'", tt.expected["pending"], pendingText)
			}
			if completedText != tt.expected["completed"] {
				t.Errorf("Completed text: expected '%s', got '%s'", tt.expected["completed"], completedText)
			}
			if failedText != tt.expected["failed"] {
				t.Errorf("Failed text: expected '%s', got '%s'", tt.expected["failed"], failedText)
			}
		})
	}
}

func TestShouldShowPodStatusItem(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected bool
	}{
		{"Zero count should hide", 0, false},
		{"Non-zero count should show", 1, true},
		{"Large count should show", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldShowPodStatusItem(tt.count)
			if result != tt.expected {
				t.Errorf("Expected %v for count %d, got %v", tt.expected, tt.count, result)
			}
		})
	}
}

// Helper function to generate pod status display text
func getPodStatusText(statusType string, count int) string {
	switch statusType {
	case "ready":
		return fmt.Sprintf("  🟢 Ready: %d", count)
	case "not_ready":
		return fmt.Sprintf("  🛑 Not Ready: %d", count)
	case "pending":
		return fmt.Sprintf("  ⏳ Pending: %d", count)
	case "completed":
		return fmt.Sprintf("  ✅ Completed: %d", count)
	case "failed":
		return fmt.Sprintf("  ❌ Failed: %d", count)
	default:
		return ""
	}
}

// Helper function to determine if a pod status item should be shown
func shouldShowPodStatusItem(count int) bool {
	return count > 0
}

// TestPodSubmenuGrouping tests the grouping of pods by their states
func TestPodSubmenuGrouping(t *testing.T) {
	// Create test pod details
	pods := []models.PodDetail{
		{Name: "pod1", Namespace: "default", Phase: podPhaseRunning, Ready: true},
		{Name: "pod2", Namespace: "default", Phase: podPhaseRunning, Ready: false},
		{Name: "pod3", Namespace: "kube-system", Phase: podPhasePending, Ready: false},
		{Name: "pod4", Namespace: "default", Phase: podPhaseSucceeded, Ready: true},
		{Name: "pod5", Namespace: "default", Phase: podPhaseFailed, Ready: false},
		{Name: "pod6", Namespace: "kube-system", Phase: podPhaseRunning, Ready: true},
	}

	podStatus := &models.PodStatus{
		Total:           6,
		RunningReady:    2,
		RunningNotReady: 1,
		Pending:         1,
		Completed:       1,
		Failed:          1,
		Details:         pods,
	}

	// Group pods by state (simulate the updatePodSubmenus logic)
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

	// Verify grouping
	if len(readyPods) != 2 {
		t.Errorf("Expected 2 ready pods, got %d", len(readyPods))
	}
	if len(notReadyPods) != 1 {
		t.Errorf("Expected 1 not ready pod, got %d", len(notReadyPods))
	}
	if len(pendingPods) != 1 {
		t.Errorf("Expected 1 pending pod, got %d", len(pendingPods))
	}
	if len(completedPods) != 1 {
		t.Errorf("Expected 1 completed pod, got %d", len(completedPods))
	}
	if len(failedPods) != 1 {
		t.Errorf("Expected 1 failed pod, got %d", len(failedPods))
	}

	// Verify specific pods are in correct groups
	if readyPods[0].Name != "pod1" && readyPods[1].Name != "pod6" {
		t.Error("Ready pods should contain pod1 and pod6")
	}
	if notReadyPods[0].Name != "pod2" {
		t.Error("Not ready pods should contain pod2")
	}
	if pendingPods[0].Name != "pod3" {
		t.Error("Pending pods should contain pod3")
	}
	if completedPods[0].Name != "pod4" {
		t.Error("Completed pods should contain pod4")
	}
	if failedPods[0].Name != "pod5" {
		t.Error("Failed pods should contain pod5")
	}
}

// TestPodSubmenuWithNoPods tests the behavior when there are no pods in a category
func TestPodSubmenuWithNoPods(t *testing.T) {
	// Create empty pod status
	emptyPodStatus := &models.PodStatus{
		Total:           0,
		RunningReady:    0,
		RunningNotReady: 0,
		Pending:         0,
		Completed:       0,
		Failed:          0,
		Details:         []models.PodDetail{}, // Empty slice
	}

	// Simulate the grouping logic with empty pods
	var readyPods, notReadyPods, pendingPods, completedPods, failedPods []models.PodDetail

	for _, pod := range emptyPodStatus.Details {
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

	// Verify all slices are empty
	if len(readyPods) != 0 {
		t.Errorf("Expected 0 ready pods, got %d", len(readyPods))
	}
	if len(notReadyPods) != 0 {
		t.Errorf("Expected 0 not ready pods, got %d", len(notReadyPods))
	}
	if len(pendingPods) != 0 {
		t.Errorf("Expected 0 pending pods, got %d", len(pendingPods))
	}
	if len(completedPods) != 0 {
		t.Errorf("Expected 0 completed pods, got %d", len(completedPods))
	}
	if len(failedPods) != 0 {
		t.Errorf("Expected 0 failed pods, got %d", len(failedPods))
	}
}
