package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

var envLoaded = false

type Config struct {
	SlackBotToken          string
	SlackAppToken          string
	SlackClientID          string
	SlackClientSecret      string
	SlackRedirectURL       string
	SlackUserOAuthToken    string
	SlackSigningSecret     string
	SlackVerificationToken string
}

func initEnv() {
	if envLoaded {
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	envLoaded = true
}
func getEnv(key string) string {
	initEnv()

	value := os.Getenv(key)
	if len(value) == 0 {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
	}
	return value
}

func parseConfig() Config {
	return Config{
		SlackBotToken:          getEnv("SLACK_BOT_TOKEN"),
		SlackAppToken:          getEnv("SLACK_APP_TOKEN"),
		SlackClientID:          getEnv("SLACK_CLIENT_ID"),
		SlackClientSecret:      getEnv("SLACK_CLIENT_SECRET"),
		SlackRedirectURL:       getEnv("SLACK_REDIRECT_URL"),
		SlackUserOAuthToken:    getEnv("SLACK_USER_OAUTH_TOKEN"),
		SlackSigningSecret:     getEnv("SLACK_SIGNING_SECRET"),
		SlackVerificationToken: getEnv("SLACK_VERIFICATION_TOKEN"),
	}
}

func main() {
	config := parseConfig()
	client := slack.New(config.SlackUserOAuthToken, slack.OptionAppLevelToken(config.SlackAppToken))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	http.HandleFunc("/oauth/start", startOAuth)
	http.HandleFunc("/oauth/callback", handleOAuthCallback)
	// http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
	// 	handleEvents(w, r, client)
	// })

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
