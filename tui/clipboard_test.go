package tui

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestSetClipboard(t *testing.T) {
	// Test text to copy
	testText := "SecretShare test content"

	// Try to set clipboard
	err := SetClipboard(testText)

	// Check results based on platform
	switch runtime.GOOS {
	case "darwin": // macOS
		// On macOS, check if pbcopy exists
		_, lookErr := exec.LookPath("pbcopy")
		if lookErr != nil {
			// If pbcopy doesn't exist, SetClipboard should return an error
			if err == nil {
				t.Error("Expected error on macOS when pbcopy is not available, but got nil")
			}
		} else {
			// If pbcopy exists, SetClipboard should not return an error
			if err != nil {
				t.Errorf("Unexpected error on macOS: %v", err)
			}
		}
	case "linux": // Linux
		// On Linux, check if xclip exists
		_, lookErr := exec.LookPath("xclip")
		if lookErr != nil {
			// If xclip doesn't exist, SetClipboard should return an error
			if err == nil {
				t.Error("Expected error on Linux when xclip is not available, but got nil")
			}
		} else {
			// If xclip exists, SetClipboard should not return an error
			if err != nil {
				t.Errorf("Unexpected error on Linux: %v", err)
			}
		}
	case "windows": // Windows
		// check if clip exists
		_, lookErr := exec.LookPath("clip")
		if lookErr != nil {
			// If clip doesn't exist, SetClipboard should return an error
			if err == nil {
				t.Error("Expected error on Windows when clip is not available, but got nil")
			}
		} else {
			// If clip exists, SetClipboard should not return an error
			if err != nil {
				t.Errorf("Unexpected error on Linux: %v", err)
			}
		}
	default: // other platforms
		// On unsupported platforms, SetClipboard should always return an error
		if err == nil {
			t.Error("Expected error on unsupported platforms, but got nil")
		}
	}
}
