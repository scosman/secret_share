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

	// Test case 7: Missing both brackets
	input7 := "secret_share_key>TEST_KEY_CONTENT</secret_share_key"
	expected7 := "TEST_KEY_CONTENT"
	result7 := ExtractPublicKey(input7)
	if result7 != expected7 {
		t.Errorf("Test 7 failed: Expected '%s', got '%s'", expected7, result7)
	}

	// Test case 8: Missing opening bracket and closing bracket on secret
	input8 := "secret_share_secret>TEST_SECRET_CONTENT<secret_share_secret"
	expected8 := "TEST_SECRET_CONTENT"
	result8 := ExtractSecret(input8)
	if result8 != expected8 {
		t.Errorf("Test 8 failed: Expected '%s', got '%s'", expected8, result8)
	}

	// Test case 9: Missing closing bracket on opening tag
	input9 := "<secret_share_keyTEST_KEY_CONTENT</secret_share_key>"
	expected9 := "TEST_KEY_CONTENT"
	result9 := ExtractPublicKey(input9)
	if result9 != expected9 {
		t.Errorf("Test 9 failed: Expected '%s', got '%s'", expected9, result9)
	}

	// Test case 10: Missing opening bracket on closing tag
	input10 := "<secret_share_key>TEST_KEY_CONTENT/secret_share_key>"
	expected10 := "TEST_KEY_CONTENT"
	result10 := ExtractPublicKey(input10)
	if result10 != expected10 {
		t.Errorf("Test 10 failed: Expected '%s', got '%s'", expected10, result10)
	}

	// Test case 12: Whitespace trimming
	input12 := "  <secret_share_key>  TEST_KEY_CONTENT  </secret_share_key>  "
	expected12 := "TEST_KEY_CONTENT"
	result12 := ExtractPublicKey(input12)
	if result12 != expected12 {
		t.Errorf("Test 12 failed: Expected '%s', got '%s'", expected12, result12)
	}

	// Test case 13: Extraneous data before and after tags
	input13 := "Some extra data <secret_share_key>TEST_KEY_CONTENT</secret_share_key> More extra data"
	expected13 := "TEST_KEY_CONTENT"
	result13 := ExtractPublicKey(input13)
	if result13 != expected13 {
		t.Errorf("Test 13 failed: Expected '%s', got '%s'", expected13, result13)
	}

	// Test case 14: Extraneous data with secret tags
	input14 := "Extra stuff <secret_share_secret>TEST_SECRET_CONTENT</secret_share_secret> Even more stuff"
	expected14 := "TEST_SECRET_CONTENT"
	result14 := ExtractSecret(input14)
	if result14 != expected14 {
		t.Errorf("Test 14 failed: Expected '%s', got '%s'", expected14, result14)
	}

	// Test case 15: Combination of issues - missing chars and extraneous data
	input15 := "Extra data secret_share_key>TEST_KEY_CONTENT</secret_share_key More data"
	expected15 := "TEST_KEY_CONTENT"
	result15 := ExtractPublicKey(input15)
	if result15 != expected15 {
		t.Errorf("Test 15 failed: Expected '%s', got '%s'", expected15, result15)
	}

	// Test case 16: Empty input
	input16 := ""
	expected16 := ""
	result16 := ExtractPublicKey(input16)
	if result16 != expected16 {
		t.Errorf("Test 16 failed: Expected '%s', got '%s'", expected16, result16)
	}

	// Test case 17: Input with no tags and no content
	input17 := "   "
	expected17 := ""
	result17 := ExtractPublicKey(input17)
	if result17 != expected17 {
		t.Errorf("Test 17 failed: Expected '%s', got '%s'", expected17, result17)
	}

	// Additional test cases for enhanced flexibility

	// Test case 18: Partially missing opening tag characters
	input18 := "cret_share_key>TEST_KEY_CONTENT</secret_share_key>"
	expected18 := "TEST_KEY_CONTENT"
	result18 := ExtractPublicKey(input18)
	if result18 != expected18 {
		t.Errorf("Test 18 failed: Expected '%s', got '%s'", expected18, result18)
	}

	// Test case 19: Partially missing closing tag characters
	input19 := "<secret_share_key>TEST_KEY_CONTENT</secret_share_sec"
	expected19 := "TEST_KEY_CONTENT"
	result19 := ExtractPublicKey(input19)
	if result19 != expected19 {
		t.Errorf("Test 19 failed: Expected '%s', got '%s'", expected19, result19)
	}

	// Test case 20: Partially missing both opening and closing tag characters
	input20 := "cret_share_key>TEST_KEY_CONTENT</secret_share_sec"
	expected20 := "TEST_KEY_CONTENT"
	result20 := ExtractPublicKey(input20)
	if result20 != expected20 {
		t.Errorf("Test 20 failed: Expected '%s', got '%s'", expected20, result20)
	}

	// Test case 22: Missing key part of closing tag
	input22 := "<secret_share_key>TEST_KEY_CONTENT</secret_share>"
	expected22 := "TEST_KEY_CONTENT"
	result22 := ExtractPublicKey(input22)
	if result22 != expected22 {
		t.Errorf("Test 22 failed: Expected '%s', got '%s'", expected22, result22)
	}
}
