package core

import (
	"crypto/rsa"
	"fmt"
)

// ReceiverSession represents a session where the user is receiving a secret
type ReceiverSession struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// SenderSession represents a session where the user is sending a secret
type SenderSession struct {
	receiverPublicKey *rsa.PublicKey
}

// NewReceiverSession creates a new receiver session with a fresh key pair
func NewReceiverSession() (*ReceiverSession, error) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	return &ReceiverSession{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// NewSenderSession creates a new sender session with the receiver's public key
func NewSenderSession(receiverPublicKey *rsa.PublicKey) *SenderSession {
	return &SenderSession{
		receiverPublicKey: receiverPublicKey,
	}
}

// GetPublicKey returns the public key for sharing (receiver session only)
func (rs *ReceiverSession) GetPublicKey() *rsa.PublicKey {
	return rs.publicKey
}

// EncryptSecret encrypts a secret using the receiver's public key
func (ss *SenderSession) EncryptSecret(secret []byte) ([]byte, error) {
	if ss.receiverPublicKey == nil {
		return nil, fmt.Errorf("receiver public key is not set")
	}

	encryptedData, err := HybridEncrypt(ss.receiverPublicKey, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	return encryptedData, nil
}

// DecryptSecret decrypts a secret using the receiver's private key
func (rs *ReceiverSession) DecryptSecret(encryptedSecret []byte) ([]byte, error) {
	if rs.privateKey == nil {
		return nil, fmt.Errorf("private key is not set")
	}

	decryptedData, err := HybridDecrypt(rs.privateKey, encryptedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return decryptedData, nil
}
