# Go-Do-It

A simple terminal-based Todo List application written in Go, using the Bubble Tea TUI framework.

## Features

- Add, view, and delete todos from a text-based interface
- Todos are saved to a local file (`todolist.txt`)
- Keyboard navigation and controls
- Modern, colorful TUI using Bubble Tea and Lip Gloss

## Controls

- `j` / `down arrow`: Move cursor down
- `k` / `up arrow`: Move cursor up
- `space`: Toggle completion
- `a`: Add a new todo
- `d`: Delete the selected todo
- `r`: Reload todos from file
- `q`: Quit the application
- `e`: Edit a todo

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

## License

MIT License

---

Inspired by the Bubble Tea TUI framework by Charmbracelet.
