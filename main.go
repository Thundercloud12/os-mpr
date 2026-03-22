package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	InitHistory()
	defer SaveHistory()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\033[32mmyshell>\033[0m ") // Green colored prompt
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		AddHistory(input)

		commands := ParsePipes(input)
		if len(commands) == 0 {
			continue
		}

		if len(commands) == 1 {
			args := ParseCommand(commands[0])
			if len(args) == 0 {
				continue
			}

			// Check built-ins
			if IsBuiltin(args[0]) {
				RunBuiltin(args)
				continue
			}

			// Check custom commands
			if IsCustom(args[0]) {
				RunCustom(args)
				continue
			}

			// Check AI command
			if args[0] == "ai" {
				if len(args) > 1 {
					HandleAICommand(strings.Join(args[1:], " "))
				} else {
					fmt.Println("Usage: ai <query>")
				}
				continue
			}

			ExecuteSingle(args)
		} else {
			ExecutePipeline(commands)
		}
	}
}
