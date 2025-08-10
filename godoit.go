package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const todoFile = "todolist.txt"

type mode int

const (
	modeView mode = iota
	modeAdd
	modeConfirmDelete
	modeEdit
)

type model struct {
	todos      []string
	cursor     int
	mode       mode
	textInput  textinput.Model
	status     string
	width      int
	height     int
	confirmIdx int
	editIdx    int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type a todo and press Enter"
	ti.CharLimit = 256
	ti.Width = 50

	return model{
		todos:     loadTodos(),
		cursor:    0,
		mode:      modeView,
		textInput: ti,
		status:    "Press 'a' to add, 'd' to delete, 'r' to reload, 'q' to quit.",
	}
}

func loadTodos() []string {
	f, err := os.Open(todoFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}
		}
		log.Fatal(err)
	}
	defer f.Close()

	var todos []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			todos = append(todos, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return todos
}

func saveTodos(todos []string) {
	f, err := os.Create(todoFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, t := range todos {
		_, err := writer.WriteString(t + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		k := msg.String()

		switch m.mode {
		case modeView:
			switch k {
			case "ctrl+c", "q":

				return m, tea.Quit
			case "j", "down":

				if m.cursor < len(m.todos)-1 {
					m.cursor++
				}
			case "k", "up":

				if m.cursor > 0 {
					m.cursor--
				}
			case "a":

				m.mode = modeAdd
				m.textInput.SetValue("")
				m.textInput.Focus()
				m.status = "Adding a new todo"
				return m, m.textInput.Focus()
			case "d":

				if len(m.todos) > 0 {
					m.mode = modeConfirmDelete
					m.confirmIdx = m.cursor
					m.status = fmt.Sprintf("Delete todo #%d? (y/n)", m.confirmIdx+1)
				}
			case "e":

				if len(m.todos) > 0 {
					m.mode = modeEdit
					m.editIdx = m.cursor

					todo := m.todos[m.editIdx]
					if idx := strings.LastIndex(todo, " #"); idx != -1 {
						todo = todo[:idx]
					}
					m.textInput.SetValue(todo)
					m.textInput.Focus()
					m.status = "Editing todo"
					return m, m.textInput.Focus()
				}
			case "r":

				m.todos = loadTodos()
				m.status = "Reloaded todos from file"
			}

		case modeAdd:

			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)

			if k == "enter" {

				val := strings.TrimSpace(m.textInput.Value())
				if val != "" {
					id := strconv.FormatInt(time.Now().UnixNano(), 10)[8:]
					entry := val + " #" + id
					m.todos = append(m.todos, entry)
					saveTodos(m.todos)
					m.status = "Added todo"
				} else {
					m.status = "Empty todo not added"
				}
				m.mode = modeView
			} else if k == "esc" {

				m.status = "Add cancelled"
				m.mode = modeView
			}
			return m, cmd

		case modeEdit:

			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if k == "enter" {
				val := strings.TrimSpace(m.textInput.Value())
				if val != "" && m.editIdx >= 0 && m.editIdx < len(m.todos) {

					old := m.todos[m.editIdx]
					id := ""
					if idx := strings.LastIndex(old, " #"); idx != -1 {
						id = old[idx:]
					}
					m.todos[m.editIdx] = val + id
					saveTodos(m.todos)
					m.status = "Todo edited"
				} else {
					m.status = "Edit cancelled or empty"
				}
				m.mode = modeView
			} else if k == "esc" {
				m.status = "Edit cancelled"
				m.mode = modeView
			}
			return m, cmd
		case modeConfirmDelete:

			switch k {
			case "y", "enter":

				if m.confirmIdx >= 0 && m.confirmIdx < len(m.todos) {
					m.todos = append(m.todos[:m.confirmIdx], m.todos[m.confirmIdx+1:]...)
					saveTodos(m.todos)
					m.status = "Todo deleted"
					if m.cursor >= len(m.todos) && m.cursor > 0 {
						m.cursor--
					}
				}
				m.mode = modeView
			case "n", "esc":

				m.mode = modeView
				m.status = "Delete cancelled"
			}
		}

	case tea.WindowSizeMsg:

		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7CCB"))
	statusStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888"))

	var b strings.Builder
	b.WriteString(headerStyle.Render(" Go-Do-It — Bubble Tea TUI ") + "\n\n")

	if len(m.todos) == 0 {

		b.WriteString("No todos yet — press 'a' to add one.\n\n")
	} else {

		for i, t := range m.todos {
			prefix := "  "
			if i == m.cursor && m.mode == modeView {
				prefix = cursorStyle.Render("> ")
			}
			display := t
			if len(display) > 80 {
				display = display[:77] + "..."
			}
			b.WriteString(fmt.Sprintf("%s%d: %s\n", prefix, i+1, display))
		}
		b.WriteString("\n")
	}

	switch m.mode {
	case modeAdd:
		b.WriteString("Add mode — press Enter to save, Esc to cancel\n")
		b.WriteString(m.textInput.View() + "\n")
	case modeEdit:
		b.WriteString("Edit mode — press Enter to save, Esc to cancel\n")
		b.WriteString(m.textInput.View() + "\n")
	case modeConfirmDelete:
		b.WriteString(m.status + "\n")
	}

	b.WriteString("\n")
	b.WriteString(statusStyle.Render(m.status))
	b.WriteString("\n\n")
	b.WriteString("Controls: j/down k/up a:add d:delete e:edit r:reload q:quit\n")

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
