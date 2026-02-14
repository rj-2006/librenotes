package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rj-2006/librenotes/internal/ui/styles"
)

// Homepage represents the welcome screen component
type Homepage struct {
	theme       styles.Theme
	recentFiles []string
	selectedIdx int
	width       int
	height      int
}

// MenuItem represents a menu option
type MenuItem struct {
	Label   string
	KeyHint string
	Action  string
}

// NewHomepage creates a new homepage component
func NewHomepage(theme styles.Theme) *Homepage {
	return &Homepage{
		theme:       theme,
		recentFiles: []string{},
		selectedIdx: 0,
	}
}

// SetRecentFiles updates the recent files list
func (h *Homepage) SetRecentFiles(files []string) {
	h.recentFiles = files
}

// SetSize sets the component dimensions
func (h *Homepage) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// SelectNext moves selection down
func (h *Homepage) SelectNext() {
	h.selectedIdx++
	if h.selectedIdx > 3 {
		h.selectedIdx = 0
	}
}

// SelectPrev moves selection up
func (h *Homepage) SelectPrev() {
	h.selectedIdx--
	if h.selectedIdx < 0 {
		h.selectedIdx = 3
	}
}

// GetSelectedAction returns the currently selected action
func (h *Homepage) GetSelectedAction() string {
	items := h.getMenuItems()
	if h.selectedIdx < len(items) {
		return items[h.selectedIdx].Action
	}
	return ""
}

// Render renders the homepage
func (h *Homepage) Render() string {
	// Two-column layout: Main content (left) | Recent files (right)
	mainContent := h.renderMainContent()
	recentFiles := h.renderRecentFilesSidebar()

	// Calculate heights and pad to match
	mainHeight := lipgloss.Height(mainContent)
	sidebarHeight := lipgloss.Height(recentFiles)
	maxHeight := mainHeight
	if sidebarHeight > maxHeight {
		maxHeight = sidebarHeight
	}

	// Pad main content to match sidebar height
	if mainHeight < maxHeight {
		padding := maxHeight - mainHeight
		mainContent = mainContent + strings.Repeat("\n", padding)
	}

	// Pad sidebar to match main content height
	if sidebarHeight < maxHeight {
		padding := maxHeight - sidebarHeight
		recentFiles = recentFiles + strings.Repeat("\n", padding)
	}

	// Join horizontally with separator
	separator := lipgloss.NewStyle().
		Foreground(h.theme.Muted).
		Render("│")

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainContent,
		"  ",
		separator,
		" ",
		recentFiles,
	)
}

func (h *Homepage) renderMainContent() string {
	// Logo section
	logo := h.renderLogo()

	// Spacing
	spacing := "\n\n"

	// Menu section
	menu := h.renderMenu()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		spacing,
		menu,
	)
}

func (h *Homepage) renderLogo() string {
	// ASCII art for wretched.dev - aligned properly
	asciiLogo := `██     ██ ██████  ███████ ████████  ██████  ██   ██ ███████ ██████  
██     ██ ██   ██ ██         ██    ██    ██ ██   ██ ██      ██   ██ 
██  █  ██ ██████  █████      ██    ██    ██ ███████ █████   ██   ██ 
██ ███ ██ ██   ██ ██         ██    ██    ██ ██   ██ ██      ██   ██      ██████  ███████ ██    ██ 
 ███ ███  ██   ██ ███████    ██     ██████  ██   ██ ███████ ██████       ██   ██ ██      ██    ██ 
                                                                          ██████  ███████  ██████  `

	// Style the ASCII art
	logoStyle := lipgloss.NewStyle().
		Foreground(h.theme.Primary).
		Bold(true)

	// LibreNotes text below
	textStyle := lipgloss.NewStyle().
		Foreground(h.theme.Secondary).
		Bold(true).
		MarginTop(1)

	logo := logoStyle.Render(asciiLogo)
	text := textStyle.Render("LibreNotes")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		text,
	)
}

func (h *Homepage) renderRecentFilesSidebar() string {
	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(h.theme.Primary).
		Bold(true).
		Underline(true)

	title := titleStyle.Render("Recent Files")

	// Files list
	var files []string

	if len(h.recentFiles) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(h.theme.Muted).
			Italic(true)
		files = append(files, emptyStyle.Render("No recent files"))
	} else {
		for i, file := range h.recentFiles {
			if i >= 10 { // Show max 10 recent files
				break
			}
			// Truncate if too long
			display := file
			if len(display) > 25 {
				display = display[:22] + "..."
			}

			// Numbered list
			numStyle := lipgloss.NewStyle().
				Foreground(h.theme.Muted)
			nameStyle := lipgloss.NewStyle().
				Foreground(h.theme.Foreground)

			line := numStyle.Render(fmt.Sprintf("%d. ", i+1)) + nameStyle.Render(display)
			files = append(files, line)
		}
	}

	// Join all files
	filesContent := lipgloss.JoinVertical(
		lipgloss.Left,
		files...,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		filesContent,
	)
}

func (h *Homepage) renderMenu() string {
	items := h.getMenuItems()

	var renderedItems []string

	for i, item := range items {
		isSelected := i == h.selectedIdx
		renderedItems = append(renderedItems, h.renderMenuItem(item, isSelected))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		renderedItems...,
	)
}

func (h *Homepage) renderMenuItem(item MenuItem, isSelected bool) string {
	if isSelected {
		// Selected: bright color with arrow
		style := lipgloss.NewStyle().
			Foreground(h.theme.Primary).
			Bold(true)

		return style.Render(fmt.Sprintf(" > %s  %s", item.Label, lipgloss.NewStyle().Foreground(h.theme.Muted).Render(item.KeyHint)))
	}

	// Unselected: dim color with space
	style := lipgloss.NewStyle().
		Foreground(h.theme.Muted)

	return style.Render(fmt.Sprintf("   %s  %s", item.Label, lipgloss.NewStyle().Foreground(h.theme.Muted).Faint(true).Render(item.KeyHint)))
}

func (h *Homepage) getMenuItems() []MenuItem {
	return []MenuItem{
		{Label: "New File", KeyHint: "Ctrl+N", Action: "new"},
		{Label: "Open File", KeyHint: "Enter", Action: "open"},
		{Label: "View List", KeyHint: "Ctrl+L", Action: "list"},
		{Label: "Change Theme", KeyHint: "Ctrl+T", Action: "settings"},
	}
}
