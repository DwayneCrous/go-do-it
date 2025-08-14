# Go-Do-It

A simple terminal-based Todo List application written in Go, using the Bubble Tea TUI framework.

## Features

- Add, view, edit, and delete todos from a modern, colorful terminal interface
- **Edit mode**: Edit any todo, including its text, due date, and priority
- **Due dates**: Assign an optional due date (YYYY-MM-DD) to each todo
- **Overdue highlighting**: Todos past their due date are shown in red (unless completed)
- **Priority selection**: Choose between **urgent** (red), **medium** (yellow), or **low** (green) for each task
- **Delete all**: Remove all todos at once, with confirmation
- **Reload**: Instantly reload todos from file without restarting
- **Persistent storage**: Todos are saved to a local file (`todolist.txt`)
- **Table-like formatting**: Todos are displayed with columns for number, task, due date, and priority
- **Keyboard navigation and controls**: Fast, Vim-like navigation and shortcuts
- Built with Bubble Tea, Bubbles, and Lip Gloss for a beautiful TUI

## Controls

- `j` / `down arrow`: Move cursor down
- `k` / `up arrow`: Move cursor up
- `space`: Toggle completion (tick/untick)
- `a`: Add a new todo (enter text, then due date, then select priority with ←/→ and Enter)
- `d`: Delete the selected todo
- `D`: Delete all todos (with confirmation)
- `e`: Edit a todo (edit text, due date, and priority)
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

- **Edit mode**: Edit todo text, due date, and priority
- **Due dates**: Add/edit due dates for todos; overdue tasks are highlighted
- **Improved add flow**: Add text, then due date, then priority
- **Reload**: Reload todos from file with `r`
- **Overdue highlighting**: Overdue tasks are shown in red
- **Table columns**: Now includes due date and priority columns
- **Delete all**: Press `D` to delete all todos (with confirmation)

## License

MIT License

---

Inspired by the Bubble Tea TUI framework by Charmbracelet.
