package main

import (
	"fmt"
	"os"
	"os/exec"
)

// ExecutePipeline handles sequences of commands connected by pipes
func ExecutePipeline(commands []string) {
	var cmds []*exec.Cmd

	for _, cmdStr := range commands {
		args := ParseCommand(cmdStr)
		if len(args) == 0 {
			continue
		}

		// Assuming no built-ins in the pipeline for simplicity
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}

	if len(cmds) == 0 {
		return
	}

	// Connect stdout to stdin for contiguous commands
	for i := 0; i < len(cmds)-1; i++ {
		stdout, err := cmds[i].StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "myshell: pipe error: %v\n", err)
			return
		}
		cmds[i+1].Stdin = stdout
	}

	// Connect last command's output to the shell's stdout
	cmds[len(cmds)-1].Stdout = os.Stdout
	cmds[len(cmds)-1].Stderr = os.Stderr

	// Start all commands
	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "myshell: pipeline start error: %v\n", err)
			return
		}
	}

	// Wait for all to finish
	for _, cmd := range cmds {
		cmd.Wait() // Ignore individual exit statuses inside pipeline for now
	}
}
