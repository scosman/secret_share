package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"secret_share/core"
	"secret_share/tui"
	"syscall"
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
			tui.PrintMessage("Shutting down SecretShare...")
			return ""
		}

		role := tui.ParseRoleInput(input)
		if role == "" {
			tui.PrintError("Invalid input. Please enter 's' for sending or 'r' for receiving.")
			tui.PrintMessage("You can quit by typing 'q'")
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
		tui.PrintInfo("Copied to clipboard. Send this key to the person who wants to share a secret with you.")
	} else {
		tui.PrintInfo("You must send this key to the person who wants to share a secret with you.")
	}

	tui.PrintMessage("")

	// Get encrypted secret from sender with retry logic
	var secretStr string
	for {
		input := tui.PromptUser("Enter the secret from the other person (can be just the secret or wrapped in <secret_share_secret> tags): ")
		if tui.IsQuit(input) {
			tui.PrintMessage("Shutting down SecretShare...")
			return
		}

		// Extract secret from tags
		secretStr = tui.ExtractSecret(input)
		if secretStr == "" {
			tui.PrintError("Could not extract secret from input. Please make sure it's properly formatted.")
			tui.PrintMessage("You can enter just the base64-encoded secret, or wrap it in tags like <secret_share_secret>YOUR_SECRET_HERE</secret_share_secret>")
			tui.PrintMessage("If you're having trouble, try copying and pasting the entire message from the sender")
			continue
		}
		break
	}

	// Decode base64 secret
	encryptedSecret, err := base64.StdEncoding.DecodeString(secretStr)
	if err != nil {
		tui.PrintError("Invalid base64 encoding in secret.")
		return
	}

	// Decrypt the secret
	decryptedSecret, err := session.DecryptSecret(encryptedSecret)
	if err != nil {
		tui.PrintError(fmt.Sprintf("Failed to decrypt secret: %v", err))
		return
	}

	// Display the decrypted secret
	tui.PrintSuccess(fmt.Sprintf("Here's your secret: %s", string(decryptedSecret)))
}

func handleSender() {
	// Get receiver's public key with retry logic
	var publicKeyStr string
	for {
		input := tui.PromptUser("Enter the secret key from the other person (can be just the key or wrapped in <secret_share_key> tags): ")
		if tui.IsQuit(input) {
			tui.PrintMessage("Shutting down SecretShare...")
			return
		}

		// Extract public key from tags
		publicKeyStr = tui.ExtractPublicKey(input)
		if publicKeyStr == "" {
			tui.PrintError("Could not extract public key from input. Please make sure it's properly formatted.")
			tui.PrintMessage("You can enter just the base64-encoded key, or wrap it in tags like <secret_share_key>YOUR_KEY_HERE</secret_share_key>")
			tui.PrintMessage("If you're having trouble, try copying and pasting the entire message from the receiver")
			continue
		}
		break
	}

	// Decode base64 public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		tui.PrintError("Invalid base64 encoding in public key.")
		return
	}

	// Parse public key
	receiverPublicKey, err := core.BytesToPublicKey(publicKeyBytes)
	if err != nil {
		tui.PrintError(fmt.Sprintf("Failed to parse public key: %v", err))
		return
	}

	// Create sender session
	session := core.NewSenderSession(receiverPublicKey)

	// Get secret to share
	secret := tui.PromptSecret("Enter the secret you want to share: ")
	if tui.IsQuit(secret) {
		tui.PrintMessage("Shutting down SecretShare...")
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
