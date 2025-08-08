package core

import (
	"crypto/rsa"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if privateKey == nil {
		t.Error("Private key should not be nil")
	}

	if publicKey == nil {
		t.Error("Public key should not be nil")
	}

	if privateKey.PublicKey.N.Cmp(publicKey.N) != 0 {
		t.Error("Generated private and public keys do not match")
	}

	if privateKey.PublicKey.E != publicKey.E {
		t.Error("Generated private and public keys do not match")
	}
}

func TestGenerateSymmetricKey(t *testing.T) {
	key, err := GenerateSymmetricKey()
	if err != nil {
		t.Fatalf("Failed to generate symmetric key: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected key length of 32 bytes, got %d", len(key))
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce, err := GenerateNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}

	if len(nonce) != 12 {
		t.Errorf("Expected nonce length of 12 bytes, got %d", len(nonce))
	}
}

func TestHybridEncryptDecrypt(t *testing.T) {
	// Generate key pair
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Test data
	testData := []byte("This is a secret message for testing hybrid encryption")

	// Encrypt
	encryptedData, err := HybridEncrypt(publicKey, testData)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	if len(encryptedData) == 0 {
		t.Error("Encrypted data should not be empty")
	}

	// Check that the encrypted data starts with "ssv1"
	if len(encryptedData) < 4 || string(encryptedData[0:4]) != "ssv1" {
		t.Error("Encrypted data should start with 'ssv1' format version")
	}

	// Decrypt
	decryptedData, err := HybridDecrypt(privateKey, encryptedData)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// Compare
	if string(testData) != string(decryptedData) {
		t.Errorf("Decrypted data does not match original. Expected: %s, Got: %s",
			string(testData), string(decryptedData))
	}
}

func TestPublicKeyToBytesAndBack(t *testing.T) {
	// Generate key pair
	_, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Convert to bytes
	publicKeyBytes, err := PublicKeyToBytes(publicKey)
	if err != nil {
		t.Fatalf("Failed to convert public key to bytes: %v", err)
	}

	if len(publicKeyBytes) == 0 {
		t.Error("Public key bytes should not be empty")
	}

	// Convert back to public key
	reconstructedPublicKey, err := BytesToPublicKey(publicKeyBytes)
	if err != nil {
		t.Fatalf("Failed to convert bytes to public key: %v", err)
	}

	if reconstructedPublicKey == nil {
		t.Error("Reconstructed public key should not be nil")
	}

	// Compare
	if publicKey.N.Cmp(reconstructedPublicKey.N) != 0 {
		t.Error("Original and reconstructed public keys do not match")
	}

	if publicKey.E != reconstructedPublicKey.E {
		t.Error("Original and reconstructed public keys do not match")
	}
}

func TestReceiverSenderSession(t *testing.T) {
	// Create receiver session
	receiverSession, err := NewReceiverSession()
	if err != nil {
		t.Fatalf("Failed to create receiver session: %v", err)
	}

	if receiverSession.privateKey == nil {
		t.Error("Receiver session private key should not be nil")
	}

	if receiverSession.publicKey == nil {
		t.Error("Receiver session public key should not be nil")
	}

	// Get public key
	publicKey := receiverSession.GetPublicKey()
	if publicKey == nil {
		t.Error("Public key should not be nil")
	}

	// Create sender session
	senderSession := NewSenderSession(publicKey)
	if senderSession == nil {
		t.Error("Sender session should not be nil")
	}

	if senderSession.receiverPublicKey == nil {
		t.Error("Sender session receiver public key should not be nil")
	}

	// Test secret encryption
	secret := []byte("Test secret message")
	encryptedSecret, err := senderSession.EncryptSecret(secret)
	if err != nil {
		t.Fatalf("Failed to encrypt secret: %v", err)
	}

	if len(encryptedSecret) == 0 {
		t.Error("Encrypted secret should not be empty")
	}

	// Test secret decryption
	decryptedSecret, err := receiverSession.DecryptSecret(encryptedSecret)
	if err != nil {
		t.Fatalf("Failed to decrypt secret: %v", err)
	}

	if string(secret) != string(decryptedSecret) {
		t.Errorf("Decrypted secret does not match original. Expected: %s, Got: %s",
			string(secret), string(decryptedSecret))
	}
}

func TestSenderSessionWithNilPublicKey(t *testing.T) {
	// Create sender session with nil public key
	senderSession := NewSenderSession(nil)

	// Try to encrypt secret
	_, err := senderSession.EncryptSecret([]byte("test"))
	if err == nil {
		t.Error("Expected error when encrypting with nil public key")
	}
}

func TestReceiverSessionWithNilPrivateKey(t *testing.T) {
	// Create receiver session with nil private key
	receiverSession := &ReceiverSession{
		privateKey: nil,
		publicKey:  &rsa.PublicKey{},
	}

	// Try to decrypt secret
	_, err := receiverSession.DecryptSecret([]byte("test"))
	if err == nil {
		t.Error("Expected error when decrypting with nil private key")
	}
}

func TestHybridDecryptValidFormat(t *testing.T) {
	// Generate key pair
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Test data
	testData := []byte("This is a secret message for testing hybrid encryption")

	// Encrypt data (which will add "ssv1" prefix)
	encryptedData, err := HybridEncrypt(publicKey, testData)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// Decrypt data
	decryptedData, err := HybridDecrypt(privateKey, encryptedData)
	if err != nil {
		t.Fatalf("Failed to decrypt data with valid format: %v", err)
	}

	// Compare
	if string(testData) != string(decryptedData) {
		t.Errorf("Decrypted data does not match original. Expected: %s, Got: %s",
			string(testData), string(decryptedData))
	}
}

func TestHybridDecryptNewerVersion(t *testing.T) {
	// Generate key pair
	privateKey, _, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create test data with "ssv" prefix but not "ssv1"
	testData := []byte("ssv2some data that would normally be encrypted")

	// Try to decrypt
	_, err = HybridDecrypt(privateKey, testData)
	if err == nil {
		t.Error("Expected error when decrypting data with newer version")
	}

	// Check error message
	expectedMsg := "this secret was sent using a newer version of SecretShare - please upgrade"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestHybridDecryptInvalidFormat(t *testing.T) {
	// Generate key pair
	privateKey, _, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create test data with invalid format
	testData := []byte("invalid format data")

	// Try to decrypt
	_, err = HybridDecrypt(privateKey, testData)
	if err == nil {
		t.Error("Expected error when decrypting data with invalid format")
	}

	// Check error message
	expectedMsg := "invalid encrypted data format"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}
