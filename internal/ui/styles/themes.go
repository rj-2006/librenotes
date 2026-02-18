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
	// Cyberpunk theme - improved neon colors on dark background
	CyberpunkTheme = Theme{
		Name:          "cyberpunk",
		Primary:       lipgloss.Color("#00d4ff"), // Electric cyan - vibrant but not blinding
		Secondary:     lipgloss.Color("#ff2a6d"), // Hot pink - better contrast
		Background:    lipgloss.Color("#0d0221"), // Deep purple-black
		Foreground:    lipgloss.Color("#f0f0ff"), // Soft white
		Accent:        lipgloss.Color("#ffd60a"), // Golden yellow
		Success:       lipgloss.Color("#39ff14"), // Neon green
		Error:         lipgloss.Color("#ff073a"), // Electric red
		Warning:       lipgloss.Color("#ff9f1c"), // Amber
		Muted:         lipgloss.Color("#5c5c8a"), // Desaturated purple
		Border:        lipgloss.Color("#00d4ff"), // Cyan border
		Cursor:        lipgloss.Color("#ff2a6d"), // Hot pink cursor
		Selection:     lipgloss.Color("#ff2a6d"), // Hot pink selection
		StatusBar:     lipgloss.Color("#1a1a3e"), // Dark purple
		StatusBarText: lipgloss.Color("#00d4ff"), // Electric cyan
		TitleBg:       lipgloss.Color("#ff2a6d"), // Hot pink background
		TitleFg:       lipgloss.Color("#0d0221"), // Deep purple-black text
		HelpText:      lipgloss.Color("#7a7ab8"), // Muted purple-blue
	}

	// Catppuccin Mocha theme - improved warm dark theme
	CatppuccinMochaTheme = Theme{
		Name:          "catppuccin-mocha",
		Primary:       lipgloss.Color("#8aadf4"), // Soft blue - more readable
		Secondary:     lipgloss.Color("#b4befe"), // Periwinkle - distinct from primary
		Background:    lipgloss.Color("#181825"), // Darker background for contrast
		Foreground:    lipgloss.Color("#cdd6f4"), // Text
		Accent:        lipgloss.Color("#f5e0dc"), // Rosewater - softer accent
		Success:       lipgloss.Color("#a6e3a1"), // Green
		Error:         lipgloss.Color("#f38ba8"), // Red
		Warning:       lipgloss.Color("#fab387"), // Peach
		Muted:         lipgloss.Color("#6c7086"), // Overlay0
		Border:        lipgloss.Color("#8aadf4"), // Soft blue border
		Cursor:        lipgloss.Color("#b4befe"), // Periwinkle cursor
		Selection:     lipgloss.Color("#45475a"), // Surface0 - lighter for visibility
		StatusBar:     lipgloss.Color("#11111b"), // Crust - darker
		StatusBarText: lipgloss.Color("#cdd6f4"), // Text
		TitleBg:       lipgloss.Color("#8aadf4"), // Soft blue background
		TitleFg:       lipgloss.Color("#181825"), // Dark text
		HelpText:      lipgloss.Color("#9399b2"), // Overlay2
	}

	// Catppuccin Latte theme - light theme
	CatppuccinLatteTheme = Theme{
		Name:          "catppuccin-latte",
		Primary:       lipgloss.Color("#1e66f5"), // Blue
		Secondary:     lipgloss.Color("#8839ef"), // Mauve
		Background:    lipgloss.Color("#eff1f5"), // Light background
		Foreground:    lipgloss.Color("#4c4f69"), // Text
		Accent:        lipgloss.Color("#df8e1f"), // Yellow
		Success:       lipgloss.Color("#40a02b"), // Green
		Error:         lipgloss.Color("#d20f39"), // Red
		Warning:       lipgloss.Color("#fe640b"), // Peach
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
