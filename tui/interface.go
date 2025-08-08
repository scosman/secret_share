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

// PromptUserSingleChar displays a prompt and waits for a single character input
// Adds proper spacing and styling
func PromptUserSingleChar(prompt string) string {
	fmt.Println()
	fmt.Print(promptText(prompt))

	// Put terminal in raw mode to read single character
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	// Read a single character
	bytes := make([]byte, 1)
	_, err = os.Stdin.Read(bytes)
	if err != nil {
		return ""
	}

	// Convert to string and return
	char := string(bytes)

	// Echo the character to the terminal since we read it directly
	fmt.Print(char)
	fmt.Println()

	return char
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
	// Trim whitespace from the entire input
	input = strings.TrimSpace(input)

	// Create properly formatted tags
	openTag := fmt.Sprintf("<%s>", tag)
	closeTag := fmt.Sprintf("</%s>", tag)

	// First, check if we have content without any tags (plain content)
	if !strings.Contains(input, "<") && !strings.Contains(input, ">") {
		return input
	}

	// First, attempt to find the end tag. If found, remove it and the content after it.
	endIdx := strings.Index(input, closeTag)
	if endIdx == -1 {
		// If end tag was not found, find the first '</' and remove it and content after it.
		closeTagStart := "</"
		closeTagIdx := strings.Index(input, closeTagStart)
		if closeTagIdx != -1 {
			input = input[:closeTagIdx]
		}
	} else {
		input = input[:endIdx]
	}

	// Then, attempt to find the start tag. If found, remove it and the content before it.
	startIdx := strings.Index(input, openTag)
	if startIdx == -1 {
		// If start tag was not found, find the first '>' and remove it and content before it.
		gtIdx := strings.Index(input, ">")
		if gtIdx != -1 {
			input = input[gtIdx+1:]
		}
	} else {
		input = input[startIdx+len(openTag):]
	}

	// Return whatever string is left
	return strings.TrimSpace(input)
}
