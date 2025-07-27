package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/systray"
	"github.com/mattlqx/k8s-tray/internal/config"
	"github.com/mattlqx/k8s-tray/internal/kubernetes"
	"github.com/mattlqx/k8s-tray/internal/tray"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Kubernetes client
	k8sClient, err := kubernetes.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Create tray manager
	trayManager := tray.NewManager(k8sClient, cfg)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Received shutdown signal, cleaning up...")
		cancel()
		systray.Quit()
	}()

	// Debug: Check if we're running in a proper GUI environment
	log.Printf("Starting k8s-tray...")
	log.Printf("PID: %d", os.Getpid())

	// Start the system tray
	log.Printf("Initializing system tray...")
	systray.Run(func() {
		log.Printf("System tray ready, initializing manager...")
		trayManager.OnReady(ctx)
	}, func() {
		log.Printf("System tray exiting...")
		trayManager.OnExit()
	})

	log.Printf("Application exiting...")
}
