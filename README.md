# Go-Do-It

A simple terminal-based Todo List application written in Go, using the Bubble Tea TUI framework.

## Features

- Add, view, edit, and delete todos from a modern, colorful terminal interface
- **Undo delete**: Accidentally deleted a todo? Press `u` to restore the last deleted item
- **Help menu**: Press `h` to view a dedicated help screen with all keybindings
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
- `u`: Undo the last todo deletion
- `D`: Delete all todos (with confirmation)
- `e`: Edit a todo (edit text, due date, and priority)
- `r`: Reload todos from file
- `h`: Show the help menu with all keybindings
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

- **Help menu**: Press `h` to view a dedicated help screen with all keybindings
- **Undo delete**: Press `u` to restore the last deleted todo

## License

MIT License

---

Inspired by the Bubble Tea TUI framework by Charmbracelet.
