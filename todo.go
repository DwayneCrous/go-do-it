package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

const todoFile = "todolist.txt"

func loadTodos() []Todo {
	f, err := os.Open(todoFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Todo{}
		}
		log.Fatal(err)
	}
	defer f.Close()

	var todos []Todo
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var todo Todo
		if err := json.Unmarshal([]byte(line), &todo); err == nil {
			todos = append(todos, todo)
		} else {
			// fallback: try to parse old format (if any)
			todo = parseLegacyTodo(line)
			todos = append(todos, todo)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return todos
}

// parseLegacyTodo tries to parse a legacy todo string into a Todo struct
func parseLegacyTodo(line string) Todo {
	todo := Todo{Priority: "medium"}
	if strings.HasPrefix(line, "[x]") {
		todo.Done = true
		line = strings.TrimSpace(line[3:])
	} else if strings.HasPrefix(line, "[ ]") {
		todo.Done = false
		line = strings.TrimSpace(line[3:])
	}
	// Priority
	if strings.Contains(line, "[urgent]") {
		todo.Priority = "urgent"
		line = strings.Replace(line, "[urgent]", "", 1)
	} else if strings.Contains(line, "[low]") {
		todo.Priority = "low"
		line = strings.Replace(line, "[low]", "", 1)
	}
	// Due date
	if atIdx := strings.Index(line, " @"); atIdx != -1 {
		end := atIdx + 12
		if end > len(line) {
			end = len(line)
		}
		todo.DueDate = strings.TrimSpace(line[atIdx+2 : end])
		line = line[:atIdx] + line[end:]
	}
	todo.Text = strings.TrimSpace(line)
	return todo
}
func saveTodos(todos []Todo) {
	f, err := os.Create(todoFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, t := range todos {
		b, err := json.Marshal(t)
		if err != nil {
			log.Fatal(err)
		}
		_, err = writer.WriteString(string(b) + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush()
}
