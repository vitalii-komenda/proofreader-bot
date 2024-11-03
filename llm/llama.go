package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (llama *LLama) Init() LLM {
	if llama.Model == "" {
		llama.Model = "lmstudio-community/Meta-Llama-3.1-8B-Instruct-GGUF"
	}
	if llama.Temperature == 0 {
		llama.Temperature = 0.7
	}
	if llama.MaxTokens == 0 {
		llama.MaxTokens = -1
	}
	if llama.URL == "" {
		llama.URL = "http://localhost:1234/v1/chat/completions"
	}
	if len(llama.Messages) == 0 {
		llama.Messages = []Message{
			{
				Role: "system",
				Content: `You are proofreader. Users will be asking to correct the text. Correct them with no explanations. 
Format like this:
*Typos*: list of words with a typo
*Proofread*: $whole_corrected_text`,
			},
		}
	}
	return llama
}
