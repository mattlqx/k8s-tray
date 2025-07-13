package tray

import (
	"runtime"
	"testing"
)

func TestCreateSimpleIcon(t *testing.T) {
	// Test icon creation
	greenIcon := createSimpleIcon(0, 255, 0)
	if len(greenIcon) == 0 {
		t.Error("Green icon should not be empty")
	}

	yellowIcon := createSimpleIcon(255, 255, 0)
	if len(yellowIcon) == 0 {
		t.Error("Yellow icon should not be empty")
	}

	redIcon := createSimpleIcon(255, 0, 0)
	if len(redIcon) == 0 {
		t.Error("Red icon should not be empty")
	}

	grayIcon := createSimpleIcon(128, 128, 128)
	if len(grayIcon) == 0 {
		t.Error("Gray icon should not be empty")
	}
}

func TestIconFormat(t *testing.T) {
	icon := createSimpleIcon(255, 0, 0)

	if runtime.GOOS == "windows" {
		// ICO format should start with specific header bytes
		if len(icon) < 6 {
			t.Error("ICO icon too small for header")
		}
		// Check ICO signature (first 4 bytes should be 0x00, 0x00, 0x01, 0x00)
		if icon[0] != 0x00 || icon[1] != 0x00 || icon[2] != 0x01 || icon[3] != 0x00 {
			t.Error("Invalid ICO header signature")
		}
	} else {
		// PNG format should start with PNG signature
		if len(icon) < 8 {
			t.Error("PNG icon too small for header")
		}
		// Check PNG signature
		pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		for i, b := range pngSignature {
			if icon[i] != b {
				t.Errorf("Invalid PNG header at byte %d: got %02x, want %02x", i, icon[i], b)
			}
		}
	}
}

func TestGetIconFunctions(t *testing.T) {
	// Test the public icon functions
	icons := []struct {
		name string
		fn   func() []byte
	}{
		{"green", getGreenIcon},
		{"yellow", getYellowIcon},
		{"red", getRedIcon},
		{"gray", getGrayIcon},
	}

	for _, icon := range icons {
		data := icon.fn()
		if len(data) == 0 {
			t.Errorf("%s icon should not be empty", icon.name)
		}
	}
}
