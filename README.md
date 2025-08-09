<p align="center">
        <picture>
            <img width="292" height="109" alt="SecretShare logo" src="https://github.com/user-attachments/assets/f66fec47-2e54-4f3c-aced-049c40881f2e" />
        </picture>
</p>

### SecretShare: Easy, secure one time secret sharing CLI

## How it works

The whole process takes about 15 seconds:

1. The receiver runs secret_share and it generates a one-time public key which send to the sender
2. The sender runs secret_share, pastes in the public key from the receiver, types the secret, and gets an encrypted response they send back
3. The receiver pastes in the encrypted response and sees the secret

<img width="1669" height="694" alt="flow chart" src="https://github.com/user-attachments/assets/9a1f0b4f-915a-4bd4-80e3-14a27066f58e" />

## Overview

### It's Easy

1. Run the app, send one message, paste the response, done. 
2. Send the messages using any chat app: you don't need to trust the communication channel since it never sees the private key.
3. No complicated security questions, just smart defaults
4. Cross platform: available for Mac, Windows and Linux

### It's Secure

1. Private key never leaves the senders device
2. Private key is never written to a file or shown on screen, they are only kept in memory
3. New random keys for every session
4. No servers, no one to trust 
5. Uses standard, strong, boring encryption: RSA-OAEP and AES-GCM 
6. Uses golang's standard crypto package (audited)
7. No dependencies except for offical Google go packages (sys & term)
8. Open source: build yourself or use public builds with checksums
9. Tiny: read all the [crypto code](core/crypto.go) in about 1 minute or the whole app in about 5 minutes

## Install/Build with Go

If you have go installed just run:

```bash
go install github.com/scosman/secret_share/cmd/secret_share@latest
```

### Mac/Linux Installer

For macOS and Linux, you can install SecretShare with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/scosman/secret_share/main/install.sh | sh
```

Feel free to download and read the installer script!

### Manual Download

You can download the latest release from [GitHub releases](https://github.com/scosman/secret_share/releases). 

### Build Integrity

Builds are created on public [Github Actions](https://github.com/scosman/secret_share/actions/workflows/release.yml). The build logs include checksums you can validate.

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

Using hybrid RSA+AES encryption allows us to share secrets of any length. RSA alone can only encrypt short payloads.

Security note: secret_send does nothing to verify the identity of the person you're sharing with. That's similar to tools that use secret links, but not as robust as something like PGP or Keybase. The tradeoff is ease of setup and complexity.

## Usability

 - User friendly TUI: clear questions, instructions and errors
 - Clipboard support: it automatically copies the keys/encrypted-secret to clipboard at the appropriate time (MacOS and Linux)
 - Flexible parsing: don't sweat it if you paste a few extra characters
 - No args: interactive terminal UI walks you through steps, no need to memorize args
 - No options/settings

## Demo GIF

![screen cast](https://github.com/user-attachments/assets/0d2f2524-38a8-4455-9e65-23c7247d67f0)

## License

This project is licensed under the MIT License - see the LICENSE.txt file for details.
