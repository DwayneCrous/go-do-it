package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		k := msg.String()

		switch m.mode {

		case modeView:
			switch k {
			case "t":
				m.mode = modeTagSearch
				m.tagSearchInput.SetValue("")
				m.tagSearchInput.Focus()
				m.status = "Tag search: type to filter tags. Press esc to return."
				return m, nil
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
				m.status = "Add a new todo. Type and press Enter."
			case "d":
				if len(m.todos) > 0 {
					m.mode = modeConfirmDelete
					m.confirmIdx = m.cursor
					m.status = "Delete this todo? (y/n)"
				}
			case "D":
				if len(m.todos) > 0 {
					m.mode = modeConfirmDeleteAll
					m.status = "Delete ALL todos? (y/n)"
				}
			case "e":
				if len(m.todos) > 0 {
					m.mode = modeEdit
					m.editIdx = m.cursor

					currentTodo := m.todos[m.cursor].Text
					m.textInput.SetValue(currentTodo)
					m.textInput.Focus()
					m.status = "Edit todo. Press Enter to continue."
				}
			case " ":
				if len(m.todos) > 0 {
					m.todos[m.cursor].Done = !m.todos[m.cursor].Done
					saveTodos(m.todos)
					m.status = "Toggled completion."
				}
			case "r":
				m.todos = loadTodos()
				m.status = "Todos reloaded."
			case "u":
				if m.canUndo {
					idx := m.lastDeletedIndex
					if idx < 0 || idx > len(m.todos) {
						idx = len(m.todos)
					}
					m.todos = append(m.todos[:idx], append([]Todo{m.lastDeletedTodo}, m.todos[idx:]...)...)
					saveTodos(m.todos)
					m.canUndo = false
					m.status = "Undo successful."
				}
			case "h":
				m.mode = modeHelp
			case "q":
				return m, tea.Quit
			}

		case modeTagSearch:
			var cmd tea.Cmd
			m.tagSearchInput, cmd = m.tagSearchInput.Update(msg)
			if k == "esc" {
				m.mode = modeView
				m.tagSearchInput.Blur()
				m.status = "Returned from tag search."
				return m, cmd
			}
			return m, cmd

		case modeAdd:
			var cmd tea.Cmd = nil

			if !m.dueDateSelect && !m.prioritySelect && !m.tagsSelect {
				m.textInput, cmd = m.textInput.Update(msg)
				if k == "enter" {
					val := strings.TrimSpace(m.textInput.Value())
					if val != "" {

						m.tempTodoText = val
						m.dueDateInput = ""
						m.priorityInput = 1
						m.tagsInput = ""
						m.dueDateSelect = true
						m.status = "Enter due date (YYYY-MM-DD) or leave blank and press Enter: "
						m.textInput.SetValue("")
						return m, cmd
					} else {
						m.status = "Empty todo not added."
						m.mode = modeView
						m.textInput.Blur()
						return m, cmd
					}
				}
				if k == "esc" {
					m.mode = modeView
					m.status = "Add cancelled."
					m.textInput.Blur()
					return m, cmd
				}
				return m, cmd
			}

			if m.dueDateSelect {
				m.textInput, cmd = m.textInput.Update(msg)
				if k == "enter" {
					m.dueDateInput = strings.TrimSpace(m.textInput.Value())
					m.dueDateSelect = false
					m.prioritySelect = true
					m.status = "Select priority with ←/→, then press Enter"
					m.textInput.SetValue("")
					return m, cmd
				}
				if k == "esc" {
					m.dueDateSelect = false
					m.mode = modeView
					m.status = "Add cancelled."
					m.textInput.Blur()
					return m, cmd
				}
				return m, cmd
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
					m.prioritySelect = false
					m.tagsSelect = true
					m.status = "Enter tags (comma separated) or leave blank and press Enter: "
					m.textInput.SetValue("")
					return m, nil
				case "esc":
					m.prioritySelect = false
					m.mode = modeView
					m.status = "Add cancelled."
					m.textInput.Blur()
					return m, nil
				}
				return m, nil
			}

			if m.tagsSelect {
				m.textInput, cmd = m.textInput.Update(msg)
				if k == "enter" {
					m.tagsInput = strings.TrimSpace(m.textInput.Value())
					priority := "medium"
					if m.priorityInput == 0 {
						priority = "urgent"
					} else if m.priorityInput == 2 {
						priority = "low"
					}
					tags := []string{}
					if m.tagsInput != "" {
						tags = strings.Split(m.tagsInput, ",")
						for i := range tags {
							tags[i] = strings.TrimSpace(tags[i])
						}
					}
					m.todos = append(m.todos, Todo{
						Text:     m.tempTodoText,
						DueDate:  m.dueDateInput,
						Priority: priority,
						Tags:     tags,
						Done:     false,
					})
					saveTodos(m.todos)
					m.status = "Todo added!"
					m.mode = modeView
					m.tagsSelect = false
					m.textInput.Blur()
					return m, cmd
				}
				if k == "esc" {
					m.tagsSelect = false
					m.mode = modeView
					m.status = "Add cancelled."
					m.textInput.Blur()
					return m, cmd
				}
				return m, cmd
			}

		case modeEdit:
			var cmd tea.Cmd

			if !m.dueDateSelect && !m.prioritySelect && !m.tagsSelect {
				m.textInput, cmd = m.textInput.Update(msg)
				if k == "enter" {
					val := strings.TrimSpace(m.textInput.Value())
					if val != "" && m.editIdx >= 0 && m.editIdx < len(m.todos) {

						m.tempTodoText = val

						m.dueDateInput = m.todos[m.editIdx].DueDate

						switch m.todos[m.editIdx].Priority {
						case "urgent":
							m.priorityInput = 0
						case "low":
							m.priorityInput = 2
						default:
							m.priorityInput = 1
						}
						m.tagsInput = strings.Join(m.todos[m.editIdx].Tags, ", ")
						m.dueDateSelect = true
						m.status = "Enter due date (YYYY-MM-DD) or leave blank and press Enter: "
						m.textInput.SetValue(m.dueDateInput)
						return m, cmd
					} else {
						m.status = "Edit cancelled or empty"
						m.mode = modeView
						m.textInput.Blur()
						return m, cmd
					}
				}
				if k == "esc" {
					m.status = "Edit cancelled"
					m.mode = modeView
					m.textInput.Blur()
					return m, cmd
				}
				return m, cmd
			}

			if m.dueDateSelect {
				m.textInput, cmd = m.textInput.Update(msg)
				if k == "enter" {
					m.dueDateInput = strings.TrimSpace(m.textInput.Value())
					m.dueDateSelect = false
					m.prioritySelect = true
					m.status = "Select priority with ←/→, then press Enter"
					return m, cmd
				}
				if k == "esc" {
					m.dueDateSelect = false
					m.mode = modeView
					m.status = "Edit cancelled"
					m.textInput.Blur()
					return m, cmd
				}
				return m, cmd
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
					m.prioritySelect = false
					m.tagsSelect = true
					m.status = "Enter tags (comma separated) or leave blank and press Enter: "
					m.textInput.SetValue(m.tagsInput)
					return m, nil
				case "esc":
					m.prioritySelect = false
					m.mode = modeView
					m.status = "Edit cancelled"
					m.textInput.Blur()
					return m, nil
				}
				return m, nil
			}

			if m.tagsSelect {
				m.textInput, cmd = m.textInput.Update(msg)
				if k == "enter" {
					m.tagsInput = strings.TrimSpace(m.textInput.Value())

					priority := "medium"
					if m.priorityInput == 0 {
						priority = "urgent"
					} else if m.priorityInput == 2 {
						priority = "low"
					}
					tags := []string{}
					if m.tagsInput != "" {
						tags = strings.Split(m.tagsInput, ",")
						for i := range tags {
							tags[i] = strings.TrimSpace(tags[i])
						}
					}

					m.todos[m.editIdx].Text = m.tempTodoText
					m.todos[m.editIdx].DueDate = m.dueDateInput
					m.todos[m.editIdx].Priority = priority
					m.todos[m.editIdx].Tags = tags
					saveTodos(m.todos)
					m.status = "Todo edited!"
					m.mode = modeView
					m.tagsSelect = false
					m.textInput.Blur()
					return m, cmd
				}
				if k == "esc" {
					m.tagsSelect = false
					m.mode = modeView
					m.status = "Edit cancelled"
					m.textInput.Blur()
					return m, cmd
				}
				return m, cmd
			}

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
				m.todos = []Todo{}
				saveTodos(m.todos)
				m.canUndo = false
				m.status = "All todos deleted"
				m.mode = modeView
				m.cursor = 0
			case "n", "esc":
				m.mode = modeView
				m.status = "Delete all cancelled"
			}

		case modeHelp:
			m.mode = modeView
			m.status = "Returned from help."
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

