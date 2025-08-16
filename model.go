package main

import (
    "github.com/charmbracelet/bubbles/textinput"
)

type mode int

const (
    modeView mode = iota
    modeAdd
    modeConfirmDelete
    modeConfirmDeleteAll
    modeEdit
    modeHelp
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
    dueDateInput   string
    dueDateSelect  bool

    lastDeletedTodo  string
    lastDeletedIndex int
    canUndo          bool
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
        dueDateInput:   "",
        dueDateSelect:  false,
    }
}