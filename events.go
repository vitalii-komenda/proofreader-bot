package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/vitalii-komenda/proofreader-bot/slash-commands"
)

func handleInteractions(w http.ResponseWriter, r *http.Request, client *slack.Client) {
	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(r.FormValue("payload")), &payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch payload.Type {
	case slack.InteractionTypeBlockActions:
		w.WriteHeader(http.StatusOK)
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
			err := slashcommands.HandleDoublecheck(s, client, db, llmInstance)
			if err != nil {
				log.Printf("Error handling slash command: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
	case "/slangify":
		w.WriteHeader(http.StatusOK)
		go func() {
			err := slashcommands.HandleSlangify(s, client, db, llmInstance)
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
