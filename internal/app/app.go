package app

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/rj-2006/librenotes/internal/config"
	"github.com/rj-2006/librenotes/internal/storage"
	"github.com/rj-2006/librenotes/internal/ui/components"
	"github.com/rj-2006/librenotes/internal/ui/styles"
)

func NewApp(cfg *config.Config) (*App, error) {
	// Initialize storage
	store := storage.NewStorage(cfg.GetVaultDir())
	if err := store.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Load theme
	theme := styles.GetTheme(cfg.GetTheme())

	// Initialize UI components
	ti := newFileInput(theme)
	ta := newTextArea(theme)
	lst := newNoteList(theme, store)
	sb := components.NewStatusBar(theme)
	hb := components.NewHelpBar(theme)
	hp := components.NewHomepage(theme)

	return &App{
		state:        StateWelcome,
		config:       cfg,
		storage:      store,
		theme:        theme,
		newFileInput: ti,
		textarea:     ta,
		list:         lst,
		statusBar:    sb,
		helpBar:      hb,
		homepage:     hp,
	}, nil
}

// Init initializes the app
func (a *App) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the app state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		// Update list size
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		a.list.SetSize(msg.Width-h, msg.Height-v-5)

		// Update homepage size
		a.homepage.SetSize(msg.Width-4, msg.Height-10)
		a.homepage.SetRecentFiles(a.recentFiles)

		// Update preview viewport if active
		if a.isPreview {
			vpWidth := msg.Width - 4
			if vpWidth < 40 {
				vpWidth = 78
			}
			vpHeight := msg.Height - 8
			if vpHeight < 10 {
				vpHeight = 20
			}
			a.previewViewport.Width = vpWidth
			a.previewViewport.Height = vpHeight
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlQ:
			return a, tea.Quit

		case tea.KeyEsc:
			return a.handleBack()

		case tea.KeyCtrlN:
			a.state = StateNewFile
			a.newFileInput.Focus()
			return a, cmd

		case tea.KeyCtrlL:
			a.state = StateList
			// Refresh the list
			a.refreshList()
			return a, cmd

		case tea.KeyUp:
			if a.state == StateWelcome {
				a.homepage.SelectPrev()
			}

		case tea.KeyDown:
			if a.state == StateWelcome {
				a.homepage.SelectNext()
			}

		case tea.KeyEnter:
			// Don't intercept Enter when editing - let textarea handle it
			if a.state == StateEditing {
				break
			}
			return a.handleEnter()

		case tea.KeyCtrlT:
			if a.state == StateWelcome {
				return a.handleHomepageAction("settings")
			}

		case tea.KeyCtrlS:
			return a.handleSave()

		case tea.KeyCtrlP:
			if a.state == StateEditing {
				return a, a.togglePreview()
			}
		}
	}

	// Update active component based on state
	switch a.state {
	case StateNewFile:
		a.newFileInput, cmd = a.newFileInput.Update(msg)
	case StateEditing:
		if a.isPreview {
			// Handle viewport scrolling
			var vpCmd tea.Cmd
			a.previewViewport, vpCmd = a.previewViewport.Update(msg)
			cmd = vpCmd
		} else {
			a.textarea, cmd = a.textarea.Update(msg)
			// Update word count
			a.statusBar.SetWordCount(a.textarea.Value())
		}
	case StateList:
		a.list, cmd = a.list.Update(msg)
	}

	return a, cmd
}

// View renders the current view
func (a *App) View() string {
	// Main content based on state
	var content string
	switch a.state {
	case StateWelcome:
		content = a.homepage.Render()
	case StateList:
		content = a.list.View()
	case StateNewFile:
		content = a.newFileInput.View()
	case StateEditing:
		if a.isPreview {
			content = a.previewViewport.View()
		} else {
			content = a.textarea.View()
		}
	default:
		content = ""
	}

	// Status bar
	status := a.statusBar.Render(a.width)

	// Help bar
	help := a.helpBar.Render()

	return fmt.Sprintf("\n%s\n\n%s\n%s", content, status, help)
}

