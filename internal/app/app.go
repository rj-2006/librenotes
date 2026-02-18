package app

import (
	"fmt"
	"path/filepath"
	"time"

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

	// Initialize search input
	searchInput := newSearchInput(theme)

	// Load all notes for search (recursively from all folders)
	allNotes, _ := store.ListAllNotesRecursively()
	items := make([]list.Item, len(allNotes))
	for i, note := range allNotes {
		displayTitle := note.Title
		if len(displayTitle) > 30 {
			displayTitle = displayTitle[:27] + "..."
		}
		items[i] = NoteItem{
			title:       "📝 " + displayTitle,
			description: note.Description,
			isFolder:    false,
			folderPath:  note.Title, // Store full path here
		}
	}

	app := &App{
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
		searchInput:  searchInput,
		allNotes:     items,
		recentFiles:  cfg.GetRecentFiles(),
	}

	// Set recent files on homepage
	app.homepage.SetRecentFiles(app.recentFiles)

	return app, nil
}

// Init initializes the app
func (a *App) Init() tea.Cmd {
	// Start autosave timer if enabled
	if a.config.GetAutoSave() {
		return tea.Batch(
			textinput.Blink,
			a.startAutoSave(),
		)
	}
	return textinput.Blink
}

// startAutoSave starts the autosave timer
func (a *App) startAutoSave() tea.Cmd {
	interval := time.Duration(a.config.GetAutoSaveInterval()) * time.Second
	return tea.Every(interval, func(t time.Time) tea.Msg {
		return autoSaveMsg{}
	})
}

// autoSaveMsg is sent when it's time to autosave
type autoSaveMsg struct{}

// autoSave saves the current note if it has unsaved changes
func (a *App) autoSave() {
	if a.state != StateEditing || a.currentFile == nil {
		return
	}

	currentContent := a.textarea.Value()
	if currentContent == a.lastSavedContent {
		return // No changes to save
	}

	// Save the content
	if err := a.storage.SaveNote(a.currentFile, currentContent); err != nil {
		return // Silent fail on autosave
	}

	a.lastSavedContent = currentContent
}

// Update handles messages and updates the app state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case autoSaveMsg:
		a.autoSave()
		return a, a.startAutoSave() // Restart timer

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
			a.previewViewport.Width = msg.Width - 4
			a.previewViewport.Height = msg.Height - 8
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
			// In search mode, select the highlighted item
			if a.state == StateSearch {
				return a.handleEnter()
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

		case tea.KeyDelete:
			if a.state == StateList {
				item, ok := a.list.SelectedItem().(NoteItem)
				if ok && !item.isFolder {
					// Use folderPath which contains the full relative path
					a.fileToDelete = item.folderPath
					a.state = StateConfirmDelete
				}
			}

		case tea.KeyRunes:
			if a.state == StateList && msg.String() == "/" {
				a.state = StateSearch
				a.searchInput.Focus()
				return a, textinput.Blink
			}
			if a.state == StateConfirmDelete {
				switch msg.String() {
				case "y", "Y":
					if err := a.storage.DeleteNote(a.fileToDelete); err == nil {
						a.removeFromRecentFiles(a.fileToDelete)
						a.refreshList()
					}
					a.fileToDelete = ""
					a.state = StateList
				case "n", "N":
					a.fileToDelete = ""
					a.state = StateList
				}
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
	case StateSearch:
		var searchCmd tea.Cmd
		a.searchInput, searchCmd = a.searchInput.Update(msg)
		// Also update the list so user can navigate with arrow keys
		var listCmd tea.Cmd
		a.list, listCmd = a.list.Update(msg)
		cmd = tea.Batch(searchCmd, listCmd)
		// Perform fuzzy search
		a.performSearch(a.searchInput.Value())
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
	case StateConfirmDelete:
		confirmStyle := lipgloss.NewStyle().
			Foreground(a.theme.Error).
			Bold(true).
			Padding(2)
		content = confirmStyle.Render(fmt.Sprintf("Delete '%s'?\n\n[y]es / [n]o", a.fileToDelete))
	case StateSearch:
		// Show search input and list
		content = a.searchInput.View() + "\n\n" + a.list.View()
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
		// If in preview mode, exit preview first
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
		a.currentFolder = "" // Reset folder navigation
		a.homepage.SetRecentFiles(a.recentFiles)
	case StateSearch:
		// Exit search mode and return to list
		a.state = StateList
		a.searchInput.SetValue("")
		a.refreshList()
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
	case StateSearch:
		// Open the selected note from search results
		item, ok := a.list.SelectedItem().(NoteItem)
		if ok && !item.isFolder {
			// Clear search and open note
			a.searchInput.SetValue("")
			a.state = StateList
			return a.openNoteByPath(item.folderPath)
		}
		return a, nil
	case StateNewFile:
		return a.createNewNote()
	}
	return a, nil
}

