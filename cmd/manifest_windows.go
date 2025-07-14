//go:build windows

package main

import (
	_ "embed"
)

// Windows application manifest for proper system tray integration
//
//go:embed k8s-tray.exe.manifest
var windowsManifest []byte
