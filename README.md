
# Go-Do-It

A simple terminal-based Todo List application written in Go, using the Bubble Tea TUI framework.

## Project Structure

* `model.go` — Data model and state
* `todo.go` — Todo file I/O and helpers
* `update.go` — All update logic (event handling)

## Features

* Add, view, edit, and delete todos from a modern, colorful terminal interface
* **Tags**: Assign tags to each todo for better organization and filtering
* **Undo delete**: Accidentally deleted a todo? Press `u` to restore the last deleted item
* **Help menu**: Press `h` to view a dedicated help screen with all keybindings
* **Edit mode**: Edit any todo, including its text, due date, priority, and tags
* **Due dates**: Assign an optional due date (YYYY-MM-DD) to each todo
* **Overdue highlighting**: Todos past their due date are shown in red (unless completed)
* **Priority selection**: Choose between **urgent** (red), **medium** (yellow), or **low** (green) for each task
* **Delete all**: Remove all todos at once, with confirmation
* **Reload**: Instantly reload todos from file without restarting
* **Persistent storage**: Todos are saved to a local file (`todolist.txt`)
* **Table-like formatting**: Todos are displayed with columns for number, task, due date, priority, and tags
* **Keyboard navigation and controls**: Fast, Vim-like navigation and shortcuts
* Tag Search: Press `t` to search and filter todos by tag in a dedicated tag search mode
* Built with Bubble Tea, Bubbles, and Lip Gloss for a beautiful TUI
* **Reload**: Instantly reload todos from file without restarting
* **Persistent storage**: Todos are saved to a local file (`todolist.txt`)

## Controls

* `j` / `down arrow`: Move cursor down
* `k` / `up arrow`: Move cursor up
* `space`: Toggle completion (tick/untick)
* `a`: Add a new todo (enter text, due date, priority, and tags)
* `d`: Delete the selected todo
* `u`: Undo the last todo deletion
* `D`: Delete all todos (with confirmation)
* `e`: Edit a todo (edit text, due date, priority, and tags)
* `r`: Reload todos from file
* `h`: Show the help menu with all keybindings
* `q`: Quit the application
* `t`: Tag search (filter todos by tag)
* `u`: Undo the last todo deletion
* `D`: Delete all todos (with confirmation)

## Requirements

* `h`: Show the help menu with all keybindings
* `q`: Quit the application

## Installation

### Requirements

* Go 1.18 or newer

1. Clone this repository or copy the files to your local machine.
2. Install dependencies:
   ```sh
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/bubbles/textinput
   go get github.com/charmbracelet/lipgloss
   ```

## Usage


```sh
go build -o godoit.exe
```

Run the program from your terminal:

./godoit.exe

### Recent Updates

* **Tag Search**: You can now search for todos by tags using the `t` keybinding
* **Tags**: You can now add tags to todos during add and edit flows

---

## License

MIT License

---

Inspired by the Bubble Tea TUI framework by Charmbracelet.
