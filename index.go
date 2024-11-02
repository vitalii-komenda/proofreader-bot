package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"goslackbot/llm"
)

var db *sql.DB
var llmInstance llm.LLM

func main() {
	config := getConfig()
	db = initDB(config.ENCRYPTION_KEY)
	defer db.Close()

	client := slack.New(config.SlackUserOAuthToken, slack.OptionAppLevelToken(config.SlackAppToken))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	llmInstance = &llm.LLama{
		URL:   "http://localhost:1234/v1/chat/completions",
		Model: "lmstudio-community/Meta-Llama-3.1-8B-Instruct-GGUF",
		Messages: []llm.Message{
			{
				Role: "system",
				Content: `You are proofreader. Users will be asking to correct the text. Correct them with no explanations. 
Format like this:
*Typos*: list of words with a typo
*Proofread*: $whole_corrected_text`,
			},
		},
		Temperature: 0.7,
		MaxTokens:   -1,
		Stream:      false,
	}

	http.HandleFunc("/oauth/start", startOAuth)
	http.HandleFunc("/oauth/callback", handleOAuthCallback)
	http.HandleFunc("/slack/slash-commands", func(w http.ResponseWriter, r *http.Request) {
		handleSlashCommand(w, r, client)
	})
	http.HandleFunc("/slack/interactions", func(w http.ResponseWriter, r *http.Request) {
		handleInteractions(w, r, client)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	go func() {
		fmt.Println("[INFO] Server listening :3000")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}()

	<-ctx.Done()
}
