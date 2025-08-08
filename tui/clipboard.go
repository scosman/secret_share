package tui

import (
	"os/exec"
	"runtime"
	"strings"
)

// SetClipboard copies the given text to the system clipboard
// Returns an error if the operation is not supported or fails
func SetClipboard(text string) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		return setClipboardMacOS(text)
	case "linux": // Linux
		return setClipboardLinux(text)
	default: // Windows and other platforms
		return exec.ErrNotFound
	}
}

// setClipboardMacOS copies text to clipboard on macOS using pbcopy
func setClipboardMacOS(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// setClipboardLinux copies text to clipboard on Linux using xclip
func setClipboardLinux(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