// handleBack navigates back to the previous screen
func (a *App) handleBack() (tea.Model, tea.Cmd) {
	switch a.state {
	case StateEditing:
		// If in preview mode, exit preview first before closing file
		if a.isPreview {
			a.isPreview = false
			a.previewViewport = viewport.Model{}
			return a, nil
		}
		// Close file without saving, go back to list
		if a.currentFile != nil {
			a.storage.CloseNote(a.currentFile)
			a.currentFile = nil
		}
		a.currentFileName = ""
		a.textarea.SetValue("")
		a.statusBar.SetFileInfo("", "")
		a.state = StateList
	case StateList, StateNewFile:
		// Go back to welcome screen
		a.state = StateWelcome
		a.homepage.SetRecentFiles(a.recentFiles)
	case StateWelcome:
		// Already at welcome, quit
		return a, tea.Quit
	}
	return a, nil
}

// handleEnter processes the Enter key press
func (a *App) handleEnter() (tea.Model, tea.Cmd) {
	switch a.state {
	case StateWelcome:
		action := a.homepage.GetSelectedAction()
		return a.handleHomepageAction(action)
	case StateList:
		return a.openSelectedNote()
	case StateNewFile:
		return a.createNewNote()
	}
	return a, nil
}

// handleHomepageAction handles menu actions from the homepage
func (a *App) handleHomepageAction(action string) (tea.Model, tea.Cmd) {
	switch action {
	case "new":
		a.state = StateNewFile
		a.newFileInput.Focus()
		return a, nil
	case "open":
		// Open the first recent file or go to list
		if len(a.recentFiles) > 0 {
			return a.openRecentFile(a.recentFiles[0])
		}
		a.state = StateList
		a.refreshList()
		return a, nil
	case "list":
		a.state = StateList
		a.refreshList()
		return a, nil
	case "settings":
		// Cycle to next theme
		themes := styles.GetAvailableThemes()
		currentIdx := 0
		for i, t := range themes {
			if t == a.config.GetTheme() {
				currentIdx = i
				break
			}
		}
		nextIdx := (currentIdx + 1) % len(themes)
		a.config.SetTheme(themes[nextIdx])
		a.theme = styles.GetTheme(themes[nextIdx])
		// Update all components with new theme
		a.homepage = components.NewHomepage(a.theme)
		a.homepage.SetRecentFiles(a.recentFiles) // Preserve recent files
		a.homepage.SetSize(a.width-4, a.height-10)
		a.statusBar = components.NewStatusBar(a.theme)
		a.helpBar = components.NewHelpBar(a.theme)
		a.list = newNoteList(a.theme, a.storage)
		// Resize the list to match current window
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		a.list.SetSize(a.width-h, a.height-v-5)
		a.newFileInput = newFileInput(a.theme)
		a.newFileInput.Focus() // Refocus the input
		a.textarea = newTextArea(a.theme)
		return a, nil
	}
	return a, nil
}

// openRecentFile opens a file from the recent files list
func (a *App) openRecentFile(filename string) (tea.Model, tea.Cmd) {
	content, err := a.storage.ReadNote(filename)
	if err != nil {
		return a, nil
	}

	file, err := a.storage.OpenNote(filename)
	if err != nil {
		return a, nil
	}

	a.currentFile = file
	a.currentFileName = filename
	a.textarea.SetValue(string(content))
	a.statusBar.SetFileInfo(filename, string(content))
	a.state = StateEditing

	// Add to recent files (move to top)
	a.addToRecentFiles(filename)

	return a, nil
}

// openSelectedNote opens the selected note from the list
func (a *App) openSelectedNote() (tea.Model, tea.Cmd) {
	item, ok := a.list.SelectedItem().(NoteItem)
	if !ok {
		return a, nil
	}

	content, err := a.storage.ReadNote(item.title)
	if err != nil {
		return a, nil
	}

	file, err := a.storage.OpenNote(item.title)
	if err != nil {
		return a, nil
	}

	a.currentFile = file
	a.currentFileName = item.title
	a.textarea.SetValue(string(content))
	a.statusBar.SetFileInfo(item.title, string(content))
	a.state = StateEditing

	// Add to recent files
	a.addToRecentFiles(item.title)

	return a, nil
}

// createNewNote creates a new note
func (a *App) createNewNote() (tea.Model, tea.Cmd) {
	filename := a.newFileInput.Value()
	if filename == "" {
		return a, nil
	}

	file, err := a.storage.CreateNote(filename)
	if err != nil {
		// File already exists or other error
		return a, nil
	}

	a.currentFile = file
	a.currentFileName = filepath.Base(file.Name())
	a.newFileInput.SetValue("")
	a.statusBar.SetFileInfo(a.currentFileName, "")
	a.state = StateEditing

	// Add to recent files
	a.addToRecentFiles(a.currentFileName)

	return a, nil
}

