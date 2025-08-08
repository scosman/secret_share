package tui

import (
	"testing"
)

func TestExtractTagContent(t *testing.T) {
	// Test case 1: Properly formatted tags
	input1 := "<secret_share_key>TEST_KEY_CONTENT</secret_share_key>"
	expected1 := "TEST_KEY_CONTENT"
	result1 := ExtractPublicKey(input1)
	if result1 != expected1 {
		t.Errorf("Test 1 failed: Expected '%s', got '%s'", expected1, result1)
	}

	// Test case 2: Properly formatted secret tags
	input2 := "<secret_share_secret>TEST_SECRET_CONTENT</secret_share_secret>"
	expected2 := "TEST_SECRET_CONTENT"
	result2 := ExtractSecret(input2)
	if result2 != expected2 {
		t.Errorf("Test 2 failed: Expected '%s', got '%s'", expected2, result2)
	}

	// Test case 3: No XML tags, just content
	input3 := "TEST_KEY_CONTENT"
	expected3 := "TEST_KEY_CONTENT"
	result3 := ExtractPublicKey(input3)
	if result3 != expected3 {
		t.Errorf("Test 3 failed: Expected '%s', got '%s'", expected3, result3)
	}

	// Test case 4: No XML tags, just secret content
	input4 := "TEST_SECRET_CONTENT"
	expected4 := "TEST_SECRET_CONTENT"
	result4 := ExtractSecret(input4)
	if result4 != expected4 {
		t.Errorf("Test 4 failed: Expected '%s', got '%s'", expected4, result4)
	}

	// Test case 5: Missing opening bracket
	input5 := "secret_share_key>TEST_KEY_CONTENT</secret_share_key>"
	expected5 := "TEST_KEY_CONTENT"
	result5 := ExtractPublicKey(input5)
	if result5 != expected5 {
		t.Errorf("Test 5 failed: Expected '%s', got '%s'", expected5, result5)
	}

	// Test case 6: Missing closing bracket
	input6 := "<secret_share_key>TEST_KEY_CONTENT</secret_share_key"
	expected6 := "TEST_KEY_CONTENT"
	result6 := ExtractPublicKey(input6)
	if result6 != expected6 {
		t.Errorf("Test 6 failed: Expected '%s', got '%s'", expected6, result6)
	}

	// Test case 7: Whitespace trimming
	input7 := "  <secret_share_key>  TEST_KEY_CONTENT  </secret_share_key>  "
	expected7 := "TEST_KEY_CONTENT"
	result7 := ExtractPublicKey(input7)
	if result7 != expected7 {
		t.Errorf("Test 7 failed: Expected '%s', got '%s'", expected7, result7)
	}

	// Test case 8: Extraneous data before and after tags
	input8 := "Some extra data <secret_share_key>TEST_KEY_CONTENT</secret_share_key> More extra data"
	expected8 := "TEST_KEY_CONTENT"
	result8 := ExtractPublicKey(input8)
	if result8 != expected8 {
		t.Errorf("Test 8 failed: Expected '%s', got '%s'", expected8, result8)
	}

	// Test case 9: Extraneous data with secret tags
	input9 := "Extra stuff <secret_share_secret>TEST_SECRET_CONTENT</secret_share_secret> Even more stuff"
	expected9 := "TEST_SECRET_CONTENT"
	result9 := ExtractSecret(input9)
	if result9 != expected9 {
		t.Errorf("Test 9 failed: Expected '%s', got '%s'", expected9, result9)
	}

	// Test case 10: Empty input
	input10 := ""
	expected10 := ""
	result10 := ExtractPublicKey(input10)
	if result10 != expected10 {
		t.Errorf("Test 10 failed: Expected '%s', got '%s'", expected10, result10)
	}

	// Test case 11: Input with no tags and no content
	input11 := "   "
	expected11 := ""
	result11 := ExtractPublicKey(input11)
	if result11 != expected11 {
		t.Errorf("Test 11 failed: Expected '%s', got '%s'", expected11, result11)
	}

	// Test case 12: Valid content with only end tag missing opening bracket
	input12 := "VALIDKEY</secret_share_key>"
	expected12 := "VALIDKEY"
	result12 := ExtractPublicKey(input12)
	if result12 != expected12 {
		t.Errorf("Test 12 failed: Expected '%s', got '%s'", expected12, result12)
	}
}
