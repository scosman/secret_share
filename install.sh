#!/bin/bash

# SecretShare installer script
# Detects platform and architecture, downloads the appropriate binary, and installs it

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print functions
print_error() {
    echo -e "${RED}Error: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${YELLOW}$1${NC}"
}

# Detect platform and architecture
detect_platform() {
    # Check if we're on Windows (any Windows environment)
    if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ -n "$WINDIR" ]]; then
        OS="windows"
        ARCH=$(uname -m 2>/dev/null || echo "x86_64")
        BINARY_NAME="secret_share.exe"
    else
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        ARCH=$(uname -m)
        BINARY_NAME="secret_share"
    fi
    
    # Normalize architecture names
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    # Normalize OS names
    case $OS in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        windows)
            OS="windows"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}-${ARCH}"
    
    # Map platform to release file name (must match release.yml)
    case $PLATFORM in
        linux-amd64)
            RELEASE_NAME="Linux-amd64"
            ;;
        linux-arm64)
            RELEASE_NAME="Linux-arm64"
            ;;
        darwin-amd64)
            RELEASE_NAME="MacOS-Intel-x64"
            ;;
        darwin-arm64)
            RELEASE_NAME="MacOS-Apple-Silicon"
            ;;
        windows-amd64)
            RELEASE_NAME="Windows-amd64"
            ;;
        windows-arm64)
            RELEASE_NAME="Windows-arm64"
            ;;
        *)
            print_error "Unsupported platform: $PLATFORM"
            exit 1
            ;;
    esac
}

# Get the latest release URL
get_latest_release_url() {
    # Use GitHub API to get the latest release
    LATEST_RELEASE=$(curl -s https://api.github.com/repos/scosman/secret_share/releases/latest)
    
    # Extract the tag name
    TAG_NAME=$(echo "$LATEST_RELEASE" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$TAG_NAME" ]; then
        print_error "Could not determine the latest release tag"
        exit 1
    fi
    
    # Construct the download URL
    DOWNLOAD_URL="https://github.com/scosman/secret_share/releases/download/${TAG_NAME}/secret_share-${RELEASE_NAME}.zip"
}

# Download and install the binary
install_binary() {
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    print_info "Downloading SecretShare for $PLATFORM..."
    
    # Download the zip file
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o secret_share.zip "$DOWNLOAD_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -q -O secret_share.zip "$DOWNLOAD_URL"
    else
        print_error "Neither curl nor wget found. Please install one of them and try again."
        exit 1
    fi
    
    # Extract the zip file
    if command -v unzip >/dev/null 2>&1; then
        unzip -q secret_share.zip
    else
        print_error "unzip command not found. Please install unzip and try again."
        exit 1
    fi
    
    # Determine installation directory
    if [ "$OS" = "windows" ]; then
        # On Windows, use the standard user programs directory
        # %LOCALAPPDATA%\Programs is the Windows standard for user-installed programs
        LOCALAPPDATA="${LOCALAPPDATA:-$HOME/AppData/Local}"
        INSTALL_DIR="$LOCALAPPDATA/Programs"
        mkdir -p "$INSTALL_DIR"
        
        # Convert to Windows path format for display
        WIN_INSTALL_DIR=$(cygpath -w "$INSTALL_DIR" 2>/dev/null || echo "$INSTALL_DIR")
    else
        # Unix-like systems (Linux, macOS)
        # Try /usr/local/bin first
        if [ -w "/usr/local/bin" ]; then
            INSTALL_DIR="/usr/local/bin"
        else
            # Fallback to ~/.local/bin
            INSTALL_DIR="$HOME/.local/bin"
            mkdir -p "$INSTALL_DIR"
        fi
    fi
    
    # Move binary to installation directory
    print_info "Installing SecretShare to $INSTALL_DIR..."
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    
    # Make binary executable (not needed on Windows)
    if [ "$OS" != "windows" ]; then
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Clean up
    rm -rf "$TEMP_DIR"
    
    print_success "SecretShare installed successfully to $INSTALL_DIR/$BINARY_NAME"
    
    # Provide platform-specific instructions
    if [ "$OS" = "windows" ]; then
        echo "SecretShare has been installed to the standard Windows user programs directory."
        echo ""
        echo "Option 1 - Run from programs directory (no setup required):"
        echo "  cd \"${WIN_INSTALL_DIR:-$INSTALL_DIR}\""
        echo "  .\\secret_share.exe"
        echo ""
        echo "Option 2 - Add to PATH for global access (PowerShell):"
        echo "  1. Open PowerShell and run:"
        echo "     \$env:PATH += \";${WIN_INSTALL_DIR:-$INSTALL_DIR}\""
        echo "     [Environment]::SetEnvironmentVariable(\"PATH\", \$env:PATH, \"User\")"
        echo "  2. Close and reopen PowerShell/Command Prompt/Windows Terminal"
        echo ""
        echo "Option 3 - Add to PATH via System Properties GUI:"
        echo "  1. Press Win+X → System → Advanced system settings"
        echo "  2. Click 'Environment Variables'"
        echo "  3. Under 'User variables', edit 'Path'"
        echo "  4. Add: ${WIN_INSTALL_DIR:-$INSTALL_DIR}"
        echo "  5. Close and reopen your shell (PowerShell/CMD/Windows Terminal)"
    else
        echo "You can now run 'secret_share' from your terminal."
        echo "If the command is not found you may need to add it to your PATH:"
        echo "  echo 'export PATH=\$PATH:\$HOME/.local/bin' >> ~/.bashrc"
        echo "  source ~/.bashrc"
    fi
}

# Main execution
main() {
    echo "SecretShare Installer"
    echo "===================="
    
    detect_platform
    get_latest_release_url
    install_binary
}

# Run main function
main "$@"
