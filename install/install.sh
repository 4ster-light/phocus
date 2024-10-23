#!/bin/bash

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go before continuing."
    echo "Visit https://golang.org/doc/install for installation instructions."
    exit 1
fi

# Install phocus
echo "Installing phocus..."
go install github.com/4ster-light/phocus@latest

# Check if installation was successful
if [ $? -ne 0 ]; then
    echo "Installation failed"
    exit 1
fi

# Get the binary path from GOPATH
GOPATH=$(go env GOPATH)
SOURCE_BINARY="$GOPATH/bin/phocus"

# Ensure the binary exists
if [ ! -f "$SOURCE_BINARY" ]; then
    echo "Binary not found after installation"
    exit 1
fi

# Try to move to /usr/local/bin first (preferred), fall back to /usr/bin if needed
if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
    TARGET_DIR="/usr/local/bin"
elif [ -d "/usr/bin" ] && [ -w "/usr/bin" ]; then
    TARGET_DIR="/usr/bin"
else
    echo "Error: Need sudo privileges to move binary to system directory"
    echo "Please run: sudo mv $SOURCE_BINARY /usr/local/bin/phocus"
    exit 1
fi

# Move the binary
if mv "$SOURCE_BINARY" "$TARGET_DIR/phocus"; then
    echo "Phocus has been successfully installed to $TARGET_DIR/phocus"
else
    echo "Failed to move binary. Please run with sudo:"
    echo "sudo mv $SOURCE_BINARY $TARGET_DIR/phocus"
    exit 1
fi

# Make it executable
chmod +x "$TARGET_DIR/phocus"
echo "No you can run the program from anywhere in your system with: sudo phocus"
