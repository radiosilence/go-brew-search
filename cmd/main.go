package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/user/go-brew-search/internal/api"
	"github.com/user/go-brew-search/internal/brewfile"
	"github.com/user/go-brew-search/internal/cache"
	"github.com/user/go-brew-search/internal/ui"
)

func main() {
	// Initialize cache directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("❌ Failed to get home directory:", err)
	}

	cacheDir := filepath.Join(homeDir, ".cache", "go-brew-search")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Fatal("❌ Failed to create cache directory:", err)
	}

	// Initialize components
	cacheManager := cache.New(cacheDir, 24*time.Hour)
	apiClient := api.New(cacheManager)
	brewfileManager := brewfile.New(filepath.Join(homeDir, "Brewfile"))

	// Load existing Brewfile packages
	existing, err := brewfileManager.LoadExisting()
	if err != nil {
		log.Printf("⚠️  Warning: Could not load Brewfile: %v", err)
		existing = make(map[string]bool)
	}

	// Fetch packages
	fmt.Println("🔄 Fetching Homebrew packages...")
	packages, err := apiClient.FetchAllPackages()
	if err != nil {
		log.Fatal("❌ Failed to fetch packages:", err)
	}

	fmt.Printf("✅ Loaded %d packages\n", len(packages))

	// Show interactive UI
	selected, err := ui.ShowPackageSelector(packages, existing)
	if err != nil {
		log.Fatal("❌ Error in package selector:", err)
	}

	if len(selected) == 0 {
		fmt.Println("👋 No packages selected")
		return
	}

	// Filter out already installed packages
	newPackages := []api.Package{}
	for _, pkg := range selected {
		if !existing[pkg.Token] {
			newPackages = append(newPackages, pkg)
		}
	}

	if len(newPackages) == 0 {
		fmt.Println("✅ All selected packages are already in Brewfile")
		return
	}

	// Add new packages to Brewfile
	fmt.Printf("📝 Adding %d new packages to Brewfile...\n", len(newPackages))
	if err := brewfileManager.AddPackages(newPackages); err != nil {
		log.Fatal("❌ Failed to update Brewfile:", err)
	}

	// Run brew bundle
	fmt.Println("🚀 Running brew bundle...")
	if err := brewfileManager.RunBundle(); err != nil {
		log.Fatal("❌ Failed to run brew bundle:", err)
	}

	fmt.Println("✨ Done!")
}