package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecuteSingle runs a command that does not involve pipes
func ExecuteSingle(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Printf("myshell: command not found: %s\n", args[0])
			// Fallback to Groq AI if available, otherwise local suggestions
			if os.Getenv("GROQ_API_KEY") != "" {
				HandleAICommand(strings.Join(args, " "))
			} else {
				SuggestCommand(args[0])
			}
		} else {
			fmt.Fprintf(os.Stderr, "myshell: %v\n", err)
		}
	}
}
