package app

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/rj-2006/librenotes/internal/config"
	"github.com/rj-2006/librenotes/internal/storage"
	"github.com/rj-2006/librenotes/internal/ui/components"
	"github.com/rj-2006/librenotes/internal/ui/styles"
)

type ViewState int

const (
	StateWelcome ViewState = iota
	StateList
	StateNewFile
	StateEditing
)

type App struct {
	// State
	state ViewState

	// Config and Storage
	config  *config.Config
	storage *storage.Storage
	theme   styles.Theme

	// UI Components
	newFileInput textinput.Model
	textarea     textarea.Model
	list         list.Model
	statusBar    *components.StatusBar
	helpBar      *components.HelpBar
	homepage     *components.Homepage

	// Current file
	currentFile     *os.File
	currentFileName string

	// Recent files tracking
	recentFiles []string

	// Window dimensions
	width  int
	height int
}

// NoteItem represents a note in the list
type NoteItem struct {
	title       string
	description string
}

// Title returns the item title for the list
func (i NoteItem) Title() string { return i.title }

// Description returns the item description for the list
func (i NoteItem) Description() string { return i.description }

// FilterValue returns the filter value for the list
func (i NoteItem) FilterValue() string { return i.title }
