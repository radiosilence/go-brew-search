#!/bin/bash

set -e

echo "🔍 Checking for Task installation..."

if command -v task &> /dev/null; then
    echo "✅ Task is already installed ($(task --version))"
    exit 0
fi

echo "📦 Task not found. Installing..."

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    if command -v brew &> /dev/null; then
        echo "🍺 Installing Task via Homebrew..."
        brew install go-task/tap/go-task
    else
        echo "❌ Homebrew not found. Please install Homebrew first:"
        echo "   /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
        exit 1
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    echo "🐧 Installing Task for Linux..."
    
    # Detect architecture
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7l) ARCH="arm" ;;
        *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    
    # Download and install
    VERSION=$(curl -s https://api.github.com/repos/go-task/task/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
    curl -sL "https://github.com/go-task/task/releases/download/v${VERSION}/task_linux_${ARCH}.tar.gz" | tar -xz -C /tmp
    sudo mv /tmp/task /usr/local/bin/
    
    echo "✅ Task installed successfully!"
else
    echo "❌ Unsupported operating system: $OSTYPE"
    echo "   Please install Task manually: https://taskfile.dev/installation/"
    exit 1
fi

echo "✅ Task installation complete!"
echo "📝 You can now use 'task' to build the project"