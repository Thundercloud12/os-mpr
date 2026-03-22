package main

import (
	"fmt"
	"os"
	"os/exec"
)

// IsCustom checks if the command is a custom shortcut
func IsCustom(cmd string) bool {
	switch cmd {
	case "open", "myip":
		return true
	}
	return false
}

// RunCustom executes a custom shortcut
func RunCustom(args []string) {
	switch args[0] {
	case "open":
		// Example: open google -> opens google.com
		if len(args) > 1 && args[1] == "google" {
			fmt.Println("Opening Google in browser...")
			// Assuming Linux running xdg-open. On macOS it would be 'open', Windows 'start'
			err := exec.Command("xdg-open", "https://google.com").Start()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error opening browser:", err)
			}
		} else {
			fmt.Println("Usage: open google")
		}
	case "myip":
		// Prints public IP using curl
		fmt.Println("Fetching public IP...")
		cmd := exec.Command("curl", "-s", "ifconfig.me")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching IP:", err)
		}
		fmt.Println() // Add a trailing newline
	}
}
