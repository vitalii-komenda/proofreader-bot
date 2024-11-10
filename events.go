package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func handleInteractions(w http.ResponseWriter, r *http.Request, client *slack.Client) {
	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(r.FormValue("payload")), &payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch payload.Type {
	case slack.InteractionTypeBlockActions:
		err := handleInteraction(payload, client)
		if err != nil {
			log.Printf("Error handling shortcut: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleSlashCommand(w http.ResponseWriter, r *http.Request, client *slack.Client) {
	config := getConfig()
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(config.SlackVerificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/doublecheck":
		w.WriteHeader(http.StatusOK)
		go func() {
			err := handleDoublecheck(s, client)
			if err != nil {
				log.Printf("Error handling slash command: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
	case "/slangify":
		w.WriteHeader(http.StatusOK)
		go func() {
			err := handleSlangify(s, client)
			if err != nil {
				log.Printf("Error handling slash command: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func addBlockButtons(response string) slack.Blocks {
	return slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", response, false, false),
				nil,
				nil,
			),
			slack.NewActionBlock(
				"",
				slack.NewButtonBlockElement("approve", "approve", slack.NewTextBlockObject("plain_text", "Send", false, false)),
				slack.NewButtonBlockElement("reject", "reject", slack.NewTextBlockObject("plain_text", "Delete", false, false)),
			),
		},
	}
}
