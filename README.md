<p align="center">
        <picture>
            <img width="292" height="109" alt="SecretShare logo" src="https://github.com/user-attachments/assets/f66fec47-2e54-4f3c-aced-049c40881f2e" />
        </picture>
</p>

### SecretShare: An easy and secure way to share secrets

## How it works

The whole process takes about 15 seconds:

1. The receiver runs secret_share and it generates a one-time public key they can send to the sender
2. The sender runs secret_share, pastes in the public key from the receiver, types the secret, and gets an encrypted response they can send back
3. The receiver pastes in the encrypted response and sees the secret

![output_420](https://github.com/user-attachments/assets/0d2f2524-38a8-4455-9e65-23c7247d67f0)

## Overview

### It's Easy

1. Run the app, send one message, paste the response, done. 
2. Send the messages using your normal chat app: you don't need to trust the communication channel, they never see the private keys.
3. No complicated security questions, just smart defaults
4. Cross platform: available for Mac, Windows and Linux

### It's Secure

1. Private key never leaves the senders device
2. Private key is never written to a file or shown on screen, they are only kept in memory
3. New random keys for every session
4. No servers, no one to trust 
5. Uses standard, strong, boring encryption: RSA-OAEP and AES-GCM 
6. Uses golang's standard crypto package (formally audited)
7. No dependencies except golang.org/x/term (Google maintained)
8. Open source: build yourself or use public builds from Github Actions with checksums
9. Tiny: read all the code in about 5 minutes

## Technical Details

SecretShare is a golang command-line tool. They encryption flow works as follows:

1. The receiver generates a one-time RSA key pair
2. The receiver shares only the public key with the sender
3. The sender uses the public key to encrypt a randomly generated AES key
4. The sender uses the AES key to encrypt the actual secret with AES-GCM
5. The sender shares the encrypted data with the receiver
6. The receiver uses their private key to decrypt the AES key, then decrypts the secret
7. The app ends, removing the keys from memory

The private key never leaves the receiver's machine and is never exposed to the communication channel.

Using hybrid encryption allows us to share secrets of any length. RSA can only encrypt short data.

## Usability

 - User friendly TUI: clear questions, instructions and errors
 - Clipboard support: it automatically copies the keys/encrypted-secret to clipboard at the appropriate time (MacOS and Linux)
 - Flexible parsing: don't sweat it if you paste a few extra characters
 - No args: interactive terminal UI, no need to memorize args
 - No options/settings

## Installation

1. Ensure you have Go installed (version 1.21 or later)
2. Clone this repository:
   ```bash
   git clone git@github.com:scosman/secret_share.git
   cd secret_share
   ```
3. Build the application:
   ```bash
   go build -o secret_share cmd/secret_share/main.go
   ```
4. (Optional) Add the binary to your PATH:
   ```bash
   sudo cp secret_share /usr/local/bin/
   ```
5. Run the app
   ```bash
   ./secret_share
   ```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
