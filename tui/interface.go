package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Yellow = "\033[33m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
)

// Styled text functions
func headerText(text string) string {
	return Bold + Cyan + text + Reset
}

func promptText(text string) string {
	return Green + "❯ " + text + Reset
}

func errorText(text string) string {
	return Red + "✖ Error: " + text + Reset
}

func successText(text string) string {
	return Green + "✔ " + text + Reset
}

func infoText(text string) string {
	return Yellow + "ℹ " + text + Reset
}

// PromptUser displays a prompt and waits for user input
// Adds proper spacing and styling
func PromptUser(prompt string) string {
	fmt.Println()
	fmt.Print(promptText(prompt))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// PromptSecret displays a prompt and waits for user input, masking the characters
// Adds proper spacing and styling
func PromptSecret(prompt string) string {
	fmt.Println()
	fmt.Print(promptText(prompt))

	// Read password (masked input)
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add a newline after the masked input

	if err != nil {
		return ""
	}

	return string(bytes)
}

// PrintMessage displays a message to the user with better formatting
func PrintMessage(message string) {
	fmt.Println()
	fmt.Println(message)
}

// PrintHeader displays a header message with styling
func PrintHeader(message string) {
	fmt.Println()
	fmt.Println(headerText(message))
}

// PrintError displays an error message to the user with better formatting and colors
func PrintError(message string) {
	fmt.Println()
	fmt.Println(errorText(message))
}

// PrintSuccess displays a success message to the user with better formatting and colors
func PrintSuccess(message string) {
	fmt.Println()
	fmt.Println(successText(message))
}

// PrintInfo displays an info message to the user with better formatting and colors
func PrintInfo(message string) {
	fmt.Println()
	fmt.Println(infoText(message))
}

// IsQuit checks if the user input is a quit command
func IsQuit(input string) bool {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower == "q" || lower == "quit" || lower == "[q]" || lower == "exit"
}

// ParseRoleInput parses the user's role selection input
func ParseRoleInput(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)

	// Accept various forms of "send" input
	if lower == "s" || lower == "[s]" || lower == "send" || lower == "sender" {
		return "sender"
	}

	// Accept various forms of "receive" input
	if lower == "r" || lower == "[r]" || lower == "receive" || lower == "receiver" || lower == "recv" {
		return "receiver"
	}

	return ""
}

// ExtractPublicKey extracts the public key from XML-like tags
func ExtractPublicKey(input string) string {
	return extractTagContent(input, "secret_share_key")
}

// ExtractSecret extracts the secret from XML-like tags
func ExtractSecret(input string) string {
	return extractTagContent(input, "secret_share_secret")
}

// extractTagContent extracts content from XML-like tags with tolerance for formatting errors
func extractTagContent(input, tag string) string {
	// Normalize input by removing extra spaces
	normalized := strings.ReplaceAll(input, " ", "")

	// Try to find the opening tag with more tolerance
	openTag := fmt.Sprintf("<%s>", tag)
	startIdx := strings.Index(normalized, openTag)
	if startIdx == -1 {
		// Try variations with missing characters
		startIdx = strings.Index(normalized, tag+">")
		if startIdx == -1 {
			startIdx = strings.Index(normalized, tag)
			if startIdx == -1 {
				return ""
			}
		}
		// Adjust start index to after the tag
		if strings.HasSuffix(normalized[:startIdx], "<") {
			startIdx += len(tag) + 1 // +1 for >
		} else {
			startIdx += len(tag)
			// Find the closing >
			gtIdx := strings.Index(normalized[startIdx:], ">")
			if gtIdx == -1 {
				return ""
			}
			startIdx += gtIdx + 1
		}
	} else {
		startIdx += len(openTag)
	}

	// Try to find the closing tag with more tolerance
	closeTag := fmt.Sprintf("</%s>", tag)
	endIdx := strings.Index(normalized[startIdx:], closeTag)
	if endIdx == -1 {
		// Try variations with missing characters
		endIdx = strings.Index(normalized[startIdx:], "<"+tag)
		if endIdx == -1 {
			endIdx = strings.Index(normalized[startIdx:], "<")
			if endIdx == -1 {
				return normalized[startIdx:] // Return everything if no closing tag found
			}
		}
		endIdx += startIdx
	} else {
		endIdx += startIdx
	}

	if endIdx <= startIdx {
		return normalized[startIdx:] // Return everything from start if no proper end
	}

	return normalized[startIdx:endIdx]
}

// FormatPublicKey formats a public key with XML-like tags for sharing
func FormatPublicKey(key []byte) string {
	return fmt.Sprintf("<secret_share_key>%s</secret_share_key>", string(key))
}

// FormatSecret formats an encrypted secret with XML-like tags for sharing
func FormatSecret(secret []byte) string {
	return fmt.Sprintf("<secret_share_secret>%s</secret_share_secret>", string(secret))
}
