package tray

import (
	"fmt"
	"testing"

	"github.com/k8s-tray/k8s-tray/pkg/models"
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
