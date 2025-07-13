package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// Kubernetes configuration
	KubeConfig      string `yaml:"kubeconfig"`
	Context         string `yaml:"context"`
	Namespace       string `yaml:"namespace"`

	// Polling configuration
	PollInterval    time.Duration `yaml:"poll_interval"`

	// UI configuration
	ShowNotifications bool `yaml:"show_notifications"`
	Theme            string `yaml:"theme"`

	// Feature flags
	ShowMetrics      bool `yaml:"show_metrics"`
	ShowLogs         bool `yaml:"show_logs"`
	ShowEvents       bool `yaml:"show_events"`
}

// Constants for namespace selection
const (
	AllNamespaces = "<all>"
)

// Default configuration values
var defaultConfig = Config{
	KubeConfig:        getDefaultKubeConfig(),
	Context:          "",
	Namespace:        AllNamespaces,
	PollInterval:     15 * time.Second,
	ShowNotifications: true,
	Theme:            "auto",
	ShowMetrics:      true,
	ShowLogs:         false,
	ShowEvents:       true,
}

// Load loads the configuration from file or returns default configuration
func Load() (*Config, error) {
	cfg := defaultConfig

	// Try to load from config file
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal configuration to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0644)
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.PollInterval < time.Second {
		c.PollInterval = time.Second
	}

	if c.PollInterval > 5*time.Minute {
		c.PollInterval = 5 * time.Minute
	}

	return nil
}

// getConfigPath returns the path to the configuration file
var getConfigPath = func() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".k8s-tray.yaml")
	}

	return filepath.Join(homeDir, ".k8s-tray.yaml")
}

// getDefaultKubeConfig returns the default kubeconfig path
func getDefaultKubeConfig() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".kube", "config")
}
