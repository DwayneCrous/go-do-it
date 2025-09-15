package main

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func newTextInputModel() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search tags..."
	ti.CharLimit = 50
	ti.Width = 30
	return ti
}

type mode int

const (
	modeView mode = iota
	modeAdd
	modeConfirmDelete
	modeConfirmDeleteAll
	modeEdit
	modeHelp
	modeTagSearch
)

type Todo struct {
	Text     string
	Priority string
	DueDate  string
	Done     bool
	Tags     []string
}

type model struct {
	todos          []Todo
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
	dueDateInput   string
	dueDateSelect  bool
	tagsInput      string
	tagsSelect     bool
	tempTodoText   string
	tagSearchInput textinput.Model

	lastDeletedTodo  Todo
	lastDeletedIndex int
	canUndo          bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type a todo and press Enter"
	ti.CharLimit = 256
	ti.Width = 50

	return model{
		todos:            loadTodos(),
		cursor:           0,
		mode:             modeView,
		textInput:        ti,
		status:           "Welcome to Go-Do-It! Press 'a' to add a todo.",
		width:            0,
		height:           0,
		confirmIdx:       -1,
		editIdx:          -1,
		priorityInput:    1,
		prioritySelect:   false,
		dueDateInput:     "",
		dueDateSelect:    false,
		tagsInput:        "",
		tagsSelect:       false,
		tagSearchInput:   newTextInputModel(),
		lastDeletedTodo:  Todo{},
		lastDeletedIndex: -1,
		canUndo:          false,
	}
}

