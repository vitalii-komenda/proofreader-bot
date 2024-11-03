package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAI struct {
	Model       string
	Temperature float64
	MaxTokens   int
	Stream      bool
	Messages    []Message
	URL         string
	Token       string
}

func (openai *OpenAI) SendRequest(text string) (string, error) {
	requestBody := RequestBody{
		Model: openai.Model,
		Messages: append(openai.Messages, Message{
			Role:    "user",
			Content: text,
		}),
		Temperature: openai.Temperature,
		MaxTokens:   openai.MaxTokens,
		Stream:      openai.Stream,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", openai.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openai.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	var responseBody ResponseBody
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(responseBody.Choices) > 0 {
		return responseBody.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no choices in response")
}

func (openai *OpenAI) Init() LLM {
	if openai.Model == "" {
		openai.Model = "gpt-4o-mini"
	}
	if openai.Temperature == 0 {
		openai.Temperature = 0.7
	}
	if openai.MaxTokens == 0 {
		openai.MaxTokens = -1
	}
	if openai.URL == "" {
		openai.URL = "https://api.openai.com/v1/chat/completions"
	}
	if len(openai.Messages) == 0 {
		openai.Messages = []Message{
			{
				Role: "system",
				Content: `You are proofreader. Users will be asking to correct the text. Correct them with no explanations. 
Format like this:
*Typos*: list of words with a typo
*Proofread*: $whole_corrected_text`,
			},
		}
	}
	return openai
}
