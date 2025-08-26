# Paste Detection Feature

This document describes the paste detection functionality implemented in SecretShare to handle multi-line secrets and keys that contain newline characters.

## Problem

Previously, when pasting secrets or keys that contained newline characters, the input would be terminated at the first newline, truncating the content. This was problematic for:
- Multi-line certificates
- Base64 encoded content with line breaks
- Secrets that naturally contain newlines

## Solution

A timing-based approach is used to distinguish between:
1. **Manual Enter**: User manually presses Enter (≥15ms gap from previous character) → Terminates input
2. **Pasted Enter**: Newline character in pasted content (<15ms gap from previous character) → Preserved as part of input

## Implementation Details

### Timing Threshold
- **15 milliseconds** is used as the threshold
- Manual typing typically has 100-200ms+ gaps between characters
- Paste operations typically have <1ms gaps between characters
- 15ms provides a safe buffer for various typing speeds and system performance

### Functions Modified
- `PromptUser()`: Now uses `readInputWithPasteDetection()`
- `PromptSecret()`: Now uses `readSecretWithPasteDetection()`

### Features
- Character echoing for regular input (masked with `*` for secrets)
- Backspace support
- Ctrl+C handling (treated as quit)
- Fallback to original behavior if terminal raw mode fails

### Character Handling
- Printable characters (ASCII 32+) and tabs are accepted
- Newlines (`\n`) and carriage returns (`\r`) are handled based on timing
- Backspace (127, 8) removes last character
- Ctrl+C (3) terminates with quit signal

## Testing

To test the functionality:

1. **Manual typing test**:
   - Type text character by character
   - Press Enter manually
   - Should terminate input

2. **Paste test**:
   - Copy multi-line text to clipboard
   - Paste into the prompt
   - Should preserve all newlines as part of the input

## Backward Compatibility

The implementation maintains backward compatibility:
- Single-line input continues to work as before
- Manual Enter still terminates input
- Fallback to original behavior if raw terminal mode is unavailable
