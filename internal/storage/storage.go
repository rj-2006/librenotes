package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Note represents a single note file
type Note struct {
	Title       string
	Path        string
	ModTime     time.Time
	Description string
}

// Storage handles all file operations
type Storage struct {
	vaultDir string
}

// NewStorage creates a new Storage instance
func NewStorage(vaultDir string) *Storage {
	return &Storage{vaultDir: vaultDir}
}

// Initialize creates the vault directory if it doesn't exist
func (s *Storage) Initialize() error {
	return os.MkdirAll(s.vaultDir, 0750)
}

// GetVaultDir returns the vault directory path
func (s *Storage) GetVaultDir() string {
	return s.vaultDir
}

// ListNotes returns all notes in the vault
func (s *Storage) ListNotes() ([]Note, error) {
	entries, err := os.ReadDir(s.vaultDir)
	if err != nil {
		return nil, fmt.Errorf("can't read from directory: %w", err)
	}

	var notes []Note
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		modTime := info.ModTime()
		notes = append(notes, Note{
			Title:       entry.Name(),
			Path:        filepath.Join(s.vaultDir, entry.Name()),
			ModTime:     modTime,
			Description: fmt.Sprintf("Modified: %s", modTime.Format("2006-01-02 15:04")),
		})
	}

	return notes, nil
}

// ReadNote reads the content of a note
func (s *Storage) ReadNote(filename string) ([]byte, error) {
	filepath := filepath.Join(s.vaultDir, filename)
	return os.ReadFile(filepath)
}

// OpenNote opens a note for reading/writing
func (s *Storage) OpenNote(filename string) (*os.File, error) {
	filepath := filepath.Join(s.vaultDir, filename)
	return os.OpenFile(filepath, os.O_RDWR, 0644)
}

// CreateNote creates a new note file
func (s *Storage) CreateNote(filename string) (*os.File, error) {
	if filepath.Ext(filename) == "" {
		filename = filename + ".md"
	}

	filepath := filepath.Join(s.vaultDir, filename)

	// Check if file already exists
	if _, err := os.Stat(filepath); err == nil {
		return nil, fmt.Errorf("file already exists")
	}

	return os.Create(filepath)
}

// NoteExists checks if a note exists
func (s *Storage) NoteExists(filename string) bool {
	filepath := filepath.Join(s.vaultDir, filename)
	_, err := os.Stat(filepath)
	return err == nil
}

// SaveNote saves content to an open file
func (s *Storage) SaveNote(file *os.File, content string) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("can't truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("can't seek file: %w", err)
	}
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("can't write to file: %w", err)
	}
	return nil
}

// CloseNote closes a note file
func (s *Storage) CloseNote(file *os.File) error {
	return file.Close()
}

// DeleteNote deletes a note file
func (s *Storage) DeleteNote(filename string) error {
	filepath := filepath.Join(s.vaultDir, filename)
	return os.Remove(filepath)
}

// CountWords counts words in content
func CountWords(content string) int {
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

// CountChars counts characters in content
func CountChars(content string) int {
	return len([]rune(content))
}
