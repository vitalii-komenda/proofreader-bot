package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"net/http"
	"os"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	client := slack.New(token, slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleEvents(ctx, socketClient)

	err = socketClient.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func handleEvents(ctx context.Context, client *socketmode.Client) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down socketmode listener")
			return
		case event := <-client.Events:
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
					continue
				}
				client.Ack(*event.Request)
				err := handleEventMessage(eventsAPIEvent, client)
				if err != nil {
					log.Printf("Error handling message: %v", err)
				}

			case socketmode.EventTypeSlashCommand:
				cmd, ok := event.Data.(slack.SlashCommand)
				if !ok {
					log.Printf("Could not type cast the message to a SlashCommand: %v\n", event)
					continue
				}
				client.Ack(*event.Request)
				err := handleSlashCommand(cmd, client)
				if err != nil {
					log.Printf("Error handling slash command: %v", err)
				}
			}
		}
	}
}

func handleEventMessage(event slackevents.EventsAPIEvent, client *socketmode.Client) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if ev.BotID == "" {
				// Handle only non-bot messages
				err := handleMessageEvent(ev, client)
				if err != nil {
					return fmt.Errorf("error handling message event: %w", err)
				}
			}
		}
	default:
		return fmt.Errorf("unsupported event type: %s", event.Type)
	}
	return nil
}

func handleMessageEvent(ev *slackevents.MessageEvent, client *socketmode.Client) error {
	fmt.Printf("Not implemented yet %v", ev)

	return nil
}

func handleSlashCommand(cmd slack.SlashCommand, client *socketmode.Client) error {
	switch cmd.Command {
	case "/typosweep":
		return handleTypoSweep(cmd, client)
	default:
		return fmt.Errorf("unknown command: %s", cmd.Command)
	}
}

func handleTypoSweep(cmd slack.SlashCommand, client *socketmode.Client) error {
	_, _, _, err := client.Client.JoinConversation(cmd.ChannelID)
	if err != nil {
		// If there's an error joining, it might be because we're already in the channel
		// or it's a DM. We can proceed with posting the message.
		log.Printf("Could not join channel %s: %v", cmd.ChannelID, err)
	}

	if cmd.Text == "" {
		_, _, err := client.Client.PostMessage(cmd.ChannelID, slack.MsgOptionText("Please provide some text to proofread.", false))
		if err != nil {
			return fmt.Errorf("error posting message: %w", err)
		}
		return nil
	}

	proofreaded, err := proofreadText(cmd.Text)
	if err != nil {
		return fmt.Errorf("error proofreading text: %w", err)
	}

	response := fmt.Sprintf("*Original:* %s\n%s", cmd.Text, proofreaded)
	_, _, err = client.Client.PostMessage(cmd.ChannelID, slack.MsgOptionText(response, false))
	if err != nil {
		return fmt.Errorf("error posting message: %w", err)
	}
	return nil
}
func proofreadText(text string) (string, error) {
	proofread, err := sendRequest(text)
	if err != nil {
		log.Printf("Error proofreading text: %v", err)
		return "Error in proofreading", err
	}
	return proofread, nil
}

func sendRequest(text string) (string, error) {
	requestBody := RequestBody{
		Model: "lmstudio-community/Meta-Llama-3.1-8B-Instruct-GGUF",
		Messages: []Message{
			{
				Role: "system",
				Content: `You are proofreader. Users will be asking to correct the text. Correct them with no explanations. 
Format like this:
*Typos*: list of words with a typo
*Proofread*: $whole_corrected_text`,
			},
			{
				Role:    "user",
				Content: text,
			},
		},
		Temperature: 0.7,
		MaxTokens:   -1,
		Stream:      false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:1234/v1/chat/completions", "application/json", bytes.NewBuffer(jsonData))
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
