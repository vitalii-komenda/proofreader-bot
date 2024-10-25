package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

var envLoaded = false
var db *sql.DB

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
func main() {
	config := getConfig()

	db = initDB(config.ENCRYPTION_KEY)
	defer db.Close()

	client := slack.New(config.SlackUserOAuthToken, slack.OptionAppLevelToken(config.SlackAppToken))
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
		fmt.Println("[INFO] Server listening :3000")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}()

	<-ctx.Done()
}
