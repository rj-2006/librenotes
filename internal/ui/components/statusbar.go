package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/rj-2006/librenotes/internal/ui/styles"
)

// StatusBar is a reusable status bar component
type StatusBar struct {
	theme       styles.Theme
	currentFile string
	wordCount   int
	charCount   int
	showTime    bool
}

// NewStatusBar creates a new status bar
func NewStatusBar(theme styles.Theme) *StatusBar {
	return &StatusBar{
		theme:    theme,
		showTime: true,
	}
}

// SetTheme updates the theme
func (sb *StatusBar) SetTheme(theme styles.Theme) {
	sb.theme = theme
}

// SetFileInfo updates the file information
func (sb *StatusBar) SetFileInfo(filename string, content string) {
	sb.currentFile = filename
	sb.wordCount = countWords(content)
	sb.charCount = len([]rune(content))
}

// SetWordCount updates just the word count
func (sb *StatusBar) SetWordCount(content string) {
	sb.wordCount = countWords(content)
	sb.charCount = len([]rune(content))
}

// Render renders the status bar
func (sb *StatusBar) Render(width int) string {
	style := lipgloss.NewStyle().
		Background(sb.theme.StatusBar).
		Foreground(sb.theme.StatusBarText).
		Padding(0, 1).
		Width(width)

	left := sb.renderLeft()
	right := sb.renderRight()

	// Calculate padding needed
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	padding := width - leftWidth - rightWidth - 2

	if padding < 0 {
		padding = 0
	}

	content := left + fmt.Sprintf("%*s", padding, "") + right
	return style.Render(content)
}

func (sb *StatusBar) renderLeft() string {
	if sb.currentFile == "" {
		return "  LibreNotes"
	}
	return fmt.Sprintf("  %s", sb.currentFile)
}

func (sb *StatusBar) renderRight() string {
	var info string
	if sb.wordCount > 0 {
		info = fmt.Sprintf("%d words | %d chars", sb.wordCount, sb.charCount)
	}

	if sb.showTime {
		now := time.Now().Format("15:04")
		if info != "" {
			info = fmt.Sprintf("%s | %s", info, now)
		} else {
			info = now
		}
	}

	return info + "  "
}

func countWords(content string) int {
	wordCount := 0
	inWord := false
	for _, ch := range content {
		if ch == ' ' || ch == '\n' || ch == '\t' {
			if inWord {
				wordCount++
				inWord = false
			}
		} else {
			inWord = true
		}
	}
	if inWord {
		wordCount++
	}
	return wordCount
}

// HelpBar shows keyboard shortcuts
type HelpBar struct {
	theme styles.Theme
	items []HelpItem
}

// HelpItem represents a single help entry
type HelpItem struct {
	Key  string
	Desc string
}

// NewHelpBar creates a new help bar
func NewHelpBar(theme styles.Theme) *HelpBar {
	return &HelpBar{
		theme: theme,
		items: []HelpItem{
			{Key: "Ctrl+N", Desc: "new"},
			{Key: "Ctrl+L", Desc: "list"},
			{Key: "Ctrl+S", Desc: "save"},
			{Key: "Esc", Desc: "quit"},
		},
	}
}

// SetItems updates the help items
func (hb *HelpBar) SetItems(items []HelpItem) {
	hb.items = items
}

// Render renders the help bar
func (hb *HelpBar) Render() string {
	var items []string

	for _, item := range hb.items {
		keyStyle := lipgloss.NewStyle().
			Foreground(hb.theme.Primary).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(hb.theme.HelpText)

		rendered := keyStyle.Render(item.Key) + " " + descStyle.Render(item.Desc)
		items = append(items, rendered)
	}

	separator := lipgloss.NewStyle().
		Foreground(hb.theme.Muted).
		Render(" • ")

	return lipgloss.JoinHorizontal(lipgloss.Center, items...) + separator
}
