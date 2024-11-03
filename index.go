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

	client := slack.New(config.SlackAppToken)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// model := llm.LLama{}
	model := llm.OpenAI{Token: config.OPENAI_API_KEY}
	llmInstance = llm.Init(&model)

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
		fmt.Println("[INFO] Server listening :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-ctx.Done()
}
