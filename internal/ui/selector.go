package ui

import (
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/user/go-brew-search/internal/api"
)

func ShowPackageSelector(packages []api.Package, existing map[string]bool) ([]api.Package, error) {
	// Prepare display items
	items := make([]string, len(packages))
	for i, pkg := range packages {
		status := "  "
		if existing[pkg.Token] {
			status = "✓ "
		}
		
		typeIcon := "📦"
		if pkg.Type == "cask" {
			typeIcon = "🍺"
		}
		
		name := pkg.Token
		if pkg.FullName != "" && pkg.FullName != pkg.Token {
			name = fmt.Sprintf("%s (%s)", pkg.Token, pkg.FullName)
		}
		
		desc := pkg.Description
		if desc == "" {
			desc = "No description available"
		}
		if len(desc) > 60 {
			desc = desc[:60] + "..."
		}
		
		version := pkg.Version
		if version == "" {
			version = "unknown"
		}
		if len(version) > 20 {
			version = version[:20] + "..."
		}
		
		items[i] = fmt.Sprintf("%s%s %-40s %-25s %s", 
			status, 
			typeIcon,
			truncate(name, 40),
			truncate(version, 25),
			desc,
		)
	}
	
	// Show fuzzy finder with multi-select
	indices, err := fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			
			pkg := packages[i]
			var preview strings.Builder
			
			preview.WriteString(fmt.Sprintf("📦 Package: %s\n", pkg.Token))
			preview.WriteString(fmt.Sprintf("📝 Type: %s\n", pkg.Type))
			
			if pkg.FullName != "" && pkg.FullName != pkg.Token {
				preview.WriteString(fmt.Sprintf("📛 Full Name: %s\n", pkg.FullName))
			}
			
			if pkg.Version != "" {
				preview.WriteString(fmt.Sprintf("🏷️  Version: %s\n", pkg.Version))
			}
			
			if existing[pkg.Token] {
				preview.WriteString("\n✅ Already in Brewfile\n")
			} else {
				preview.WriteString("\n❌ Not in Brewfile\n")
			}
			
			if pkg.Description != "" {
				preview.WriteString(fmt.Sprintf("\n📄 Description:\n%s\n", wordWrap(pkg.Description, w-2)))
			}
			
			if pkg.Homepage != "" {
				preview.WriteString(fmt.Sprintf("\n🌐 Homepage:\n%s\n", pkg.Homepage))
			}
			
			return preview.String()
		}),
		fuzzyfinder.WithPromptString("🔍 Search packages (TAB to select, ENTER to confirm): "),
		fuzzyfinder.WithHeader("Use TAB to select/deselect, ENTER to confirm, ESC to cancel\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"),
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
		selected[i] = packages[idx]
	}
	
	return selected, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
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