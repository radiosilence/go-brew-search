package main

import (
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/user/go-brew-search/internal/api"
)

func main() {
	fmt.Println("🧪 Testing UI Display")
	fmt.Println("====================")
	
	// Create test packages
	packages := []api.Package{
		{Token: "btop", Type: "formula", Version: "4.0.2", Description: "Resource monitor that shows usage and stats for processor, memory, disks, network and processes"},
		{Token: "htop", Type: "formula", Version: "3.3.0", Description: "Improved top (interactive process viewer)"},
		{Token: "mpv", Type: "formula", Version: "0.37.0", Description: "Media player based on MPlayer and mplayer2"},
		{Token: "iterm2", Type: "cask", Version: "3.4.23", Description: "Terminal emulator as alternative to Apple's Terminal app"},
		{Token: "firefox", Type: "cask", Version: "121.0", Description: "Web browser"},
		{Token: "c", Type: "formula", Version: "0.14", Description: "Compile and execute C \"scripts\" in one go"},
		{Token: "ibkr", Type: "cask", Version: "10.13.0g", Description: "Trading software"},
	}
	
	existing := map[string]bool{
		"htop":    true,
		"firefox": true,
	}
	
	// Test different formatting approaches
	fmt.Println("\n1. Current format with dots:")
	fmt.Println("─────────────────────────────")
	testFormat(packages, existing, "%s %s %-30s · %-15s · %s")
	
	fmt.Println("\n2. Tab-separated format:")
	fmt.Println("────────────────────────")
	testFormat(packages, existing, "%s %s %-30s\t%-15s\t%s")
	
	fmt.Println("\n3. Double-space format:")
	fmt.Println("───────────────────────")
	testFormat(packages, existing, "%s %s %-30s  %-15s  %s")
	
	fmt.Println("\n4. Custom separator format:")
	fmt.Println("───────────────────────────")
	testFormat(packages, existing, "%s %s %-30s :: %-15s :: %s")
	
	// Test fuzzy finder display
	fmt.Println("\n\n🔍 Testing Fuzzy Finder Display")
	fmt.Println("================================")
	fmt.Println("(Press Ctrl+C to exit)\n")
	
	testFuzzyFinder(packages, existing)
}

func testFormat(packages []api.Package, existing map[string]bool, format string) {
	for _, pkg := range packages {
		statusIcon := "  "
		if existing[pkg.Token] {
			statusIcon = "✅"
		}
		
		typeIcon := "⚡"
		if pkg.Type == "cask" {
			typeIcon = "🖥️"
		}
		
		name := truncate(pkg.Token, 30)
		version := truncate(pkg.Version, 15)
		desc := pkg.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		
		line := fmt.Sprintf(format, statusIcon, typeIcon, name, version, desc)
		fmt.Println(line)
	}
}

func testFuzzyFinder(packages []api.Package, existing map[string]bool) {
	// Create display items
	items := make([]string, len(packages))
	for i, pkg := range packages {
		statusIcon := "  "
		if existing[pkg.Token] {
			statusIcon = "✅"
		}
		
		typeIcon := "⚡"
		if pkg.Type == "cask" {
			typeIcon = "🖥️"
		}
		
		name := truncate(pkg.Token, 30)
		version := truncate(pkg.Version, 15)
		desc := pkg.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		
		// Test with double spaces
		items[i] = fmt.Sprintf("%s %s %-30s  %-15s  %s", statusIcon, typeIcon, name, version, desc)
	}
	
	// Show fuzzy finder
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPromptString("Test search: "),
	)
	
	if err != nil {
		fmt.Printf("Error or cancelled: %v\n", err)
		return
	}
	
	fmt.Printf("\nSelected: %s\n", items[idx])
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}