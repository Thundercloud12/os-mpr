package main

import "strings"

// ParsePipes splits the input by the pipe character '|'
func ParsePipes(input string) []string {
	parts := strings.Split(input, "|")
	var commands []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			commands = append(commands, trimmed)
		}
	}
	return commands
}

// ParseCommand splits a single command string by spaces into arguments
func ParseCommand(cmd string) []string {
	parts := strings.Fields(cmd)
	return parts
}
