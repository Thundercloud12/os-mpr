# Smart Mini Shell in Go

A clean, modular, and demo-ready CLI-based shell application written in Go. It supports standard command execution along with enhanced usability features like pipes, command history, custom shortcuts, and an embedded AI assistant powered by the Groq API.

## Features

- **Interactive Shell Loop:** Custom green prompt `myshell>`, graceful empty input handling.
- **Command Execution:** Parses space-delimited inputs and executes system binaries (`ls`, `pwd`) natively.
- **Pipe Support:** Connects standard I/O streams seamlessly between sequential commands (e.g., `ls | grep go`).
- **Built-in Commands:**
  - `cd <dir>`: Change the current working directory.
  - `history`: View previously executed commands with timestamps.
  - `exit`: Terminate the shell.
- **History Persistence:** Automatically saves and loads your command history to/from `~/.myshell_history`.
- **Smart Correction:** Levenshtein-distance typo detection for unknown commands (e.g., `gti` -> `git`).
- **Groq AI Assistant:** Built-in natural language to CLI translations using LLaMA models. 
  - The `-exec` flag: Append `-exec` to your natural language prompt and the AI will auto-generate and immediately execute the suggested command in your shell loop.
- **Custom Scripts:** Shortcuts like `open google` or `myip` directly integrated into the parser.

## Architecture

The source code is organized clearly for a highly modular feature addition experience:
- `main.go`: Entry point, input loop, and core routing logic.
- `parser.go`: Command parsing split logic mapping commands sequentially or breaking piped chunks.
- `executor.go`: Subprocess control and standard program loading.
- `builtins.go`: Native Go overrides simulating standard core utilities.
- `pipeline.go`: I/O orchestrator for executing sequential programs safely.
- `suggestions.go`: Autocomplete and correction logic implementation.
- `custom.go`: Custom workflow mappings.
- `ai.go`: Integration with Groq API.

## Prerequisites

- **Go 1.21+** installed on your system.
- *(Optional)* A valid [Groq API Key](https://console.groq.com/keys) to utilize the AI Assistant feature.

## Building and Running

1. Clone or navigate to the project directory:
   ```bash
   cd os-mpr
   ```

2. Compile the binary:
   ```bash
   go build -o myshell
   ```

3. Add your Groq API key:
   ```bash
   export GROQ_API_KEY="your_api_key_here"
   ```

4. Launch the shell:
   ```bash
   ./myshell
   ```

## Example Usage of AI Integration

You can easily translate plain English text directly into Linux commands inside the shell. Try opening `./myshell` and running:
```bash
myshell> ai check what process is taking the highest memory -exec
```

The AI will output and run the pipeline securely:
```bash
❖ AI Assistant [command] (100% confidence)
Show which process is taking the most memory. Sorted by memory usage in descending order and show the top result.
Executing: ps aux | sort -rn -k 4 | head -1
...
```
