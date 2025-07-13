package kubernetes

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/k8s-tray/k8s-tray/internal/config"
	"github.com/k8s-tray/k8s-tray/pkg/models"
)

// Client wraps the Kubernetes client with additional functionality
type Client struct {
	clientset *kubernetes.Clientset
	config    *config.Config
	namespace string
}

// NewClient creates a new Kubernetes client
func NewClient(cfg *config.Config) (*Client, error) {
	// Build config from kubeconfig
	config, err := buildConfig(cfg.KubeConfig, cfg.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{
		clientset: clientset,
		config:    cfg,
		namespace: cfg.Namespace,
	}, nil
}

// buildConfig builds the Kubernetes configuration
func buildConfig(kubeconfig, context string) (*rest.Config, error) {
	// Try in-cluster config first
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	// Use kubeconfig
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: context},
	).ClientConfig()
}

// GetClusterStatus returns the overall cluster status
func (c *Client) GetClusterStatus(ctx context.Context) (*models.ClusterStatus, error) {
	// Get server version
	version, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	// Get current context
	currentContext, err := c.GetCurrentContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get current context: %w", err)
	}

	// Get pod status
	podStatus, err := c.GetPodStatus(ctx, c.config.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod status: %w", err)
	}

	return &models.ClusterStatus{
		ClusterName:   currentContext,
		ServerVersion: version.String(),
		PodStatus:     podStatus,
		LastUpdated:   time.Now(),
		HealthStatus:  calculateHealthStatus(podStatus),
	}, nil
}

// GetPodStatus returns pod status for the specified namespace
func (c *Client) GetPodStatus(ctx context.Context, namespace string) (*models.PodStatus, error) {
	// Determine which namespace to query
	var queryNamespace string
	if namespace == config.AllNamespaces {
		queryNamespace = "" // Empty string means all namespaces
	} else {
		queryNamespace = namespace
	}

	// List pods in namespace
	pods, err := c.clientset.CoreV1().Pods(queryNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	status := &models.PodStatus{
		Total:           len(pods.Items),
		Running:         0,
		RunningReady:    0,
		RunningNotReady: 0,
		Pending:         0,
		Failed:          0,
		Unknown:         0,
		Completed:       0,
		Details:         make([]models.PodDetail, 0, len(pods.Items)),
	}

	// Process each pod
	for _, pod := range pods.Items {
		detail := models.PodDetail{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Phase:     string(pod.Status.Phase),
			Ready:     isPodReady(&pod),
			Restarts:  getRestartCount(&pod),
			Age:       time.Since(pod.CreationTimestamp.Time),
		}

		status.Details = append(status.Details, detail)

		// Update counters
		switch pod.Status.Phase {
		case corev1.PodRunning:
			status.Running++
			if isPodReady(&pod) {
				status.RunningReady++
			} else {
				status.RunningNotReady++
			}
		case corev1.PodPending:
			status.Pending++
		case corev1.PodSucceeded:
			status.Completed++
		case corev1.PodFailed:
			status.Failed++
		default:
			status.Unknown++
		}
	}

	return status, nil
}

// GetAllNamespaces returns all namespaces in the cluster
func (c *Client) GetAllNamespaces(ctx context.Context) ([]string, error) {
	namespaces, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	names := make([]string, len(namespaces.Items))
	for i, ns := range namespaces.Items {
		names[i] = ns.Name
	}

	return names, nil
}

// GetEvents returns recent events in the namespace
func (c *Client) GetEvents(ctx context.Context, namespace string) ([]models.Event, error) {
	// Determine which namespace to query
	var queryNamespace string
	if namespace == config.AllNamespaces {
		queryNamespace = "" // Empty string means all namespaces
	} else {
		queryNamespace = namespace
	}

	events, err := c.clientset.CoreV1().Events(queryNamespace).List(ctx, metav1.ListOptions{
		Limit: 50,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	result := make([]models.Event, len(events.Items))
	for i, event := range events.Items {
		result[i] = models.Event{
			Type:      event.Type,
			Reason:    event.Reason,
			Message:   event.Message,
			Object:    event.InvolvedObject.Name,
			Timestamp: event.LastTimestamp.Time,
		}
	}

	return result, nil
}

// TestConnection tests the connection to the Kubernetes cluster
func (c *Client) TestConnection(ctx context.Context) error {
	_, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}
	return nil
}

// Helper functions

// GetCurrentContext returns the current context name
func (c *Client) GetCurrentContext() (string, error) {
	config, err := clientcmd.LoadFromFile(c.config.KubeConfig)
	if err != nil {
		return "", err
	}

	if c.config.Context != "" {
		return c.config.Context, nil
	}

	return config.CurrentContext, nil
}

// GetAllContexts returns all available contexts from the kubeconfig
func (c *Client) GetAllContexts() ([]string, error) {
	config, err := clientcmd.LoadFromFile(c.config.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	contexts := make([]string, 0, len(config.Contexts))
	for contextName := range config.Contexts {
		contexts = append(contexts, contextName)
	}

	return contexts, nil
}

// isPodReady checks if a pod is ready
func isPodReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

// getRestartCount returns the total restart count for a pod
func getRestartCount(pod *corev1.Pod) int32 {
	var restarts int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restarts += containerStatus.RestartCount
	}
	return restarts
}

// calculateHealthStatus determines the overall health status
func calculateHealthStatus(podStatus *models.PodStatus) models.HealthStatus {
	if podStatus.Failed > 0 {
		return models.HealthCritical
	}
	if podStatus.Pending > 0 || podStatus.Unknown > 0 || podStatus.RunningNotReady > 0 {
		return models.HealthWarning
	}
	if podStatus.RunningReady > 0 {
		return models.HealthHealthy
	}
	return models.HealthUnknown
}
