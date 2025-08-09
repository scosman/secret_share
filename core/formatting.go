package core

import (
	"fmt"
)

// FormatPublicKey formats a public key with XML-like tags for sharing
func FormatPublicKey(key []byte) string {
	// We add a version number for future upgradeability
	return fmt.Sprintf("<secret_share_key>ssv1%s</secret_share_key>", string(key))
}

// FormatSecret formats an encrypted secret with XML-like tags for sharing
func FormatSecret(secret []byte) string {
	return fmt.Sprintf("<secret_share_secret>%s</secret_share_secret>", string(secret))
}
