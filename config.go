package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var envLoaded = false

func initEnv() {
	if envLoaded {
		return
	}
	var err error

	if len(os.Args) > 1 {
		err = godotenv.Load(os.Args[1])
	}

	err = godotenv.Load()
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

type Config struct {
	SlackBotToken          string
	SlackAppToken          string
	SlackClientID          string
	SlackClientSecret      string
	SlackRedirectURL       string
	SlackUserOAuthToken    string
	SlackSigningSecret     string
	SlackVerificationToken string
	EncryptionKey          string
	Environment            string
	OPENAI_API_KEY         string
	GCP_PROJECT_ID         string
	PORT                   string
}

func getConfig() Config {
	config := Config{
		SlackBotToken:          getEnv("SLACK_BOT_TOKEN"),
		SlackAppToken:          getEnv("SLACK_APP_TOKEN"),
		SlackClientID:          getEnv("SLACK_CLIENT_ID"),
		SlackClientSecret:      getEnv("SLACK_CLIENT_SECRET"),
		SlackRedirectURL:       getEnv("SLACK_REDIRECT_URL"),
		SlackUserOAuthToken:    getEnv("SLACK_USER_OAUTH_TOKEN"),
		SlackSigningSecret:     getEnv("SLACK_SIGNING_SECRET"),
		SlackVerificationToken: getEnv("SLACK_VERIFICATION_TOKEN"),
		EncryptionKey:          getEnv("ENCRYPTION_KEY"),
		Environment:            getEnv("ENVIRONMENT"),
		PORT:                   getEnv("PORT"),
	}

	// Type switch to handle specific configurations
	switch v := config.Environment; v {
	case "gcp":
		config.OPENAI_API_KEY = getEnv("OPENAI_API_KEY")
		config.GCP_PROJECT_ID = getEnv("GCP_PROJECT_ID")
		return config
	case "local":
		return config
	default:
		panic("Unsupported config type")
	}
}
