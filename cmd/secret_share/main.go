package main

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/scosman/secret_share/core"
	"github.com/scosman/secret_share/tui"
)

const TitleCard = `
           ▄▖        ▗ ▄▖▌       
           ▚ █▌▛▘▛▘█▌▜▘▚ ▛▌▀▌▛▘█▌
           ▄▌▙▖▙▖▌ ▙▖▐▖▄▌▌▌█▌▌ ▙▖
  
      Secure One Time Secret Sharing`

func main() {
	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		tui.PrintMessage("\nShutting down SecretShare...")
		os.Exit(0)
	}()

	// Welcome message
	tui.PrintHeader(TitleCard)

	// Get user role
	role := getUserRole()
	if role == "" {
		return
	}

	// Handle based on role
	if role == "receiver" {
		handleReceiver()
	} else {
		handleSender()
	}
}

func getUserRole() string {
	for {
		input := tui.PromptUserSingleChar("Are you [s]ending or [r]eceiving a secret? ")
		if tui.IsQuit(input) {
			tui.PrintMessage("Quiting SecretShare")
			return ""
		}

		role := tui.ParseRoleInput(input)
		if role == "" {
			tui.PrintError("Invalid input. Please enter 's' for sending or 'r' for receiving (or 'q' to quit).")
			continue
		}

		return role
	}
}

func handleReceiver() {
	// Create a new receiver session
	session, err := core.NewReceiverSession()
	if err != nil {
		tui.PrintError(fmt.Sprintf("Failed to create receiver session: %v", err))
		return
	}

	// Get public key bytes
	publicKeyBytes, err := core.PublicKeyToBytes(session.GetPublicKey())
	if err != nil {
		tui.PrintError(fmt.Sprintf("Failed to serialize public key: %v", err))
		return
	}

	// Display public key for sharing
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKeyBytes)
	publicKeyFormatted := core.FormatPublicKey([]byte(publicKeyStr))
	tui.PrintInfo("Here's a new the public key:")
	tui.PrintMessage(publicKeyFormatted)

	// Try to copy public key to clipboard
	err = tui.SetClipboard(publicKeyFormatted)
	if err == nil {
		tui.PrintInfo("Copied to clipboard.")
	}

	// Get encrypted secret from sender with retry logic
	var decryptedSecret []byte
	for {
		input := tui.PromptUser("Send the key above to the person who wants to share a secret with you. When they reply back with the encrypted secret, enter it here: ")
		if tui.IsQuit(input) {
			tui.PrintMessage("Quiting SecretShare")
			return
		}

		// Extract secret from tags
		secretStr := tui.ExtractSecret(input)
		// Decode base64 secret
		encryptedSecret, err := base64.StdEncoding.DecodeString(secretStr)
		// Decrypt the secret
		if err == nil {
			decryptedSecret, err = session.DecryptSecret(encryptedSecret)
		}

		if err != nil || secretStr == "" {
			tui.PrintError("Could not extract secret from input.")
			tui.PrintMessage("Ensure you are pasting the exact encrypted secret from the sender. It should be a string wrapped in tags like '<secret_share_secret>'.")
			continue
		}

		break
	}

	// Display the decrypted secret
	tui.PrintSuccess(fmt.Sprintf("Here's your secret: %s", string(decryptedSecret)))
}

func handleSender() {
	// Get receiver's public key with retry logic
	var receiverPublicKey *rsa.PublicKey
	for {
		input := tui.PromptUser("Enter the secret key from the other person. It should be a string wrapped in <secret_share_key> tags: ")
		if tui.IsQuit(input) {
			tui.PrintMessage("Quiting SecretShare")
			return
		}

		// Extract public key from tags
		publicKeyStr := tui.ExtractPublicKey(input)
		// Decode base64 public key
		publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
		// Parse public key
		if err == nil {
			receiverPublicKey, err = core.BytesToPublicKey(publicKeyBytes)
		}

		if err != nil || publicKeyStr == "" {
			tui.PrintError("Could not extract public key from input.")
			tui.PrintMessage("Ensure you are pasting the exact secret key from the sender. It should be a string wrapped in tags like '<secret_share_key>'.")
			continue
		}

		break
	}

	// Create sender session
	session := core.NewSenderSession(receiverPublicKey)

	// Get secret to share
	secret := tui.PromptSecret("Enter the secret you want to share: ")
	if tui.IsQuit(secret) {
		tui.PrintMessage("Quiting SecretShare")
		return
	}

	// Encrypt the secret
	encryptedSecret, err := session.EncryptSecret([]byte(secret))
	if err != nil {
		tui.PrintError(fmt.Sprintf("Failed to encrypt secret: %v", err))
		return
	}

	// Encode encrypted secret as base64
	encryptedSecretStr := base64.StdEncoding.EncodeToString(encryptedSecret)
	encryptedSecretFormatted := core.FormatSecret([]byte(encryptedSecretStr))

	// Display the encrypted secret for sharing
	tui.PrintSuccess("Here's the secret encrypted so only they can decrypt it:")
	tui.PrintMessage(encryptedSecretFormatted)

	// Try to copy encrypted secret to clipboard
	err = tui.SetClipboard(encryptedSecretFormatted)
	if err == nil {
		tui.PrintInfo("Copied to clipboard. Send this secret back to the person who shared their key with you.")
	} else {
		tui.PrintInfo("Send this secret back to the person who shared their key with you.")
	}
}
