# Go-Do-It

A simple terminal-based Todo List application written in Go, using the Bubble Tea TUI framework.

## Features

- Add, view, edit, and delete todos from a text-based interface
- Delete all todos at once, with confirmation
- Set a priority for each task: **urgent** (red), **medium** (yellow), or **low** (green)
- Todos are saved to a local file (`todolist.txt`)
- Table-like formatting for todos: number, todo text (with tick status), and priority columns
- Keyboard navigation and controls
- Modern, colorful TUI using Bubble Tea and Lip Gloss

## Controls

- `j` / `down arrow`: Move cursor down
- `k` / `up arrow`: Move cursor up
- `space`: Toggle completion (tick/untick)
- `a`: Add a new todo (then select priority with ←/→ and Enter)
- `d`: Delete the selected todo
- `D`: Delete all todos (with confirmation)
- `e`: Edit a todo (then select priority with ←/→ and Enter)
- `r`: Reload todos from file
- `q`: Quit the application

## Requirements

- Go 1.18 or newer

## Installation

1. Clone this repository or copy the files to your local machine.
2. Install dependencies:
   ```sh
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/bubbles/textinput
   go get github.com/charmbracelet/lipgloss
   ```
3. Build the program:
   ```sh
   go build -o godoit.exe godoit.go
   ```

## Usage

Run the program from your terminal:

```sh
./godoit.exe
```

Your todos will be saved in `todolist.txt` in the same directory.

---

### Recent Updates

- Added ability to delete all todos at once (press `D` in view mode, with confirmation prompt)
- Todos are now displayed in a table-like format with columns for number, todo (with tick status), and priority
- Removed multi-select feature for simplicity

## License

MIT License

---

Inspired by the Bubble Tea TUI framework by Charmbracelet.
