package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gookit/color"
)

const (
	apiKey       = "<YOUR_CHATGPT_API_KEY>"
	model        = "gpt-3.5-turbo"
	maxTokens    = 100
	temperature  = 0.5 // Adjust the temperature value (between 0.2 and 0.8) to control creativity vs. accuracy
	chatEndpoint = "https://api.openai.com/v1/chat/completions"
)

type ChatRequest struct {
	Model       string   `json:"model"`
	MaxTokens   int      `json:"max_tokens"`
	Temperature float32  `json:"temperature"`
	Prompt      []string `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func generateCode(input string) (string, error) {
	chatInput := []string{
		{"user", input},
		{"assistant", "```go\n"},
	}

	reqBody, err := json.Marshal(ChatRequest{
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		Prompt:      chatInput,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", chatEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return "", err
	}

	code := chatResponse.Choices[0].Message.Content

	// Wrap the generated code in code block if not already wrapped
	code = wrapInCodeBlock(code)

	return code, nil
}

func wrapInCodeBlock(code string) string {
	// Check if code is already wrapped in a code block
	if strings.HasPrefix(code, "```") && strings.HasSuffix(code, "```") {
		return code
	}

	// Wrap the code in a code block
	return "```\n" + code + "\n```"
}

func main() {
	fmt.Println("Welcome to Codebuddy!")
	fmt.Println("Enter your desired code request:")
	var input string
	fmt.Scanln(&input)

	code, err := generateCode(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Generated Code:")
	color.Style{color.FgGreen}.Printf(code)
}
