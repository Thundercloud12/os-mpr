package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GroqRequest struct {
	Model          string            `json:"model"`
	Messages       []Message         `json:"messages"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type AIResponse struct {
	Type       string  `json:"type"`
	Command    string  `json:"command"`
	Message    string  `json:"message"`
	Confidence float64 `json:"confidence"`
}

const systemPrompt = `You are an AI systems agent embedded inside a custom Go-based terminal (Smart Shell). You are accessed via the Groq LLM API and must behave deterministically, safely, and return strictly structured JSON.

## 🎯 PRIMARY OBJECTIVE
Given a user input string, classify and respond as one of:
1. command → generate a safe Linux command
2. explanation → explain a given command
3. suggestion → correct a mistyped command
4. history → infer a past command query
5. none → unsafe or unsupported request

## ⚠️ STRICT OUTPUT FORMAT (MANDATORY)
You MUST return ONLY valid JSON. No text outside JSON.
{
"type": "<command | explanation | suggestion | history | none>",
"command": "<string or null>",
"message": "<clear human-readable explanation>",
"confidence": <float between 0 and 1>
}

## 🧠 BEHAVIOR RULES
### 1. Natural Language → Command
* Convert plain English into a valid Linux command
* Prefer simple, readable, commonly used commands
* Avoid overengineering (no complex one-liners unless necessary)

### 2. Command Explanation
Trigger: input starts with "explain"
* Do NOT generate a new command
* Only explain the given command clearly

### 3. Suggestions (Error Correction)
* Detect typos or near matches
* Suggest corrected command

### 4. History Queries
* If input refers to past commands
* Return best inferred command

### 5. 🚨 SAFETY CONSTRAINTS (CRITICAL)
NEVER generate or suggest dangerous commands, including but not limited to:
* rm -rf /
* shutdown / reboot
* mkfs / disk formatting
* fork bombs (:(){ :|:& };:)
* destructive wildcard deletions
If detected: return type "none".

### 6. EXECUTION AWARENESS
* Assume commands will be executed on a real Linux system
* Prefer safe flags
* Avoid irreversible actions

### 7. SIMPLICITY RULE
* Prefer: ls | grep file
* Over: find . -type f -name "*.file" (unless explicitly required)

## 🔌 GROQ API USAGE CONTEXT
You are being called via Groq Chat Completions API.
* Always respond in pure JSON (no markdown, no backticks)
* Keep responses concise and structured
* Do not include explanations outside the JSON "message" field

## 🧪 EDGE CASE HANDLING
If input is:
* empty → return type "none"
* unrelated to shell usage → return type "none"
* ambiguous → choose safest interpretation.`

func executeAICommand(cmdStr string) {
	commands := ParsePipes(cmdStr)
	if len(commands) == 0 {
		return
	}
	if len(commands) == 1 {
		args := ParseCommand(commands[0])
		if len(args) == 0 {
			return
		}
		if IsBuiltin(args[0]) {
			RunBuiltin(args)
		} else if IsCustom(args[0]) {
			RunCustom(args)
		} else {
			ExecuteSingle(args)
		}
	} else {
		ExecutePipeline(commands)
	}
}

func HandleAICommand(query string) {
	autoExec := false
	queryTrimmed := strings.TrimSpace(query)
	if strings.HasSuffix(queryTrimmed, "-exec") {
		autoExec = true
		query = strings.TrimSpace(strings.TrimSuffix(queryTrimmed, "-exec"))
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: GROQ_API_KEY environment variable is not set.")
		return
	}

	reqBody := GroqRequest{
		Model: "llama-3.1-8b-instant", // Fast and capable for this task
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: query},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error calling Groq API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "API Error: %s\n", string(body))
		return
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing Groq response:", err)
		return
	}

	if len(groqResp.Choices) == 0 {
		fmt.Fprintln(os.Stderr, "No response from Groq.")
		return
	}

	resultText := groqResp.Choices[0].Message.Content

	var aiResp AIResponse
	if err := json.Unmarshal([]byte(resultText), &aiResp); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing AI JSON output:", err)
		return
	}

	// Output Formatting
	fmt.Printf("\033[36m❖ AI Assistant [%s] (%.0f%% confidence)\033[0m\n", aiResp.Type, aiResp.Confidence*100)

	if aiResp.Message != "" {
		fmt.Printf("\033[33m%s\033[0m\n", aiResp.Message)
	}

	if aiResp.Type == "command" && aiResp.Command != "" {
		if autoExec {
			fmt.Printf("\033[32mExecuting:\033[0m %s\n", aiResp.Command)
			executeAICommand(aiResp.Command)
		} else {
			fmt.Printf("\033[32m%s\033[0m\n", aiResp.Command)
		}
	} else if aiResp.Type == "suggestion" && aiResp.Command != "" {
		if autoExec {
			fmt.Printf("\033[32mExecuting Suggestion:\033[0m %s\n", aiResp.Command)
			executeAICommand(aiResp.Command)
		} else {
			fmt.Printf("\033[32mTry running:\033[0m %s\n", aiResp.Command)
		}
	} else if aiResp.Type == "history" && aiResp.Command != "" {
		fmt.Printf("\033[32mInferred Command:\033[0m %s\n", aiResp.Command)
	}
}
