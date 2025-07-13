package tray

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// createSimpleIcon creates a simple colored square icon
func createSimpleIcon(r, g, b uint8) []byte {
	const size = 16
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Fill with transparent background
	transparent := color.RGBA{0, 0, 0, 0}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, transparent)
		}
	}

	// Create a simple circle/dot in the center
	centerX, centerY := size/2, size/2
	radius := size/4

	iconColor := color.RGBA{r, g, b, 255}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, iconColor)
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

// getGreenIcon returns a green circle icon
func getGreenIcon() []byte {
	return createSimpleIcon(0, 255, 0) // Green
}

// getYellowIcon returns a yellow circle icon
func getYellowIcon() []byte {
	return createSimpleIcon(255, 255, 0) // Yellow
}

// getRedIcon returns a red circle icon
func getRedIcon() []byte {
	return createSimpleIcon(255, 0, 0) // Red
}

// getGrayIcon returns a gray circle icon
func getGrayIcon() []byte {
	return createSimpleIcon(128, 128, 128) // Gray
}
