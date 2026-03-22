package main

import (
	"fmt"
	"strings"
)

var knownCommands = []string{
	"ls", "cd", "pwd", "whoami", "echo", "cat", "grep", "git", "history", "exit", "clear",
}

// SuggestCommand uses a simple Levenshtein distance to find the closest match
func SuggestCommand(input string) {
	bestMatch := ""
	minDistance := -1

	for _, cmd := range knownCommands {
		dist := levenshtein(input, cmd)
		if dist <= 2 { // Threshold for suggestion
			if minDistance == -1 || dist < minDistance {
				minDistance = dist
				bestMatch = cmd
			}
		}
	}

	if bestMatch != "" {
		fmt.Printf("\033[33mDid you mean: %s?\033[0m\n", bestMatch)
	} else {
		// Fallback to simple prefix match
		for _, cmd := range knownCommands {
			if strings.HasPrefix(cmd, input) {
				fmt.Printf("\033[33mDid you mean: %s?\033[0m\n", cmd)
				return
			}
		}
	}
}

// levenshtein distance calculation
func levenshtein(a, b string) int {
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
	}
	for i := 0; i <= len(a); i++ {
		d[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		d[0][j] = j
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = minInt(
				d[i-1][j]+1,
				d[i][j-1]+1,
				d[i-1][j-1]+cost,
			)
		}
	}
	return d[len(a)][len(b)]
}

func minInt(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
