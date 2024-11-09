package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/vitalii-komenda/proofreader-bot/llm"
	"github.com/vitalii-komenda/proofreader-bot/repository"
)

var db repository.AccessToken
var llmInstance llm.LLM

func main() {
	config := getConfig()

	var model llm.LLM
	if config.Environment == "local" {
		model = &llm.LLama{}
		db = repository.InitSQLite(config.EncryptionKey)
		fmt.Println("[INFO] Running in development environment")
	} else {
		model = &llm.OpenAI{Token: config.OPENAI_API_KEY}
		db = repository.Init(config.EncryptionKey, config.GCP_PROJECT_ID)
		fmt.Println("[INFO] Running in production environment")
	}
	llmInstance = llm.Init(model)

	client := slack.New(config.SlackAppToken)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		fmt.Println("[INFO] Server listening :" + config.PORT)
		log.Fatal(http.ListenAndServe(":"+config.PORT, nil))
	}()

	<-ctx.Done()
}
