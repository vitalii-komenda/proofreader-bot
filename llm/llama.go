package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

type LLama struct {
	Model       string
	Temperature float64
	MaxTokens   int
	Stream      bool
	Messages    []Message
	URL         string
}

func (llama *LLama) SendRequest(text string) (string, error) {
	requestBody := RequestBody{
		Model: llama.Model,
		Messages: append(llama.Messages, Message{
			Role:    "user",
			Content: text,
		}),
		Temperature: llama.Temperature,
		MaxTokens:   llama.MaxTokens,
		Stream:      llama.Stream,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(llama.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var responseBody ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(responseBody.Choices) > 0 {
		return responseBody.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no choices in response")
}