// handleSave saves the current note
func (a *App) handleSave() (tea.Model, tea.Cmd) {
	if a.currentFile == nil {
		return a, nil
	}

	content := a.textarea.Value()
	if err := a.storage.SaveNote(a.currentFile, content); err != nil {
		return a, nil
	}

	if err := a.storage.CloseNote(a.currentFile); err != nil {
		return a, nil
	}

	a.currentFile = nil
	a.currentFileName = ""
	a.textarea.SetValue("")
	a.state = StateWelcome
	a.statusBar.SetFileInfo("", "")
	a.homepage.SetRecentFiles(a.recentFiles)

	return a, nil
}

// refreshList refreshes the note list
func (a *App) refreshList() {
	notes, err := a.storage.ListNotes()
	if err != nil {
		return
	}

	items := make([]list.Item, len(notes))
	for i, note := range notes {
		items[i] = NoteItem{
			title:       note.Title,
			description: note.Description,
		}
	}

	a.list.SetItems(items)
}

// addToRecentFiles adds a file to recent files list
func (a *App) addToRecentFiles(filename string) {
	// Remove if already exists
	for i, f := range a.recentFiles {
		if f == filename {
			a.recentFiles = append(a.recentFiles[:i], a.recentFiles[i+1:]...)
			break
		}
	}

	// Add to front
	a.recentFiles = append([]string{filename}, a.recentFiles...)

	// Keep only last 10
	if len(a.recentFiles) > 10 {
		a.recentFiles = a.recentFiles[:10]
	}
}

// Helper functions for creating UI components
func newFileInput(theme styles.Theme) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "What do you want to name the file?"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40

	cursorStyle := lipgloss.NewStyle().Foreground(theme.Cursor)
	textStyle := lipgloss.NewStyle().Foreground(theme.Foreground)
	placeholderStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = textStyle
	ti.PlaceholderStyle = placeholderStyle

	return ti
}

func newTextArea(theme styles.Theme) textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Start typing your note..."
	ta.Focus()
	ta.ShowLineNumbers = false

	// Style the textarea with theme colors
	ta.Prompt = ""
	cursorStyle := lipgloss.NewStyle().Foreground(theme.Cursor)
	ta.Cursor.Style = cursorStyle

	return ta
}

func newNoteList(theme styles.Theme, store *storage.Storage) list.Model {
	notes, _ := store.ListNotes()

	items := make([]list.Item, len(notes))
	for i, note := range notes {
		items[i] = NoteItem{
			title:       note.Title,
			description: note.Description,
		}
	}

	// Create delegate with theme colors
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(theme.Foreground)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Foreground(theme.Muted)
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(theme.Secondary)

	lst := list.New(items, delegate, 0, 0)
	lst.Title = "All Notes"
	lst.Styles.Title = lipgloss.NewStyle().
		Foreground(theme.TitleFg).
		Background(theme.TitleBg).
		Padding(0, 1)

	return lst
}

// rendering the markdown file
func (a *App) renderMarkdown(content string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	if err != nil {
		return "", err
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return "", err
	}

	return rendered, nil
}

func (a *App) togglePreview() tea.Cmd {
	if a.isPreview {
		a.isPreview = false
		a.previewViewport = viewport.Model{}
	} else {
		// Calculate viewport dimensions
		vpWidth := a.width - 4
		if vpWidth < 40 {
			vpWidth = 78
		}
		vpHeight := a.height - 8
		if vpHeight < 10 {
			vpHeight = 20
		}

		// Create viewport with themed border
		vp := viewport.New(vpWidth, vpHeight)
		vp.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(a.theme.Border).
			PaddingRight(2)

		// Calculate glamour render width accounting for viewport frame and gutter
		glamourGutter := 2
		glamourRenderWidth := vpWidth - vp.Style.GetHorizontalFrameSize() - glamourGutter

		// Create glamour renderer with word wrap
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(glamourRenderWidth),
		)
		if err != nil {
			return nil
		}

		// Render markdown content
		rendered, err := renderer.Render(a.textarea.Value())
		if err != nil {
			return nil
		}

		vp.SetContent(rendered)
		a.previewViewport = vp
		a.isPreview = true
	}
	return nil
}
