package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		k := msg.String()

		switch m.mode {
		case modeView:
			switch k {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "h":
				m.mode = modeHelp
				m.status = "Help menu (press any key or 'esc' to return)"
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
				m.dueDateInput = ""
				m.dueDateSelect = false
				return m, m.textInput.Focus()
			case "d":
				if len(m.todos) > 0 {
					m.mode = modeConfirmDelete
					m.confirmIdx = m.cursor
					m.status = fmt.Sprintf("Delete todo #%d? (y/n)", m.confirmIdx+1)
				}
			case "u":
				if m.canUndo {
					if m.lastDeletedIndex < 0 || m.lastDeletedIndex > len(m.todos) {
						m.lastDeletedIndex = len(m.todos)
					}
					before := m.todos[:m.lastDeletedIndex]
					after := m.todos[m.lastDeletedIndex:]
					m.todos = append(append(before, m.lastDeletedTodo), after...)
					saveTodos(m.todos)
					m.status = "Undo: restored deleted todo"
					m.cursor = m.lastDeletedIndex
					m.canUndo = false
				} else {
					m.status = "Nothing to undo"
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
		case modeHelp:
			m.mode = modeView
			m.status = "Exited help menu"
		case modeAdd:
			if m.dueDateSelect {
				switch k {
				case "enter":
					due := strings.TrimSpace(m.dueDateInput)
					if due != "" {
						_, err := time.Parse("2006-01-02", due)
						if err != nil {
							m.status = "Invalid date format. Use YYYY-MM-DD."
							return m, nil
						}
					}
					m.dueDateSelect = false
					m.prioritySelect = true
					m.status = "Select priority: ←/→ and Enter (urgent, medium, low)"
				case "esc":
					m.status = "Add cancelled"
					m.mode = modeView
					m.dueDateSelect = false
				case "backspace":
					if len(m.dueDateInput) > 0 {
						m.dueDateInput = m.dueDateInput[:len(m.dueDateInput)-1]
					}
				default:
					if len(k) == 1 && ((k[0] >= '0' && k[0] <= '9') || k[0] == '-') && len(m.dueDateInput) < 10 {
						m.dueDateInput += k
					}
				}
				return m, nil
			}
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
						due := ""
						if m.dueDateInput != "" {
							due = " @" + m.dueDateInput
						}
						entry := "[ ] " + val + due + " " + prio + " #" + id
						m.todos = append(m.todos, entry)
						saveTodos(m.todos)
						m.status = "Added todo"
					} else {
						m.status = "Empty todo not added"
					}
					m.mode = modeView
					m.prioritySelect = false
					m.dueDateInput = ""
				case "esc":
					m.status = "Add cancelled"
					m.mode = modeView
					m.prioritySelect = false
					m.dueDateInput = ""
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if k == "enter" {
				m.dueDateSelect = true
				m.status = "Enter due date (YYYY-MM-DD) or leave blank and press Enter: "
			} else if k == "esc" {
				m.status = "Add cancelled"
				m.mode = modeView
			}
			return m, cmd
		case modeEdit:
			if m.dueDateSelect {
				switch k {
				case "enter":
					due := strings.TrimSpace(m.dueDateInput)
					if due != "" {
						_, err := time.Parse("2006-01-02", due)
						if err != nil {
							m.status = "Invalid date format. Use YYYY-MM-DD."
							return m, nil
						}
					}
					m.dueDateSelect = false
					m.prioritySelect = true
					m.status = "Select priority: ←/→ and Enter (urgent, medium, low)"
				case "esc":
					m.status = "Edit cancelled"
					m.mode = modeView
					m.dueDateSelect = false
				case "backspace":
					if len(m.dueDateInput) > 0 {
						m.dueDateInput = m.dueDateInput[:len(m.dueDateInput)-1]
					}
				default:
					if len(k) == 1 && ((k[0] >= '0' && k[0] <= '9') || k[0] == '-') && len(m.dueDateInput) < 10 {
						m.dueDateInput += k
					}
				}
				return m, nil
			}
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
						due := ""
						if m.dueDateInput != "" {
							due = " @" + m.dueDateInput
						} else {
							if atIdx := strings.Index(old, " @"); atIdx != -1 {
								end := atIdx + 12
								if end > len(old) {
									end = len(old)
								}
								old = old[:atIdx] + old[end:]
							} else if atIdx := strings.Index(old, "@"); atIdx != -1 {
								end := atIdx + 11
								if end > len(old) {
									end = len(old)
								}
								old = strings.TrimSpace(old[:atIdx] + old[end:])
							}
						}
						m.todos[m.editIdx] = prefix + val + due + " " + prio + id
						saveTodos(m.todos)
						m.status = "Todo edited"
					} else {
						m.status = "Edit cancelled or empty"
					}
					m.mode = modeView
					m.prioritySelect = false
					m.dueDateInput = ""
				case "esc":
					m.status = "Edit cancelled"
					m.mode = modeView
					m.prioritySelect = false
					m.dueDateInput = ""
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if k == "enter" {
				old := m.todos[m.editIdx]
				m.dueDateInput = ""
				if atIdx := strings.Index(old, " @"); atIdx != -1 {
					end := atIdx + 12
					if end > len(old) {
						end = len(old)
					}
					m.dueDateInput = strings.TrimSpace(old[atIdx+2 : end])
				}
				m.dueDateSelect = true
				m.status = "Enter due date (YYYY-MM-DD) or leave blank and press Enter: "
			} else if k == "esc" {
				m.status = "Edit cancelled"
				m.mode = modeView
			}
			return m, cmd
		case modeConfirmDelete:
			switch k {
			case "y", "enter":
				if m.confirmIdx >= 0 && m.confirmIdx < len(m.todos) {
					m.lastDeletedTodo = m.todos[m.confirmIdx]
					m.lastDeletedIndex = m.confirmIdx
					m.canUndo = true
					m.todos = append(m.todos[:m.confirmIdx], m.todos[m.confirmIdx+1:]...)
					saveTodos(m.todos)
					m.status = "Todo deleted (press 'u' to undo)"
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
				m.canUndo = false
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
