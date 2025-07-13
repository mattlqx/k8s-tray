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

	// Convert to ICO format
	return createICOFromImage(img)
}

// createICOFromImage converts an image to ICO format
func createICOFromImage(img *image.RGBA) []byte {
	const size = 16
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a buffer for the ICO file
	var buf bytes.Buffer

	// ICO header (6 bytes)
	binary.Write(&buf, binary.LittleEndian, uint16(0))    // Reserved (must be 0)
	binary.Write(&buf, binary.LittleEndian, uint16(1))    // Type (1 = ICO)
	binary.Write(&buf, binary.LittleEndian, uint16(1))    // Number of images

	// ICO directory entry (16 bytes)
	buf.WriteByte(byte(width))                            // Width (0 = 256)
	buf.WriteByte(byte(height))                           // Height (0 = 256)
	buf.WriteByte(0)                                      // Color count (0 = >256 colors)
	buf.WriteByte(0)                                      // Reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1))   // Color planes
	binary.Write(&buf, binary.LittleEndian, uint16(32))  // Bits per pixel

	// Create the bitmap data
	bitmapData := createBitmapData(img)
	binary.Write(&buf, binary.LittleEndian, uint32(len(bitmapData))) // Image size
	binary.Write(&buf, binary.LittleEndian, uint32(22))              // Offset to image data

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

	// Bitmap info header (40 bytes)
	binary.Write(&buf, binary.LittleEndian, uint32(40))           // Header size
	binary.Write(&buf, binary.LittleEndian, int32(width))        // Width
	binary.Write(&buf, binary.LittleEndian, int32(height*2))     // Height (doubled for ICO)
	binary.Write(&buf, binary.LittleEndian, uint16(1))           // Planes
	binary.Write(&buf, binary.LittleEndian, uint16(32))          // Bits per pixel
	binary.Write(&buf, binary.LittleEndian, uint32(0))           // Compression
	binary.Write(&buf, binary.LittleEndian, uint32(0))           // Image size
	binary.Write(&buf, binary.LittleEndian, int32(0))            // X pixels per meter
	binary.Write(&buf, binary.LittleEndian, int32(0))            // Y pixels per meter
	binary.Write(&buf, binary.LittleEndian, uint32(0))           // Colors used
	binary.Write(&buf, binary.LittleEndian, uint32(0))           // Important colors

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
