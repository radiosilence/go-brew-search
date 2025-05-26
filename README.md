# 🍺 go-brew-search

A blazing-fast, interactive Homebrew package search tool with local caching. Say goodbye to slow `brew search` commands!

## ✨ Features

- **⚡ Lightning Fast**: Uses local caching to provide instant search results
- **🔍 Fuzzy Search**: Interactive fzf-style interface for easy package discovery
- **📊 Rich Display**: Shows package names, versions, and descriptions in a clean table format
- **✅ Multi-Select**: Select multiple packages at once with TAB
- **🎯 Smart Integration**: Automatically adds selected packages to your `~/Brewfile`
- **🚀 Auto-Install**: Runs `brew bundle` after selection to install packages
- **📝 Status Tracking**: Shows which packages are already in your Brewfile
- **🎨 Beautiful UI**: Uses emojis and colors for better visual feedback

## 📦 Installation

### Prerequisites

- Go 1.21 or later
- Homebrew installed
- A terminal that supports Unicode (for emojis)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/go-brew-search.git
cd go-brew-search

# Build the binary
./build.sh

# Option 1: Install system-wide
sudo cp brew-search /usr/local/bin/

# Option 2: Add to PATH
export PATH="$PWD:$PATH"
```

### Quick Install (using go install)

```bash
go install github.com/yourusername/go-brew-search/cmd@latest
mv $GOPATH/bin/cmd $GOPATH/bin/brew-search
```

## 🚀 Usage

Simply run:

```bash
brew-search
```

### Interactive Controls

- **Type** to search packages
- **↑↓** to navigate
- **TAB** to select/deselect packages
- **ENTER** to confirm selection
- **ESC** to cancel

### Package Icons

- 📦 Formula (regular Homebrew packages)
- 🍺 Cask (GUI applications)
- ✓ Already in Brewfile

## 🔧 How It Works

1. **Caching**: On first run, fetches all package data from Homebrew's API and caches it locally
2. **Search**: Provides instant fuzzy search through the cached data
3. **Selection**: Shows an interactive interface with package details
4. **Brewfile Management**: Adds selected packages to `~/Brewfile` with proper formatting
5. **Installation**: Automatically runs `brew bundle` to install the selected packages

## 📁 Cache Management

The tool caches Homebrew package data in `~/.cache/go-brew-search/`. The cache expires after 24 hours.

To clear the cache manually:

```bash
rm -rf ~/.cache/go-brew-search
```

## ⚙️ Configuration

Currently, the tool uses sensible defaults:

- **Cache Location**: `~/.cache/go-brew-search/`
- **Cache TTL**: 24 hours
- **Brewfile Location**: `~/Brewfile`

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development

```bash
# Run tests
go test ./...

# Build with verbose output
go build -v -o brew-search cmd/main.go

# Run with Go
go run cmd/main.go
```

## 📝 License

MIT License - feel free to use this tool however you like!

## 🐛 Known Issues

- The tool requires a terminal that supports Unicode for emoji display
- Initial cache population might take a few seconds on first run
- Some casks might not have complete metadata

## 🎯 Motivation

The built-in `brew search` command is notoriously slow because it queries the API every time. This tool solves that problem by:

1. Caching all package data locally
2. Providing instant fuzzy search
3. Offering a better UI with multi-select
4. Automating the Brewfile workflow

Never wait for `brew search` again! 🚀