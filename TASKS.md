# Task Commands Quick Reference

## 🚀 Getting Started

### First Time Setup
```bash
# Install Task if not already installed
./install-task.sh

# Download dependencies and build
task deps
task build
```

### Installation Options
```bash
# Install to /usr/local/bin (requires sudo)
task install

# Install to ~/.local/bin (user directory)
task install-local
```

## 🏃 Common Tasks

### Build & Run
```bash
# Build the binary
task build

# Build and run immediately
task run

# Run interactive demo
task demo
```

### Development
```bash
# Format code
task fmt

# Run linters
task lint

# Run tests
task test

# Run tests with coverage report
task test-coverage

# Run all checks (fmt, lint, test)
task check

# Watch for changes and auto-rebuild
task dev
```

### Maintenance
```bash
# Clean build artifacts
task clean

# Clear package cache
task clean-cache

# Update dependencies
task tidy
```

### Release
```bash
# Build binaries for all platforms
task release
```

## 📋 All Available Tasks

| Task | Description |
|------|-------------|
| `task` | Show all available tasks |
| `task build` | Build the brew-search binary |
| `task install` | Install to /usr/local/bin |
| `task install-local` | Install to ~/.local/bin |
| `task run` | Build and run brew-search |
| `task demo` | Run interactive demo |
| `task test` | Run tests |
| `task test-coverage` | Run tests with coverage |
| `task fmt` | Format Go code |
| `task lint` | Run linters |
| `task check` | Run all checks |
| `task dev` | Watch mode with auto-rebuild |
| `task deps` | Download dependencies |
| `task tidy` | Tidy go.mod |
| `task clean` | Clean build artifacts |
| `task clean-cache` | Clear package cache |
| `task release` | Build release binaries |

## 💡 Tips

- Use `task --list` to see all available tasks with descriptions
- Tasks automatically handle dependencies (e.g., `task run` will build first)
- The `dev` task requires `fswatch` (install with `brew install fswatch`)
- Release binaries are created in the `dist/` directory