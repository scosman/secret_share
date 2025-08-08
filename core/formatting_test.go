package core

import (
	"testing"
)

func TestFormatPublicKey(t *testing.T) {
	// Test case 1: Basic formatting
	key := []byte("test_public_key")
	expected := "<secret_share_key>test_public_key</secret_share_key>"
	result := FormatPublicKey(key)
	if result != expected {
		t.Errorf("Test 1 failed: Expected '%s', got '%s'", expected, result)
	}

	// Test case 2: Empty key
	emptyKey := []byte("")
	expectedEmpty := "<secret_share_key></secret_share_key>"
	resultEmpty := FormatPublicKey(emptyKey)
	if resultEmpty != expectedEmpty {
		t.Errorf("Test 2 failed: Expected '%s', got '%s'", expectedEmpty, resultEmpty)
	}

	// Test case 3: Key with special characters
	specialKey := []byte("key_with_special_chars_!@#$%^&*()")
	expectedSpecial := "<secret_share_key>key_with_special_chars_!@#$%^&*()</secret_share_key>"
	resultSpecial := FormatPublicKey(specialKey)
	if resultSpecial != expectedSpecial {
		t.Errorf("Test 3 failed: Expected '%s', got '%s'", expectedSpecial, resultSpecial)
	}
}

func TestFormatSecret(t *testing.T) {
	// Test case 1: Basic formatting
	secret := []byte("test_secret")
	expected := "<secret_share_secret>test_secret</secret_share_secret>"
	result := FormatSecret(secret)
	if result != expected {
		t.Errorf("Test 1 failed: Expected '%s', got '%s'", expected, result)
	}

	// Test case 2: Empty secret
	emptySecret := []byte("")
	expectedEmpty := "<secret_share_secret></secret_share_secret>"
	resultEmpty := FormatSecret(emptySecret)
	if resultEmpty != expectedEmpty {
		t.Errorf("Test 2 failed: Expected '%s', got '%s'", expectedEmpty, resultEmpty)
	}

	// Test case 3: Secret with special characters
	specialSecret := []byte("secret_with_special_chars_!@#$%^&*()")
	expectedSpecial := "<secret_share_secret>secret_with_special_chars_!@#$%^&*()</secret_share_secret>"
	resultSpecial := FormatSecret(specialSecret)
	if resultSpecial != expectedSpecial {
		t.Errorf("Test 3 failed: Expected '%s', got '%s'", expectedSpecial, resultSpecial)
	}
}