// openNoteByPath opens a note given its full path
func (a *App) openNoteByPath(fullPath string) (tea.Model, tea.Cmd) {
	content, err := a.storage.ReadNote(fullPath)
	if err != nil {
		return a, nil
	}

	file, err := a.storage.OpenNote(fullPath)
	if err != nil {
		return a, nil
	}

	a.currentFile = file
	a.currentFileName = fullPath
	a.textarea.SetValue(string(content))
	a.lastSavedContent = string(content) // Track for autosave
	a.statusBar.SetFileInfo(fullPath, string(content))
	a.state = StateEditing

	// Add to recent files
	a.addToRecentFiles(fullPath)

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
		a.homepage.SetTheme(a.theme)
		a.statusBar = components.NewStatusBar(a.theme)
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
	a.lastSavedContent = string(content) // Track for autosave
	a.statusBar.SetFileInfo(filename, string(content))
	a.state = StateEditing

	// Add to recent files (move to top)
	a.addToRecentFiles(filename)

	return a, nil
}

// openSelectedNote opens the selected note from the list or navigates into a folder
func (a *App) openSelectedNote() (tea.Model, tea.Cmd) {
	item, ok := a.list.SelectedItem().(NoteItem)
	if !ok {
		return a, nil
	}

	// Handle folder navigation
	if item.isFolder {
		if item.folderPath == ".." {
			// Go up one level
			if a.currentFolder != "" {
				a.currentFolder = filepath.Dir(a.currentFolder)
				if a.currentFolder == "." {
					a.currentFolder = ""
				}
			}
		} else {
			// Enter folder
			a.currentFolder = item.folderPath
		}
		a.refreshList()
		return a, nil
	}

	// Open note - extract filename without emoji
	// Use folderPath which contains the full relative path
	fullPath := item.folderPath

	content, err := a.storage.ReadNote(fullPath)
	if err != nil {
		return a, nil
	}

	file, err := a.storage.OpenNote(fullPath)
	if err != nil {
		return a, nil
	}

	a.currentFile = file
	a.currentFileName = fullPath
	a.textarea.SetValue(string(content))
	a.lastSavedContent = string(content) // Track for autosave
	a.statusBar.SetFileInfo(fullPath, string(content))
	a.state = StateEditing

	// Add to recent files
	a.addToRecentFiles(fullPath)

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
	notes, folders, err := a.storage.ListNotes(a.currentFolder)
	if err != nil {
		return
	}

	// Combine folders and notes, folders first
	items := make([]list.Item, 0, len(folders)+len(notes))

	// Add ".." to go up if we're in a folder
	if a.currentFolder != "" {
		items = append(items, NoteItem{
			title:       "📁 ..",
			description: "Go up",
			isFolder:    true,
			folderPath:  "..",
		})
	}

	// Add folders
	for _, folder := range folders {
		items = append(items, NoteItem{
			title:       "📁 " + folder.Name,
			description: "Folder",
			isFolder:    true,
			folderPath:  folder.Path,
		})
	}

	// Add notes
	for _, note := range notes {
		displayTitle := note.Title
		if a.currentFolder != "" {
			// Show just the filename, not full path
			displayTitle = note.Title[len(a.currentFolder)+1:]
		}
		items = append(items, NoteItem{
			title:       "📝 " + displayTitle,
			description: note.Description,
			isFolder:    false,
			folderPath:  note.Title, // Store full path here
		})
	}

	a.list.SetItems(items)

	// Update list title to show current folder
	if a.currentFolder != "" {
		a.list.Title = "📁 " + a.currentFolder
	} else {
		a.list.Title = "All Notes"
	}
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

	// Save to config
	a.config.SetRecentFiles(a.recentFiles)
	a.config.Save()

	// Update homepage
	a.homepage.SetRecentFiles(a.recentFiles)
}

