package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7CCB"))
	statusStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888"))
	doneStyle := lipgloss.NewStyle().Faint(true).Strikethrough(true)
	urStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333")).Bold(true)
	medStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	lowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CC44")).Bold(true)
	overdueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Underline(true)

	if m.mode == modeHelp {
		var b strings.Builder
		b.WriteString(headerStyle.Render(" Go-Do-It — Help Menu ") + "\n\n")
		b.WriteString("Keybindings:\n\n")
		b.WriteString("  j / ↓         Move cursor down\n")
		b.WriteString("  k / ↑         Move cursor up\n")
		b.WriteString("  a             Add a new todo\n")
		b.WriteString("  d             Delete selected todo\n")
		b.WriteString("  u             Undo last todo deletion\n")
		b.WriteString("  D             Delete all todos\n")
		b.WriteString("  e             Edit selected todo\n")
		b.WriteString("  <space>       Toggle completion\n")
		b.WriteString("  r             Reload todos from file\n")
		b.WriteString("  h             Show this help menu\n")
		b.WriteString("  q             Quit the application\n")
		b.WriteString("  esc/any key   Return to todo list\n")
		b.WriteString("\n")
		b.WriteString(statusStyle.Render("Press any key or 'esc' to return to your todos."))
		b.WriteString("\n")
		return b.String()
	}

	numCol := 4
	taskCol := 44
	dueCol := 12
	prioCol := 10

	space := " "
	sep := space

	var b strings.Builder
	b.WriteString(headerStyle.Render(" Go-Do-It — Bubble Tea TUI ") + "\n\n")

	if len(m.todos) == 0 {
		b.WriteString("No todos yet — press 'a' to add one.\n\n")
	} else {
		headerLine := fmt.Sprintf("%-*s%s%-*s%s%-*s%s%-*s",
			numCol, "#", sep,
			taskCol, "Todo", sep,
			dueCol, "Due Date", sep,
			prioCol, "Priority",
		)
		b.WriteString(headerLine + "\n")
		b.WriteString(strings.Repeat("-", len(headerLine)) + "\n")

		for i, t := range m.todos {
			rowPrefix := "  "
			if i == m.cursor && m.mode == modeView {
				rowPrefix = cursorStyle.Render("> ")
			}
			display := t
			task := display

			prio := "[medium]"
			if strings.Contains(display, "[urgent]") {
				prio = "[urgent]"
			} else if strings.Contains(display, "[low]") {
				prio = "[low]"
			}

			dueDate := ""
			if atIdx := strings.Index(display, " @"); atIdx != -1 {
				end := atIdx + 12
				if end > len(display) {
					end = len(display)
				}
				dueDate = strings.TrimSpace(display[atIdx+2 : end])
			}

			if idIdx := strings.LastIndex(display, " #"); idIdx != -1 {
				task = display[:idIdx]
			} else {
				task = display
			}

			if pidx := strings.LastIndex(task, "["); pidx != -1 {
				task = strings.TrimSpace(task[:pidx])
			}
			if atIdx := strings.Index(task, " @"); atIdx != -1 {
				task = strings.TrimSpace(task[:atIdx])
			}

			if len([]rune(task)) > taskCol {
				r := []rune(task)
				if len(r) > taskCol-3 {
					r = r[:taskCol-3]
				}
				task = string(r) + "..."
			}

			isDone := strings.HasPrefix(display, "[x]")
			isOverdue := false
			if dueDate != "" && !isDone {
				if due, err := time.Parse("2006-01-02", dueDate); err == nil {
					if due.Before(time.Now()) {
						isOverdue = true
					}
				}
			}
			if isDone {
				task = doneStyle.Render(task)
			} else if isOverdue {
				task = overdueStyle.Render(task)
			}

			var prioLabel string
			switch prio {
			case "[urgent]":
				prioLabel = urStyle.Render("[urgent]")
			case "[medium]":
				prioLabel = medStyle.Render("[medium]")
			case "[low]":
				prioLabel = lowStyle.Render("[low]")
			}
			dueLabel := dueDate
			if isOverdue && !isDone && dueDate != "" {
				dueLabel = overdueStyle.Render(dueDate)
			}

			row := fmt.Sprintf("%s%-*d%s%-*s%s%-*s%s%-*s",
				rowPrefix,
				numCol, i+1, sep,
				taskCol, task, sep,
				dueCol, dueLabel, sep,
				prioCol, prioLabel,
			)
			b.WriteString(row + "\n")
		}
		b.WriteString("\n")
	}

	switch m.mode {
	case modeAdd:
		if m.dueDateSelect {
			b.WriteString("Add mode — enter due date (YYYY-MM-DD) or leave blank and press Enter\n")
			b.WriteString("Due date: " + m.dueDateInput + "\n")
		} else if m.prioritySelect {
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
		if m.dueDateSelect {
			b.WriteString("Edit mode — enter due date (YYYY-MM-DD) or leave blank and press Enter\n")
			b.WriteString("Due date: " + m.dueDateInput + "\n")
		} else if m.prioritySelect {
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
	b.WriteString("Controls: j/down k/up a:add d:delete D:delete-all e:edit <space>:toggle r:reload u:undo h:help q:quit\n")

	return b.String()
}
