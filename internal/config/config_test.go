package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".k8s-tray.yaml")

	// Override getConfigPath for testing
	originalGetConfigPath := getConfigPath
	defer func() {
		getConfigPath = originalGetConfigPath
	}()

	getConfigPath = func() string {
		return configPath
	}

	// Test loading default config when file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify default values
	if cfg.PollInterval != 15*time.Second {
		t.Errorf("Expected poll interval 15s, got %v", cfg.PollInterval)
	}

	if cfg.Namespace != AllNamespaces {
		t.Errorf("Expected namespace '%s', got %s", AllNamespaces, cfg.Namespace)
	}

	if !cfg.ShowNotifications {
		t.Error("Expected show_notifications to be true")
	}
}

func TestConfigValidation(t *testing.T) {
	cfg := &Config{
		PollInterval: 500 * time.Millisecond, // Too short
	}

	cfg.validate()

	// Should be adjusted to minimum
	if cfg.PollInterval < time.Second {
		t.Errorf("Poll interval should be adjusted to minimum 1s, got %v", cfg.PollInterval)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".k8s-tray.yaml")

	// Override getConfigPath for testing
	originalGetConfigPath := getConfigPath
	defer func() {
		getConfigPath = originalGetConfigPath
	}()

	getConfigPath = func() string {
		return configPath
	}

	// Create a config with custom values
	cfg := &Config{
		Namespace:         "test-namespace",
		PollInterval:      10 * time.Second,
		ShowNotifications: false,
		Theme:             "dark",
	}

	// Save config
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values
	if loadedCfg.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got %s", loadedCfg.Namespace)
	}

	if loadedCfg.PollInterval != 10*time.Second {
		t.Errorf("Expected poll interval 10s, got %v", loadedCfg.PollInterval)
	}

	if loadedCfg.ShowNotifications {
		t.Error("Expected show_notifications to be false")
	}

	if loadedCfg.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got %s", loadedCfg.Theme)
	}
}

func TestGetDefaultKubeConfig(t *testing.T) {
	// Test with KUBECONFIG env var
	original := os.Getenv("KUBECONFIG")
	defer os.Setenv("KUBECONFIG", original)

	testPath := "/test/kubeconfig"
	os.Setenv("KUBECONFIG", testPath)

	result := getDefaultKubeConfig()
	if result != testPath {
		t.Errorf("Expected %s, got %s", testPath, result)
	}

	// Test without KUBECONFIG env var
	os.Unsetenv("KUBECONFIG")

	result = getDefaultKubeConfig()
	if result == "" {
		t.Error("Default kubeconfig path should not be empty")
	}
}
