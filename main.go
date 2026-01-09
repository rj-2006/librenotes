package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	vaultDir    string
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	docStyle    = lipgloss.NewStyle().Margin(1, 2)
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
	list          list.Model
	showlist      bool
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlN:
			m.IsTextVisible = true
			return m, cmd
		case tea.KeyEnter:
			if m.currentFile != nil {
				break
			}
			if m.showlist {
				item, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, item.title)

					content, err := os.ReadFile(filepath)
					if err != nil {
						log.Fatal("Error opening file")
						return m, nil
					}

					m.textarea.SetValue(string(content))
					f, err := os.OpenFile(filepath, os.O_RDWR, 0644)
					if err != nil {
						log.Fatal("Can't open file!")
					}
					m.currentFile = f
					m.showlist = false
				}
				return m, nil
			}
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
		case tea.KeyCtrlS:

			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Printf("Can't save the file T-T")
				return m, nil
			}
			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Printf("Can't save the file T-T")
				return m, nil
			}
			if _, err := m.currentFile.WriteString(m.textarea.Value()); err != nil {
				fmt.Printf("Can't save the file T-T")
				return m, nil
			}
			if err := m.currentFile.Close(); err != nil {
				fmt.Println("Can't close file T-T")
			}
			m.currentFile = nil
			m.textarea.SetValue("")
			return m, nil
		case tea.KeyCtrlL:
			//todo show list
			m.showlist = true
			return m, nil
		}
	}
	if m.IsTextVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}
	if m.currentFile != nil {
		m.textarea, cmd = m.textarea.Update(msg)
	}
	if m.showlist {
		m.list, cmd = m.list.Update(msg)
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
	if m.showlist {
		view = m.list.View()
	} else if m.IsTextVisible {
		view = m.newFileInput.View()
	} else if m.currentFile != nil {
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

	//List
	notelist := listFiles()

	finalList := list.New(notelist, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "All Notes"
	finalList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("254")).
		Padding(0, 1)

	return model{
		newFileInput:  ti,
		IsTextVisible: false,
		textarea:      ta,
		list:          finalList,
	}
}

func main() {
	p := tea.NewProgram(initialzeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("ur code bad : %v", err)
		os.Exit(1)
	}
}

func listFiles() []list.Item {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatal("Can't read from directory :(")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			modTime := info.ModTime().Format("2006-01-02 15:04")

			items = append(items, item{
				title: entry.Name(),
				desc:  fmt.Sprintf("Modified: %s", modTime),
			})
		}
	}

	return items
}
