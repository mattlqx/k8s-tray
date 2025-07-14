package tray

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"runtime"
)

// createSimpleIcon creates a simple colored square icon
func createSimpleIcon(r, g, b uint8) []byte {
	if runtime.GOOS == "windows" {
		return createICOIcon(r, g, b)
	}
	return createPNGIcon(r, g, b)
}

// createPNGIcon creates a PNG format icon
func createPNGIcon(r, g, b uint8) []byte {
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
	radius := size / 4

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
	if err := png.Encode(&buf, img); err != nil {
		// If encoding fails, return empty byte slice
		return []byte{}
	}
	return buf.Bytes()
}

// createICOIcon creates an ICO format icon for Windows
func createICOIcon(r, g, b uint8) []byte {
	const size = 16

	// Create the image data
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
	radius := size / 4

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

	// Convert to ICO format
	return createICOFromImage(img)
}

// createICOFromImage converts an image to ICO format
func createICOFromImage(img *image.RGBA) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a buffer for the ICO file
	var buf bytes.Buffer

	// ICO header (6 bytes)
	if err := binary.Write(&buf, binary.LittleEndian, uint16(0)); err != nil {
		return []byte{}
	} // Reserved (must be 0)
	if err := binary.Write(&buf, binary.LittleEndian, uint16(1)); err != nil {
		return []byte{}
	} // Type (1 = ICO)
	if err := binary.Write(&buf, binary.LittleEndian, uint16(1)); err != nil {
		return []byte{}
	} // Number of images

	// ICO directory entry (16 bytes)
	buf.WriteByte(byte(width))  // Width (0 = 256)
	buf.WriteByte(byte(height)) // Height (0 = 256)
	buf.WriteByte(0)            // Color count (0 = >256 colors)
	buf.WriteByte(0)            // Reserved
	if err := binary.Write(&buf, binary.LittleEndian, uint16(1)); err != nil {
		return []byte{}
	} // Color planes
	if err := binary.Write(&buf, binary.LittleEndian, uint16(32)); err != nil {
		return []byte{}
	} // Bits per pixel

	// Create the bitmap data
	bitmapData := createBitmapData(img)
	// Check for potential overflow when converting to uint32
	if len(bitmapData) > 4294967295 {
		return []byte{}
	}
	// #nosec G115 -- Safe conversion for small icon dimensions
	if err := binary.Write(&buf, binary.LittleEndian, uint32(len(bitmapData))); err != nil {
		return []byte{}
	} // Image size
	if err := binary.Write(&buf, binary.LittleEndian, uint32(22)); err != nil {
		return []byte{}
	} // Offset to image data

	// Append the bitmap data
	buf.Write(bitmapData)

	return buf.Bytes()
}

// createBitmapData creates the bitmap data for the ICO
func createBitmapData(img *image.RGBA) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var buf bytes.Buffer

	// Helper function to handle binary.Write errors
	writeOrReturn := func(data interface{}) bool {
		if err := binary.Write(&buf, binary.LittleEndian, data); err != nil {
			return false
		}
		return true
	}

	// Bitmap info header (40 bytes)
	// Check for potential integer overflow before converting
	if width < 0 || width > 2147483647 || height < 0 || height > 2147483647 {
		return []byte{}
	}

	// #nosec G115 -- Safe conversion for small icon dimensions
	if !writeOrReturn(uint32(40)) || // Header size
		!writeOrReturn(int32(width)) || // Width
		!writeOrReturn(int32(height*2)) || // Height (doubled for ICO)
		!writeOrReturn(uint16(1)) || // Planes
		!writeOrReturn(uint16(32)) || // Bits per pixel
		!writeOrReturn(uint32(0)) || // Compression
		!writeOrReturn(uint32(0)) || // Image size
		!writeOrReturn(int32(0)) || // X pixels per meter
		!writeOrReturn(int32(0)) || // Y pixels per meter
		!writeOrReturn(uint32(0)) || // Colors used
		!writeOrReturn(uint32(0)) { // Important colors
		return []byte{}
	}

	// Pixel data (BGRA format, bottom-up)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			c := img.RGBAAt(x, y)
			buf.WriteByte(c.B) // Blue
			buf.WriteByte(c.G) // Green
			buf.WriteByte(c.R) // Red
			buf.WriteByte(c.A) // Alpha
		}
	}

	// AND mask (1 bit per pixel, padded to 32-bit boundary)
	maskBytesPerRow := (width + 31) / 32 * 4
	for y := 0; y < height; y++ {
		for x := 0; x < maskBytesPerRow; x++ {
			buf.WriteByte(0) // No mask (all pixels visible)
		}
	}

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
