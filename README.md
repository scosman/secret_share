# SecretShare

Securely share secrets over untrusted communication channels using hybrid encryption (RSA-OAEP + AES-GCM).

## Overview

SecretShare is a command-line tool that enables two parties to securely exchange sensitive information even when using untrusted communication methods like unencrypted chat platforms. It employs a hybrid encryption approach where:

1. The receiver generates a one-time RSA key pair
2. The receiver shares only the public key with the sender
3. The sender uses the public key to encrypt a randomly generated AES key
4. The sender uses the AES key to encrypt the actual secret with AES-GCM
5. The sender shares the encrypted data with the receiver
6. The receiver uses their private key to decrypt the AES key, then decrypts the secret

The private key never leaves the receiver's machine and is never exposed to the communication channel.

## Features

- **Secure Hybrid Encryption**: Uses RSA-OAEP for key exchange and AES-GCM for data encryption
- **One-Time Keys**: Generates new RSA key pairs for each session
- **User-Friendly TUI**: Clean terminal interface with flexible input parsing
- **XML-like Tags**: Secrets and keys are wrapped in easy-to-identify tags
- **Tolerance for Formatting Errors**: Handles minor typos in XML tags gracefully
- **Graceful Shutdown**: Supports 'q' or Ctrl+C to quit at any point

## Installation

1. Ensure you have Go installed (version 1.21 or later)
2. Clone this repository:
   ```
   git clone https://github.com/yourusername/secret_share.git
   cd secret_share
   ```
3. Build the application:
   ```
   go build -o secret_share cmd/secret_share/main.go
   ```
4. (Optional) Add the binary to your PATH:
   ```
   sudo cp secret_share /usr/local/bin/
   ```

## Usage

### Receiver Session

```
$: secret_share

Are you [s]ending or [r]eceiving a secret? 
> r

Here's the public key to share with the person sending you the secret:
<secret_share_key>ONE_TIME_PUBLIC_KEY</secret_share_key>

Enter the secret from the other person (should start and end with <secret_share_secret>): 
> <secret_share_secret>CIPHERTEXT</secret_share_secret>

Here's your secret: 12345
```

### Sender Session

```
$: secret_share

Are you [s]ending or [r]eceiving a secret? 
> s

Enter the secret key from the other person (should start and end with <secret_share_key>): 
> <secret_share_key>ONE_TIME_PUBLIC_KEY</secret_share_key>

Enter the secret you want to share: 
> ●●●●●●●●●

Here's the secret you can send back. Only they can decrypt it and see the secret:
<secret_share_secret>CIPHERTEXT</secret_share_secret>
```

## Security Implementation

SecretShare uses a well-established hybrid encryption approach:

- **RSA-OAEP**: 2048-bit RSA keys with Optimal Asymmetric Encryption Padding for secure key exchange
- **AES-GCM**: 256-bit AES encryption in Galois/Counter Mode for authenticated encryption of the secret
- **Secure Random Generation**: All keys and nonces are generated using Go's crypto/rand package
- **One-Time Keys**: New RSA key pairs are generated for each session to prevent reuse attacks
- **Memory Safety**: Private keys exist only in memory and are never written to disk

## Input Flexibility

SecretShare accepts various forms of input for user convenience:

- **Role Selection**: Accepts 's', 'S', '[s]', 'send', 'sender' for sending and 'r', 'R', '[r]', 'receive', 'receiver', 'recv' for receiving
- **Quit Commands**: Accepts 'q', 'Q', '[q]', 'quit', 'exit' at any prompt
- **XML Tag Tolerance**: Handles minor formatting errors in the XML-like tags

## License

This project is licensed under the MIT License - see the LICENSE file for details.
