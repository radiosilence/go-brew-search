#!/bin/bash

# Demo script for go-brew-search

set -e

echo "🍺 go-brew-search Demo"
echo "====================="
echo ""

# Build if not already built
if [ ! -f "./brew-search" ]; then
    echo "🔨 Building brew-search..."
    ./build.sh
    echo ""
fi

echo "📝 This demo will show you how to use go-brew-search"
echo ""
echo "Features demonstrated:"
echo "  - Fast fuzzy search through all Homebrew packages"
echo "  - Multi-select capability (use TAB)"
echo "  - Shows which packages are already in your Brewfile"
echo "  - Automatically updates ~/Brewfile"
echo "  - Runs brew bundle to install selected packages"
echo ""
echo "Press ENTER to start the interactive search..."
read

# Run the tool
./brew-search

echo ""
echo "✨ Demo complete!"
echo ""
echo "💡 Tips:"
echo "  - The first run caches all packages (may take a few seconds)"
echo "  - Subsequent runs use the cache and are instant"
echo "  - Cache expires after 24 hours"
echo "  - Clear cache with: rm -rf ~/.cache/go-brew-search"
echo ""
echo "📦 To install system-wide:"
echo "  sudo cp brew-search /usr/local/bin/"