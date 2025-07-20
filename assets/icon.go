package assets

import (
	_ "embed"
	"log"
)

// Windows resource file - this will be automatically included when building for Windows
// The .syso file is automatically linked by the Go toolchain when present in the same package

//go:embed icon-256.png
var iconPNG []byte

//go:embed icon.ico
var iconICO []byte

//go:embed icon-512.png
var iconLarge []byte

// GetIconData returns the appropriate icon data for the current platform
func GetIconData() []byte {
	if len(iconPNG) == 0 {
		log.Println("Warning: icon PNG data not embedded")
		return nil
	}
	return iconPNG
}

// GetIconICO returns Windows ICO format icon
func GetIconICO() []byte {
	if len(iconICO) == 0 {
		log.Println("Warning: icon ICO data not embedded")
		return nil
	}
	return iconICO
}

// GetIconLarge returns large icon for high-DPI displays
func GetIconLarge() []byte {
	if len(iconLarge) == 0 {
		log.Println("Warning: large icon data not embedded")
		return nil
	}
	return iconLarge
}
