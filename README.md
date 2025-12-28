<p align="center">
        <picture>
            <img width="323" height="167" alt="SecretShare logo" src="https://github.com/user-attachments/assets/1aa94780-92b5-4a90-9713-7abd172f4e1c" />
        </picture>
</p>
<h3 align="center">
    Share Secrets/Passwords Securely with with a CLI — in 15 Seconds
</h3>
<p align="center">
  <a href="#how-it-works"><strong>How it Works</strong></a> •
  <a href="#overview"><strong>Overview: Security and Use</strong></a> • 
  <a href="#install"><strong>Install</strong></a> • 
  <a href="#technical-details"><strong>Technical Details</strong></a> • 
  <a href="#usability"><strong>Usability</strong></a> • 
  <a href="#demo-gif"><strong>Demo</strong></a>
</p>

### Quick Start: Install and Run

```bash
curl -fsSL https://raw.githubusercontent.com/scosman/secret_share/main/install.sh | sh
secret_share
```

**Windows Users**: run this in Windows Subsystem for Linux (WSL), not PowerShell.

Or check out all [install options](#install).

## How it works

The whole process takes about 15 seconds:

1. The receiver runs `secret_share`, it generates a one-time public key, they send the key to the sender
2. The sender runs `secret_share`, pastes in the public key from the receiver, type/paste a secret, it generates an encrypted response they send back
3. The receiver pastes in the encrypted response and sees the secret

[[Demo GIF]](#demo-gif)

<img width="1416" height="677" alt="asdf" src="https://github.com/user-attachments/assets/5825fce2-46c5-4598-b929-ac94d86fe176" />

## Overview

### It's Easy

1. Run the app, send one message, paste the response, done. 
2. Send the messages using any chat app: you don't need to trust the communication channel since it never sees the private key.
3. No complicated security questions, just smart defaults
4. Cross platform: available for Mac, Windows and Linux

### It's Secure

1. Private key never leaves the senders device
2. Private key is never written to a file or shown on screen, it is only kept in memory
3. New random keys for every session
4. No servers, no one to trust 
5. Uses standard, strong, boring encryption: RSA-OAEP and AES-GCM 
6. Uses golang's standard crypto package (audited)
7. No dependencies except for offical Google go packages (sys & term)
8. Open source: build yourself or use public builds with checksums
9. Tiny: read all the [crypto code](core/crypto.go) in about 1 minute or the whole app in about 5 minutes

## Install

### Automatic Installer

You can install SecretShare with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/scosman/secret_share/main/install.sh | sh
```

**Important:** Windows users should run this in Windows Subsystem for Linux (WSL), not PowerShell.

Obviously feel free to download and read the installer script before running! It's downloading the latest offical release from Github.

### Install from Source

Run the following to install secret_share from source. You'll need to have [golang installed](https://go.dev/doc/install).

```bash
go install github.com/scosman/secret_share/cmd/secret_share@latest
```

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

Security note: secret_send does nothing to verify the identity of the person you're sharing with. That is similar to tools which use secret links, but not as robust as something like PGP or Keybase. The tradeoff is ease of setup and complexity.

Being an interactive CLI and not having arguments is an intentional security+usability choice. Other tools like [age](https://github.com/FiloSottile/age) allow you to generate private key files, but also make it the user's responsibility to securely manage those keys (keeping track of them, deleting them, time-based expiration, etc). SecretSend keeps it simple: no one ever sees the private key, it's never written to disk, and it's cleared from memory as soon as the app ends. This makes it great for one-time secret sharing between people. If you want long-term secret management with long lived keys, check out [age](https://github.com/FiloSottile/age).

## Usability

 - User friendly TUI: clear questions, instructions, and errors
 - Clipboard support: it automatically copies the keys/encrypted-secret to clipboard at the appropriate time
 - Flexible parsing: don't sweat it if you paste a few extra characters
 - No args: interactive terminal UI walks you through steps, no need to memorize args, no multi-step processes
 - No options/settings: just secure defaults
 - User isn't responsible for security: we don't show them the private key, there's no key files to delete, we don't ask them to choose key-length or algorithms.

## Demo GIF

![screen cast](https://github.com/user-attachments/assets/0d2f2524-38a8-4455-9e65-23c7247d67f0)

## License

This project is licensed under the MIT License - see the LICENSE.txt file for details.
