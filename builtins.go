package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var history []string
var historyFile string

func init() {
	home, err := os.UserHomeDir()
	if err == nil {
		historyFile = filepath.Join(home, ".myshell_history")
	}
}

// IsBuiltin checks if a command is a built-in shell command
func IsBuiltin(cmd string) bool {
	switch cmd {
	case "cd", "exit", "history":
		return true
	}
	return false
}

// RunBuiltin executes a built-in command
func RunBuiltin(args []string) {
	switch args[0] {
	case "cd":
		if len(args) < 2 {
			home, _ := os.UserHomeDir()
			os.Chdir(home)
			return
		}
		err := os.Chdir(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "myshell: cd: %s: %v\n", args[1], err)
		}
	case "exit":
		SaveHistory()
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "history":
		for i, cmd := range history {
			fmt.Printf("%d  %s\n", i+1, cmd)
		}
	}
}

// InitHistory loads command history from file
func InitHistory() {
	if historyFile == "" {
		return
	}
	data, err := os.ReadFile(historyFile)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line != "" {
				history = append(history, line)
			}
		}
	}
}

// AddHistory adds a new command to history with a timestamp
func AddHistory(cmd string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] %s", timestamp, cmd)
	history = append(history, entry)
}

// SaveHistory writes the command history back to the file
func SaveHistory() {
	if historyFile == "" {
		return
	}
	data := strings.Join(history, "\n") + "\n"
	os.WriteFile(historyFile, []byte(data), 0644)
}