// removeFromRecentFiles removes a file from recent files list
func (a *App) removeFromRecentFiles(filename string) {
	for i, f := range a.recentFiles {
		if f == filename {
			a.recentFiles = append(a.recentFiles[:i], a.recentFiles[i+1:]...)
			break
		}
	}

	// Save to config
	a.config.SetRecentFiles(a.recentFiles)
	a.config.Save()

	// Update homepage
	a.homepage.SetRecentFiles(a.recentFiles)
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

func newSearchInput(theme styles.Theme) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search notes..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	cursorStyle := lipgloss.NewStyle().Foreground(theme.Cursor)
	textStyle := lipgloss.NewStyle().Foreground(theme.Foreground)
	placeholderStyle := lipgloss.NewStyle().Foreground(theme.Muted)

	ti.Cursor.Style = cursorStyle
	ti.Prompt = "/"
	ti.PromptStyle = cursorStyle
	ti.TextStyle = textStyle
	ti.PlaceholderStyle = placeholderStyle

	return ti
}

// performSearch filters notes using fuzzy matching
func (a *App) performSearch(query string) {
	if query == "" {
		// Show all notes
		a.list.SetItems(a.allNotes)
		return
	}

	// Use fuzzy matching
	var matches []list.Item
	for _, item := range a.allNotes {
		if noteItem, ok := item.(NoteItem); ok {
			// Search in both title and folder path
			searchText := noteItem.title
			if noteItem.folderPath != "" {
				searchText = noteItem.folderPath
			}
			// Remove emoji prefix for search
			if len(searchText) > 4 && (searchText[:4] == "📝 " || searchText[:4] == "📁 ") {
				searchText = searchText[4:]
			}

			// Simple contains search (fuzzy would require importing the library)
			if containsIgnoreCase(searchText, query) {
				matches = append(matches, item)
			}
		}
	}
	a.list.SetItems(matches)
}

// containsIgnoreCase checks if str contains substr (case-insensitive)
func containsIgnoreCase(str, substr string) bool {
	return len(substr) == 0 ||
		(len(str) >= len(substr) &&
			(len(substr) == 0 ||
				findSubstringIgnoreCase(str, substr)))
}

func findSubstringIgnoreCase(str, substr string) bool {
	strLower := toLower(str)
	substrLower := toLower(substr)
	for i := 0; i <= len(strLower)-len(substrLower); i++ {
		if strLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	// Simple ASCII lowercase
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
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
	notes, folders, _ := store.ListNotes("")

	// Combine folders and notes
	items := make([]list.Item, 0, len(folders)+len(notes))

	// Add folders first
	for _, folder := range folders {
		items = append(items, NoteItem{
			title:       "📁 " + folder.Name,
			description: "Folder",
			isFolder:    true,
			folderPath:  folder.Path,
		})
	}

	// Add notes
	for _, note := range notes {
		items = append(items, NoteItem{
			title:       "📝 " + note.Title,
			description: note.Description,
			isFolder:    false,
		})
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

func (a *App) togglePreview() tea.Cmd {
	if a.isPreview {
		a.isPreview = false
		a.previewViewport = viewport.Model{}
		return nil
	}

	// Calculate viewport dimensions
	width := a.width - 4
	height := a.height - 8
	if width < 40 {
		width = 78
	}
	if height < 10 {
		height = 20
	}

	// Use raw content for proper markdown rendering
	content := a.textarea.Value()

	// Create viewport with themed border
	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.Border).
		PaddingRight(2)

	// Create glamour renderer with word wrap
	glamourWidth := width - vp.Style.GetHorizontalFrameSize() - 2
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourWidth),
	)
	if err != nil {
		return nil
	}

	// Render markdown content
	rendered, err := renderer.Render(content)
	if err != nil {
		return nil
	}

	vp.SetContent(rendered)
	a.previewViewport = vp
	a.isPreview = true

	return nil
}
