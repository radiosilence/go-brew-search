#!/bin/bash

set -e

echo "🔨 Building go-brew-search..."

# Get the directory of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Change to the project directory
cd "$DIR"

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download

# Build the binary
echo "🏗️  Compiling..."
go build -o brew-search cmd/main.go

# Make it executable
chmod +x brew-search

echo "✅ Build complete! Binary created: ./brew-search"
echo ""
echo "To install system-wide, run:"
echo "  sudo cp brew-search /usr/local/bin/"
echo ""
echo "Or add to PATH:"
echo "  export PATH=\"$DIR:\$PATH\""