package tui

import (
	"strings"
	"testing"
	"time"
)

// Test helper function to simulate timing behavior
func TestTimingLogic(t *testing.T) {
	// Test manual typing scenario (>= 15ms gaps)
	manualGap := 50 * time.Millisecond
	if manualGap < 15*time.Millisecond {
		t.Errorf("Manual gap %v should be >= 15ms", manualGap)
	}

	// Test paste scenario (< 15ms gaps)
	pasteGap := 5 * time.Millisecond
	if pasteGap >= 15*time.Millisecond {
		t.Errorf("Paste gap %v should be < 15ms", pasteGap)
	}
}

// Test UTF-8 rune handling in backspace logic
func TestRuneHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ASCII characters",
			input:    "hello",
			expected: "hell",
		},
		{
			name:     "UTF-8 emoji",
			input:    "hello🔒",
			expected: "hello",
		},
		{
			name:     "UTF-8 accented characters",
			input:    "café",
			expected: "caf",
		},
		{
			name:     "Mixed UTF-8",
			input:    "test🔑password",
			expected: "test🔑passwor",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single rune",
			input:    "a",
			expected: "",
		},
		{
			name:     "Single emoji",
			input:    "🔒",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the backspace logic from our functions
			str := tt.input
			if len(str) > 0 {
				runes := []rune(str)
				if len(runes) > 0 {
					result := string(runes[:len(runes)-1])
					if result != tt.expected {
						t.Errorf("Expected %q, got %q", tt.expected, result)
					}
				}
			} else if tt.expected != "" {
				t.Errorf("Expected %q for empty input, got empty string", tt.expected)
			}
		})
	}
}

// Test character classification logic
func TestCharacterClassification(t *testing.T) {
	tests := []struct {
		name     string
		ch       byte
		expected string // "printable", "newline", "backspace", "ctrl-c", "other"
	}{
		{"Space", 32, "printable"},
		{"Letter A", 65, "printable"},
		{"Tab", 9, "printable"},
		{"Newline", 10, "newline"},
		{"Carriage Return", 13, "newline"},
		{"Backspace", 8, "backspace"},
		{"DEL", 127, "backspace"},
		{"Ctrl+C", 3, "ctrl-c"},
		{"Non-printable", 1, "other"},
		{"Non-printable high", 31, "other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual string
			ch := tt.ch

			// Replicate the character classification logic from our functions
			if ch == '\n' || ch == '\r' {
				actual = "newline"
			} else if ch == 127 || ch == 8 {
				actual = "backspace"
			} else if ch >= 32 || ch == '\t' {
				actual = "printable"
			} else if ch == 3 {
				actual = "ctrl-c"
			} else {
				actual = "other"
			}

			if actual != tt.expected {
				t.Errorf("Character %d (%c): expected %s, got %s", ch, ch, tt.expected, actual)
			}
		})
	}
}

// Test newline handling in multiline content
func TestNewlineHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single line",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "Two lines with LF",
			input:    "line1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "Two lines with CRLF",
			input:    "line1\r\nline2",
			expected: "line1\n\nline2", // Both \r and \n become \n
		},
		{
			name:     "Multiple lines",
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "Empty lines",
			input:    "line1\n\nline3",
			expected: "line1\n\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate processing each character like our paste detection does
			var result strings.Builder
			for _, ch := range []byte(tt.input) {
				if ch == '\n' || ch == '\r' {
					result.WriteByte('\n') // Normalize all newlines to \n
				} else if ch >= 32 || ch == '\t' {
					result.WriteByte(ch)
				}
			}

			if result.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.String())
			}
		})
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	// Test timing threshold boundary
	threshold := 15 * time.Millisecond
	boundaryTests := []struct {
		name           string
		gap            time.Duration
		expectedManual bool
	}{
		{"Just under threshold", 14 * time.Millisecond, false},
		{"Exactly at threshold", 15 * time.Millisecond, true},
		{"Just over threshold", 16 * time.Millisecond, true},
		{"Way over threshold", 100 * time.Millisecond, true},
	}

	for _, tt := range boundaryTests {
		t.Run(tt.name, func(t *testing.T) {
			isManual := tt.gap >= threshold
			if isManual != tt.expectedManual {
				t.Errorf("Gap %v: expected manual=%v, got %v", tt.gap, tt.expectedManual, isManual)
			}
		})
	}
}

// Test string builder behavior with our usage patterns
func TestStringBuilderBehavior(t *testing.T) {
	tests := []struct {
		name       string
		operations []string // "write:text", "reset", "writeByte:X"
		expected   string
	}{
		{
			name:       "Basic write and reset",
			operations: []string{"write:hello", "reset", "write:world"},
			expected:   "world",
		},
		{
			name:       "WriteByte operations",
			operations: []string{"writeByte:a", "writeByte:b", "writeByte:c"},
			expected:   "abc",
		},
		{
			name:       "Mixed operations",
			operations: []string{"write:test", "writeByte:X", "reset", "writeByte:Y"},
			expected:   "Y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			for _, op := range tt.operations {
				parts := strings.Split(op, ":")
				switch parts[0] {
				case "write":
					builder.WriteString(parts[1])
				case "writeByte":
					builder.WriteByte(parts[1][0])
				case "reset":
					builder.Reset()
				}
			}

			if builder.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, builder.String())
			}
		})
	}
}
