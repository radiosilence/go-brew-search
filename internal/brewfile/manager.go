package brewfile

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/go-brew-search/internal/api"
)

type Manager struct {
	path string
}

func New(path string) *Manager {
	return &Manager{
		path: path,
	}
}

// LoadExisting loads existing packages from Brewfile
func (m *Manager) LoadExisting() (map[string]bool, error) {
	existing := make(map[string]bool)

	file, err := os.Open(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			// Brewfile doesn't exist yet, that's okay
			return existing, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse brew formula
		if strings.HasPrefix(line, "brew ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				pkg := strings.Trim(parts[1], `"'`)
				existing[pkg] = true
			}
		}

		// Parse cask
		if strings.HasPrefix(line, "cask ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				pkg := strings.Trim(parts[1], `"'`)
				existing[pkg] = true
			}
		}

		// Parse tap
		if strings.HasPrefix(line, "tap ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				tap := strings.Trim(parts[1], `"'`)
				existing[tap] = true
			}
		}
	}

	return existing, scanner.Err()
}

// AddPackages adds new packages to the Brewfile
func (m *Manager) AddPackages(packages []api.Package) error {
	// Ensure directory exists
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open file in append mode
	file, err := os.OpenFile(m.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if file is empty or needs a newline
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > 0 {
		// Add newline if file doesn't end with one
		buf := make([]byte, 1)
		if _, err := file.ReadAt(buf, stat.Size()-1); err == nil && buf[0] != '\n' {
			file.WriteString("\n")
		}
		file.WriteString("\n")
	}

	// Add comment header
	file.WriteString(fmt.Sprintf("# Added by go-brew-search on %s\n", 
		time.Now().Format("2006-01-02 15:04:05")))

	// Add packages
	for _, pkg := range packages {
		var line string
		if pkg.Type == "cask" {
			line = fmt.Sprintf("cask \"%s\"", pkg.Token)
		} else {
			line = fmt.Sprintf("brew \"%s\"", pkg.Token)
		}

		if pkg.Description != "" {
			// Add comment with description
			desc := pkg.Description
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			line += fmt.Sprintf(" # %s", desc)
		}

		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// RunBundle runs brew bundle command
func (m *Manager) RunBundle() error {
	cmd := exec.Command("brew", "bundle", "--file", m.path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
}