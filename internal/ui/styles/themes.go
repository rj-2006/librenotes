package styles

import "github.com/charmbracelet/lipgloss"

// Theme defines a color theme for the application
type Theme struct {
	Name          string
	Primary       lipgloss.Color
	Secondary     lipgloss.Color
	Background    lipgloss.Color
	Foreground    lipgloss.Color
	Accent        lipgloss.Color
	Success       lipgloss.Color
	Error         lipgloss.Color
	Warning       lipgloss.Color
	Muted         lipgloss.Color
	Border        lipgloss.Color
	Cursor        lipgloss.Color
	Selection     lipgloss.Color
	StatusBar     lipgloss.Color
	StatusBarText lipgloss.Color
	TitleBg       lipgloss.Color
	TitleFg       lipgloss.Color
	HelpText      lipgloss.Color
}

// Available themes
var (
	// Cyberpunk theme - neon colors on dark background
	CyberpunkTheme = Theme{
		Name:          "cyberpunk",
		Primary:       lipgloss.Color("#00f0ff"), // Cyan
		Secondary:     lipgloss.Color("#ff00ff"), // Magenta
		Background:    lipgloss.Color("#0a0a0f"), // Very dark
		Foreground:    lipgloss.Color("#e0e0ff"), // Light bluish white
		Accent:        lipgloss.Color("#ffff00"), // Yellow
		Success:       lipgloss.Color("#00ff88"), // Bright green
		Error:         lipgloss.Color("#ff3366"), // Pink-red
		Warning:       lipgloss.Color("#ffaa00"), // Orange
		Muted:         lipgloss.Color("#666699"), // Muted purple
		Border:        lipgloss.Color("#00f0ff"), // Cyan border
		Cursor:        lipgloss.Color("#ff00ff"), // Magenta cursor
		Selection:     lipgloss.Color("#ff00ff"), // Magenta selection
		StatusBar:     lipgloss.Color("#1a1a2e"), // Dark blue-purple
		StatusBarText: lipgloss.Color("#00f0ff"), // Cyan text
		TitleBg:       lipgloss.Color("#ff00ff"), // Magenta background
		TitleFg:       lipgloss.Color("#000000"), // Black text
		HelpText:      lipgloss.Color("#8888aa"), // Muted help text
	}

	// Catppuccin Mocha theme - warm dark theme
	CatppuccinMochaTheme = Theme{
		Name:          "catppuccin-mocha",
		Primary:       lipgloss.Color("#89b4fa"), // Blue
		Secondary:     lipgloss.Color("#cba6f7"), // Lavender
		Background:    lipgloss.Color("#1e1e2e"), // Dark background
		Foreground:    lipgloss.Color("#cdd6f4"), // Text
		Accent:        lipgloss.Color("#f9e2af"), // Yellow
		Success:       lipgloss.Color("#a6e3a1"), // Green
		Error:         lipgloss.Color("#f38ba8"), // Red
		Warning:       lipgloss.Color("#fab387"), // Peach
		Muted:         lipgloss.Color("#6c7086"), // Overlay0
		Border:        lipgloss.Color("#89b4fa"), // Blue border
		Cursor:        lipgloss.Color("#cba6f7"), // Lavender cursor
		Selection:     lipgloss.Color("#585b70"), // Surface1
		StatusBar:     lipgloss.Color("#181825"), // Mantle
		StatusBarText: lipgloss.Color("#cdd6f4"), // Text
		TitleBg:       lipgloss.Color("#89b4fa"), // Blue background
		TitleFg:       lipgloss.Color("#1e1e2e"), // Dark text
		HelpText:      lipgloss.Color("#9399b2"), // Overlay2
	}

	// Catppuccin Latte theme - light theme
	CatppuccinLatteTheme = Theme{
		Name:          "catppuccin-latte",
		Primary:       lipgloss.Color("#1e66f5"), // Blue
		Secondary:     lipgloss.Color("#8839ef"), // Mauve
		Background:    lipgloss.Color("#eff1f5"), // Light background
		Foreground:    lipgloss.Color("#4c4f69"), // Text
		Accent:        lipgloss.Color("(#df8e1"), // Yellow
		Success:       lipgloss.Color("#40a02b"), // Green
		Error:         lipgloss.Color("#d20f39"), // Red
		Warning:       lipgloss.Color("(#fe640"), // Peach
		Muted:         lipgloss.Color("#9ca0b0"), // Overlay0
		Border:        lipgloss.Color("#1e66f5"), // Blue border
		Cursor:        lipgloss.Color("#8839ef"), // Mauve cursor
		Selection:     lipgloss.Color("#ccd0da"), // Surface1
		StatusBar:     lipgloss.Color("#e6e9ef"), // Mantle
		StatusBarText: lipgloss.Color("#4c4f69"), // Text
		TitleBg:       lipgloss.Color("#1e66f5"), // Blue background
		TitleFg:       lipgloss.Color("#eff1f5"), // Light text
		HelpText:      lipgloss.Color("#7c7f93"), // Overlay2
	}

	// Default theme
	DefaultTheme = CyberpunkTheme
)

// Themes map for easy lookup
var Themes = map[string]Theme{
	"cyberpunk":        CyberpunkTheme,
	"catppuccin-mocha": CatppuccinMochaTheme,
	"catppuccin-latte": CatppuccinLatteTheme,
}

// GetTheme returns a theme by name, defaults to cyberpunk if not found
func GetTheme(name string) Theme {
	if theme, ok := Themes[name]; ok {
		return theme
	}
	return DefaultTheme
}

// GetAvailableThemes returns a list of available theme names
func GetAvailableThemes() []string {
	names := make([]string, 0, len(Themes))
	for name := range Themes {
		names = append(names, name)
	}
	return names
}
