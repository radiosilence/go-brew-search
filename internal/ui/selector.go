package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/user/go-brew-search/internal/api"
)

func ShowPackageSelector(packages []api.Package, existing map[string]bool) ([]api.Package, error) {
	// Create a copy and sort packages by token length (shorter = more likely to be searched)
	sortedPackages := make([]api.Package, len(packages))
	copy(sortedPackages, packages)
	
	sort.Slice(sortedPackages, func(i, j int) bool {
		// First by token length
		if len(sortedPackages[i].Token) != len(sortedPackages[j].Token) {
			return len(sortedPackages[i].Token) < len(sortedPackages[j].Token)
		}
		// Then alphabetically
		return sortedPackages[i].Token < sortedPackages[j].Token
	})
	
	// Prepare display items
	items := make([]string, len(sortedPackages))
	for i, pkg := range sortedPackages {
		// Status indicators
		var statusIcon string
		if existing[pkg.Token] {
			statusIcon = "✅" // Already installed
		} else {
			statusIcon = "  "
		}
		
		// Package type icon
		var typeIcon string
		if pkg.Type == "cask" {
			typeIcon = "🖥️ " // GUI app
		} else {
			typeIcon = "⚡" // CLI tool
		}
		
		name := pkg.Token
		if pkg.FullName != "" && pkg.FullName != pkg.Token {
			name = fmt.Sprintf("%s (%s)", pkg.Token, pkg.FullName)
		}
		
		desc := pkg.Description
		if desc == "" {
			desc = "—"
		}
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		
		version := pkg.Version
		if version == "" {
			version = "unknown"
		}
		if len(version) > 20 {
			version = version[:20] + "..."
		}
		
		// Format with clear visual separation using box drawing characters
		nameStr := truncate(name, 32)
		versionStr := truncate(version, 12)
		descStr := desc
		if len(descStr) > 60 {
			descStr = descStr[:57] + "..."
		}
		
		items[i] = fmt.Sprintf("%s %s %-32s │ %-12s │ %s", 
			statusIcon,
			typeIcon,
			nameStr,
			versionStr,
			descStr,
		)
	}
	
	// Create searchable strings that prioritize exact matches
	searchStrings := make([]string, len(sortedPackages))
	for i, pkg := range sortedPackages {
		// Build search string with weighted components
		var parts []string
		
		// Exact token match (highest priority - repeat 10 times)
		for j := 0; j < 10; j++ {
			parts = append(parts, pkg.Token)
		}
		
		// Add full name if different (medium priority)
		if pkg.FullName != "" && pkg.FullName != pkg.Token {
			parts = append(parts, pkg.FullName, pkg.FullName)
		}
		
		// Add description (lowest priority)
		if pkg.Description != "" {
			parts = append(parts, pkg.Description)
		}
		
		searchStrings[i] = strings.Join(parts, " ")
	}
	
	// Show fuzzy finder with multi-select
	indices, err := fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			// Return the search string for scoring, but display shows items[i]
			return searchStrings[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			
			pkg := sortedPackages[i]
			var preview strings.Builder
			
			// Header with package name and type
			typeEmoji := "⚡"
			typeName := "Formula"
			if pkg.Type == "cask" {
				typeEmoji = "🖥️"
				typeName = "Cask"
			}
			
			preview.WriteString(fmt.Sprintf("%s %s\n", typeEmoji, pkg.Token))
			preview.WriteString(strings.Repeat("─", min(len(pkg.Token)+3, w)) + "\n\n")
			
			// Installation status
			if existing[pkg.Token] {
				preview.WriteString("✅ Already in Brewfile\n")
			} else {
				preview.WriteString("📦 Not in Brewfile\n")
			}
			
			// Package details
			preview.WriteString(fmt.Sprintf("📋 Type: %s\n", typeName))
			
			if pkg.Version != "" {
				preview.WriteString(fmt.Sprintf("🏷️  Version: %s\n", pkg.Version))
			}
			
			if pkg.FullName != "" && pkg.FullName != pkg.Token {
				preview.WriteString(fmt.Sprintf("📛 Full Name: %s\n", pkg.FullName))
			}
			
			// Description
			if pkg.Description != "" {
				preview.WriteString(fmt.Sprintf("\n📄 Description:\n%s\n", wordWrap(pkg.Description, w-2)))
			}
			
			// Homepage
			if pkg.Homepage != "" {
				preview.WriteString(fmt.Sprintf("\n🌐 Homepage:\n%s\n", pkg.Homepage))
			}
			
			// Installation command preview
			preview.WriteString(fmt.Sprintf("\n💻 Install command:\nbrew install %s\n", pkg.Token))
			
			return preview.String()
		}),
		fuzzyfinder.WithPromptString("🔍 Search packages: "),
		fuzzyfinder.WithHeader("\n  ⚡ Formula  🖥️  Cask  ✅ In Brewfile  │  TAB: Select  ENTER: Confirm  ESC: Cancel\n  ──────────────────────────────────────────────────────────────────────────────────────────────────\n"),
	)
	
	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return nil, nil // User cancelled
		}
		return nil, err
	}
	
	// Collect selected packages
	selected := make([]api.Package, len(indices))
	for i, idx := range indices {
		selected[i] = sortedPackages[idx]
	}
	
	return selected, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}
	
	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0
	
	for i, word := range words {
		wordLen := len(word)
		
		if i > 0 && lineLen+wordLen+1 > width {
			result.WriteString("\n")
			lineLen = 0
		} else if i > 0 {
			result.WriteString(" ")
			lineLen++
		}
		
		result.WriteString(word)
		lineLen += wordLen
	}
	
	return result.String()
}