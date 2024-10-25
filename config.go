package main

type Config struct {
	SlackBotToken          string
	SlackAppToken          string
	SlackClientID          string
	SlackClientSecret      string
	SlackRedirectURL       string
	SlackUserOAuthToken    string
	SlackSigningSecret     string
	SlackVerificationToken string
	ENCRYPTION_KEY         string
}

func getConfig() Config {
	return Config{
		SlackBotToken:          getEnv("SLACK_BOT_TOKEN"),
		SlackAppToken:          getEnv("SLACK_APP_TOKEN"),
		SlackClientID:          getEnv("SLACK_CLIENT_ID"),
		SlackClientSecret:      getEnv("SLACK_CLIENT_SECRET"),
		SlackRedirectURL:       getEnv("SLACK_REDIRECT_URL"),
		SlackUserOAuthToken:    getEnv("SLACK_USER_OAUTH_TOKEN"),
		SlackSigningSecret:     getEnv("SLACK_SIGNING_SECRET"),
		SlackVerificationToken: getEnv("SLACK_VERIFICATION_TOKEN"),
		ENCRYPTION_KEY:         getEnv("ENCRYPTION_KEY"),
	}
}
