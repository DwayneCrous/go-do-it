package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const todoFile = "todolist.txt"

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
