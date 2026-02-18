# LibreNotes

A beautiful terminal-based note-taking application for developers who prefer the command line. Built with Go and the Charm ecosystem.

<!-- TODO: Add demo GIF here -->
<!-- ![Demo](demo.gif) -->

![Go Version](https://img.shields.io/badge/go-1.25.5-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- 📝 **Create & Edit Notes** - Full markdown support with live preview
- 🎨 **Beautiful Themes** - 3 built-in themes (Cyberpunk, Catppuccin Mocha, Catppuccin Latte)
- ⌨️ **Keyboard-First Navigation** - Vim-like shortcuts, no mouse required
- 👁️ **Live Preview** - Toggle between edit and rendered markdown (Ctrl+P)
- 📁 **Organized Storage** - Notes stored in `~/.librenotes`
- 🔄 **Recent Files** - Quick access to last 10 notes
- 📊 **Live Statistics** - Real-time word and character count
- 🎯 **Clean UI** - distraction-free writing environment

## Installation

### From Source

```bash
git clone https://github.com/rj-2006/librenotes.git
cd librenotes
go build -o librenotes ./cmd/librenotes/
./librenotes
```

### Using Go Install

```bash
go install github.com/rj-2006/librenotes/cmd/librenotes@latest
```

## Usage

Launch the application:
```bash
./librenotes
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Ctrl+N` | Create new note |
| `Ctrl+L` | List all notes |
| `Ctrl+S` | Save current note |
| `Ctrl+P` | Toggle preview mode |
| `Ctrl+T` | Change theme |
| `↑/↓` | Navigate menus/lists |
| `Enter` | Select item |
| `Esc` | Go back |
| `Ctrl+Q` | Quit application |

### Workflow

1. **Welcome Screen** - Navigate with arrow keys, select "New File" or "View List"
2. **Create Note** - Press `Ctrl+N`, enter filename, start writing
3. **Preview** - Press `Ctrl+P` to see rendered markdown
4. **Save** - Press `Ctrl+S` to save and return to welcome screen
5. **Browse** - Press `Ctrl+L` to see all notes, select with Enter

## Themes

Three beautiful themes included:

- **Cyberpunk** (default) - Neon cyan & magenta on dark background
- **Catppuccin Mocha** - Warm dark theme with soft colors  
- **Catppuccin Latte** - Light theme for bright environments

Switch themes anytime with `Ctrl+T` from the welcome screen.

## Tech Stack

- **Go 1.25.5** - Core language
- **Bubble Tea** - TUI framework (Elm architecture)
- **Lipgloss** - Styling and layout
- **Glamour** - Markdown rendering
- **Bubbles** - Pre-built UI components

## Architecture

```
librenotes/
├── cmd/librenotes/      # Application entry point
├── internal/
│   ├── app/            # Main application logic & state management
│   ├── config/         # Configuration management
│   ├── storage/        # File I/O operations
│   └── ui/             # UI components & themes
├── go.mod              # Go module definition
└── Makefile            # Build automation
```

## Storage

Notes are stored as markdown files in `~/.librenotes/`:
- Plain text format (portable)
- Automatic `.md` extension
- Filename-based organization

## Development

### Build
```bash
make build
# or
go build -o librenotes ./cmd/librenotes/
```

### Run
```bash
make run
# or
./librenotes
```

### Clean
```bash
make clean
```

## License

MIT License - see [LICENSE](LICENSE) file.

## Acknowledgments

- [Charm](https://charm.sh/) - For the amazing TUI ecosystem
- [Catppuccin](https://catppuccin.com/) - For the beautiful color palettes

---

Made with ❤️ for the terminal lovers
