package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

// GenerateKeyPair generates a new RSA key pair with 2048 bits
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// GenerateSymmetricKey generates a random 256-bit AES key
func GenerateSymmetricKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	return key, nil
}

// GenerateNonce generates a random nonce for AES-GCM
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, 12) // Standard GCM nonce size
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return nonce, nil
}

// HybridEncrypt encrypts data using hybrid encryption:
// 1. Generates a random AES-256 key
// 2. Encrypts the AES key with RSA-OAEP
// 3. Encrypts the data with AES-GCM
func HybridEncrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// Generate a random symmetric key
	symmetricKey, err := GenerateSymmetricKey()
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with RSA-OAEP
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, symmetricKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt symmetric key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Generate nonce
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, err
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Combine encrypted key length, encrypted key, nonce, and ciphertext
	// Format: [keyLen][encryptedKey][nonce][ciphertext]
	keyLen := len(encryptedKey)
	result := make([]byte, 4+len(encryptedKey)+len(nonce)+len(ciphertext))

	// Store key length as 4 bytes
	result[0] = byte(keyLen >> 24)
	result[1] = byte(keyLen >> 16)
	result[2] = byte(keyLen >> 8)
	result[3] = byte(keyLen)

	// Copy encrypted key
	copy(result[4:4+keyLen], encryptedKey)

	// Copy nonce
	copy(result[4+keyLen:4+keyLen+len(nonce)], nonce)

	// Copy ciphertext
	copy(result[4+keyLen+len(nonce):], ciphertext)

	return result, nil
}

// HybridDecrypt decrypts data using hybrid encryption:
// 1. Decrypts the AES key with RSA-OAEP
// 2. Decrypts the data with AES-GCM
func HybridDecrypt(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	if len(encryptedData) < 4 {
		return nil, fmt.Errorf("invalid encrypted data format")
	}

	// Extract key length
	keyLen := int(encryptedData[0])<<24 | int(encryptedData[1])<<16 | int(encryptedData[2])<<8 | int(encryptedData[3])

	if len(encryptedData) < 4+keyLen+12 {
		return nil, fmt.Errorf("invalid encrypted data format")
	}

	// Extract encrypted key
	encryptedKey := encryptedData[4 : 4+keyLen]

	// Extract nonce (12 bytes for GCM)
	nonceStart := 4 + keyLen
	nonce := encryptedData[nonceStart : nonceStart+12]

	// Extract ciphertext
	ciphertext := encryptedData[nonceStart+12:]

	// Decrypt the symmetric key with RSA-OAEP
	symmetricKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt symmetric key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// PublicKeyToBytes converts an RSA public key to bytes
func PublicKeyToBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(publicKey)
}

// BytesToPublicKey converts bytes to an RSA public key
func BytesToPublicKey(data []byte) (*rsa.PublicKey, error) {
	publicKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPub, nil
}
