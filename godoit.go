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
	modeConfirmDeleteAll
	modeEdit
)

type model struct {
	todos          []string
	cursor         int
	mode           mode
	textInput      textinput.Model
	status         string
	width          int
	height         int
	confirmIdx     int
	editIdx        int
	priorityInput  int
	prioritySelect bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type a todo and press Enter"
	ti.CharLimit = 256
	ti.Width = 50

	return model{
		todos:          loadTodos(),
		cursor:         0,
		mode:           modeView,
		textInput:      ti,
		status:         "Press 'a' to add, 'd' to delete, 'r' to reload, 'q' to quit.",
		priorityInput:  1,
		prioritySelect: false,
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

			if !strings.Contains(line, "[urgent]") && !strings.Contains(line, "[medium]") && !strings.Contains(line, "[low]") {
				id := ""
				if idx := strings.LastIndex(line, " #"); idx != -1 {
					id = line[idx:]
					line = line[:idx]
				}
				line = "[ ] " + strings.TrimSpace(line) + " [medium]" + id
			}
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
				m.priorityInput = 1
				m.prioritySelect = false
				return m, m.textInput.Focus()
			case "d":
				if len(m.todos) > 0 {
					m.mode = modeConfirmDelete
					m.confirmIdx = m.cursor
					m.status = fmt.Sprintf("Delete todo #%d? (y/n)", m.confirmIdx+1)
				}
			case "D":
				if len(m.todos) > 0 {
					m.mode = modeConfirmDeleteAll
					m.status = "Are you sure you want to delete ALL todos? (y/n)"
				}
			case "e":
				if len(m.todos) > 0 {
					m.mode = modeEdit
					m.editIdx = m.cursor
					todo := m.todos[m.editIdx]
					if idx := strings.LastIndex(todo, " #"); idx != -1 {
						todo = todo[4:idx]
					}
					m.textInput.SetValue(todo)
					m.textInput.Focus()
					m.status = "Editing todo"
					return m, m.textInput.Focus()
				}
			case " ":
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					done := strings.HasPrefix(todo, "[x]")
					rest := todo[3:]
					if strings.HasPrefix(rest, " ") {
						rest = rest[1:]
					}
					if done {
						m.todos[m.cursor] = "[ ] " + rest
					} else {
						m.todos[m.cursor] = "[x] " + rest
					}
					saveTodos(m.todos)
					m.status = "Toggled completion"
				}
			case "r":
				m.todos = loadTodos()
				m.status = "Reloaded todos from file"
			}

		case modeAdd:

			if m.prioritySelect {

				switch k {
				case "left":
					if m.priorityInput > 0 {
						m.priorityInput--
					}
				case "right":
					if m.priorityInput < 2 {
						m.priorityInput++
					}
				case "enter":
					val := strings.TrimSpace(m.textInput.Value())
					if val != "" {
						id := strconv.FormatInt(time.Now().UnixNano(), 10)[8:]
						var prio string
						switch m.priorityInput {
						case 0:
							prio = "[urgent]"
						case 1:
							prio = "[medium]"
						case 2:
							prio = "[low]"
						}
						entry := "[ ] " + val + " " + prio + " #" + id
						m.todos = append(m.todos, entry)
						saveTodos(m.todos)
						m.status = "Added todo"
					} else {
						m.status = "Empty todo not added"
					}
					m.mode = modeView
					m.prioritySelect = false
				case "esc":
					m.status = "Add cancelled"
					m.mode = modeView
					m.prioritySelect = false
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if k == "enter" {
				m.prioritySelect = true
				m.status = "Select priority: ←/→ and Enter (urgent, medium, low)"
			} else if k == "esc" {
				m.status = "Add cancelled"
				m.mode = modeView
			}
			return m, cmd

		case modeEdit:

			if m.prioritySelect {

				switch k {
				case "left":
					if m.priorityInput > 0 {
						m.priorityInput--
					}
				case "right":
					if m.priorityInput < 2 {
						m.priorityInput++
					}
				case "enter":
					val := strings.TrimSpace(m.textInput.Value())
					if val != "" && m.editIdx >= 0 && m.editIdx < len(m.todos) {
						old := m.todos[m.editIdx]
						id := ""
						prefix := "[ ] "
						if strings.HasPrefix(old, "[x]") {
							prefix = "[x] "
						}
						if idx := strings.LastIndex(old, " #"); idx != -1 {
							id = old[idx:]
						}
						var prio string
						switch m.priorityInput {
						case 0:
							prio = "[urgent]"
						case 1:
							prio = "[medium]"
						case 2:
							prio = "[low]"
						}
						m.todos[m.editIdx] = prefix + val + " " + prio + id
						saveTodos(m.todos)
						m.status = "Todo edited"
					} else {
						m.status = "Edit cancelled or empty"
					}
					m.mode = modeView
					m.prioritySelect = false
				case "esc":
					m.status = "Edit cancelled"
					m.mode = modeView
					m.prioritySelect = false
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if k == "enter" {

				old := m.todos[m.editIdx]
				if strings.Contains(old, "[urgent]") {
					m.priorityInput = 0
				} else if strings.Contains(old, "[medium]") {
					m.priorityInput = 1
				} else {
					m.priorityInput = 2
				}
				m.prioritySelect = true
				m.status = "Select priority: ←/→ and Enter (urgent, medium, low)"
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
		case modeConfirmDeleteAll:
			switch k {
			case "y", "enter":
				m.todos = []string{}
				saveTodos(m.todos)
				m.status = "All todos deleted"
				m.mode = modeView
				m.cursor = 0
			case "n", "esc":
				m.mode = modeView
				m.status = "Delete all cancelled"
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
	doneStyle := lipgloss.NewStyle().Faint(true).Strikethrough(true)
	urStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333")).Bold(true)
	medStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	lowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CC44")).Bold(true)

	numCol := 4
	taskCol := 44
	prioCol := 10

	var b strings.Builder
	b.WriteString(headerStyle.Render(" Go-Do-It — Bubble Tea TUI ") + "\n\n")

	if len(m.todos) == 0 {
		b.WriteString("No todos yet — press 'a' to add one.\n\n")
	} else {

		b.WriteString(fmt.Sprintf("%-*s%-*s%-*s\n",
			numCol, "#", taskCol, "Todo", prioCol, "Priority"))
		b.WriteString(strings.Repeat("-", numCol+taskCol+prioCol) + "\n")
		for i, t := range m.todos {

			rowPrefix := "  "
			if i == m.cursor && m.mode == modeView {
				rowPrefix = cursorStyle.Render("> ")
			}

			display := t
			task := display
			idIdx := strings.LastIndex(display, " #")
			prio := "[medium]"
			prioIdx := strings.LastIndex(display, "] [")
			if prioIdx == -1 {
				if strings.Contains(display, "[urgent]") {
					prio = "[urgent]"
				} else if strings.Contains(display, "[low]") {
					prio = "[low]"
				}
			} else {
				prioStart := strings.LastIndex(display, "[")
				prioEnd := strings.LastIndex(display, "]")
				if prioStart != -1 && prioEnd != -1 && prioEnd > prioStart {
					prio = display[prioStart : prioEnd+1]
				}
			}
			if idIdx != -1 {
				task = display[:idIdx]
			}
			if pidx := strings.LastIndex(task, "["); pidx != -1 {
				task = strings.TrimSpace(task[:pidx])
			}
			if len(task) > taskCol {
				task = task[:taskCol-3] + "..."
			}

			if strings.HasPrefix(display, "[x]") {
				task = doneStyle.Render(task)
			}

			prioLabel := ""
			switch prio {
			case "[urgent]":
				prioLabel = urStyle.Render("[urgent]")
			case "[medium]":
				prioLabel = medStyle.Render("[medium]")
			case "[low]":
				prioLabel = lowStyle.Render("[low]")
			}

			b.WriteString(fmt.Sprintf("%s%-*d%-*s%-*s\n",
				rowPrefix,
				numCol, i+1,
				taskCol, task,
				prioCol, prioLabel,
			))
		}
		b.WriteString("\n")
	}

	switch m.mode {
	case modeAdd:
		if m.prioritySelect {
			b.WriteString("Select priority: ←/→ and Enter (urgent, medium, low)\n")
			prioNames := []string{"[urgent]", "[medium]", "[low]"}
			styles := []lipgloss.Style{urStyle, medStyle, lowStyle}
			for i, name := range prioNames {
				if i == m.priorityInput {
					b.WriteString(styles[i].Bold(true).Underline(true).Render(name) + " ")
				} else {
					b.WriteString(styles[i].Render(name) + " ")
				}
			}
			b.WriteString("\n")
		} else {
			b.WriteString("Add mode — press Enter to continue, Esc to cancel\n")
			b.WriteString(m.textInput.View() + "\n")
		}
	case modeEdit:
		if m.prioritySelect {
			b.WriteString("Select priority: ←/→ and Enter (urgent, medium, low)\n")
			prioNames := []string{"[urgent]", "[medium]", "[low]"}
			styles := []lipgloss.Style{urStyle, medStyle, lowStyle}
			for i, name := range prioNames {
				if i == m.priorityInput {
					b.WriteString(styles[i].Bold(true).Underline(true).Render(name) + " ")
				} else {
					b.WriteString(styles[i].Render(name) + " ")
				}
			}
			b.WriteString("\n")
		} else {
			b.WriteString("Edit mode — press Enter to continue, Esc to cancel\n")
			b.WriteString(m.textInput.View() + "\n")
		}
	case modeConfirmDelete:
		b.WriteString(m.status + "\n")
	case modeConfirmDeleteAll:
		b.WriteString(m.status + "\n")
	}

	b.WriteString("\n")
	b.WriteString(statusStyle.Render(m.status))
	b.WriteString("\n\n")
	b.WriteString("Controls: j/down k/up a:add d:delete D:delete-all e:edit <space>:toggle r:reload q:quit\n")

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
