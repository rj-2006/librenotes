package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory!", err)
	}

	vaultDir = fmt.Sprintf("%s/.librenotes", homeDir)
}

type model struct {
	newFileInput  textinput.Model
	IsTextVisible bool
	currentFile   *os.File
	textarea      textarea.Model
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlX:
			m.IsTextVisible = true
			return m, cmd
		case tea.KeyEnter:
			// todo: create a file with name which is given
			filename := m.newFileInput.Value()
			if filename != "" {
				filepath := fmt.Sprintf("%s/%s.md", vaultDir, filename)

				if _, err := os.Stat(filepath); err == nil {
					return m, nil
				}

				f, err := os.Create(filepath)
				if err != nil {
					log.Fatalf("%v", err)
				}
				m.currentFile = f
				m.IsTextVisible = false
				m.newFileInput.SetValue("")
			}
			return m, nil
		}
	}
	if m.IsTextVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}
	if m.currentFile != nil {
		m.textarea, cmd = m.textarea.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#16")).
		Background(lipgloss.Color("#205")).
		PaddingLeft(2).
		PaddingRight(2)

	welcome := style.Render("Welcome to Librenotes!")
	help := "Ctrl+X: new file - Ctrl+L: list - Esc: Quit - Ctrl+S: save - Ctrl+Q: quit"
	view := ""
	if m.IsTextVisible {
		view = m.newFileInput.View()
	}
	if m.currentFile != nil {
		view = m.textarea.View()
	}
	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initialzeModel() model {

	// Creating a folder first of all
	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	ti := textinput.New()
	ti.Placeholder = "What do you want to name the file?"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	//Object for Text Area
	ta := textarea.New()
	ta.Placeholder = "lorem ipsum"
	ta.Focus()
	ta.ShowLineNumbers = false

	return model{
		newFileInput:  ti,
		IsTextVisible: false,
		textarea:      ta,
	}
}

func main() {
	p := tea.NewProgram(initialzeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("ur code bad : %v", err)
		os.Exit(1)
	}
}
