#!/bin/bash

set -e

# Determine the architecture
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Get the latest release version from GitHub
LATEST_VERSION=$(curl -s https://api.github.com/repos/DepsHubHQ/depshub/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [[ -z "$LATEST_VERSION" ]]; then
    echo "Failed to fetch the latest version. Please check your internet connection or GitHub API status."
    exit 1
fi

# Define the binary URL
BINARY_URL="https://github.com/DepsHubHQ/depshub/releases/download/$LATEST_VERSION/depshub-linux-$ARCH"

# Define the target installation path
INSTALL_DIR="$HOME/.local/bin"
INSTALL_PATH="$INSTALL_DIR/depshub"

# Ensure the installation directory exists
mkdir -p "$INSTALL_DIR"

# Download the binary
echo "Downloading DepsHub $LATEST_VERSION for $ARCH..."
curl -L -o depshub "$BINARY_URL"

# Make the binary executable
chmod +x depshub

# Move the binary to the installation directory
echo "Installing DepsHub to $INSTALL_DIR..."
mv depshub "$INSTALL_PATH"

# Verify installation
if [[ -x "$INSTALL_PATH" ]]; then
    echo "DepsHub installed successfully!"
    echo "Make sure $INSTALL_DIR is in your PATH if it's not already."
    echo "Run 'depshub --help' to get started."
else
    echo "Installation failed. Please check your permissions or try again."
    exit 1
fi
