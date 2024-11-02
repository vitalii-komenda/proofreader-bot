package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleDoublecheck(cmd slack.SlashCommand, client *slack.Client) error {
	_, _, _, err := client.JoinConversation(cmd.ChannelID)
	if err != nil {
		// If there's an error joining, it might be because we're already in the channel
		// or it's a DM. We can proceed with posting the message.
		log.Printf("Could not join channel %s: %v", cmd.ChannelID, err)
	}

	if cmd.Text == "" {
		_, err := client.PostEphemeral(cmd.ChannelID, cmd.UserID, slack.MsgOptionText("Please provide some text to proofread.", false))
		if err != nil {
			return fmt.Errorf("error posting message: %w", err)
		}
		return nil
	}

	proofreaded, err := proofreadText(cmd.Text)
	if err != nil {
		return fmt.Errorf("error proofreading text: %w", err)
	}

	if idx := strings.Index(proofreaded, "Proofread"); idx != -1 {
		onlyProofreaded := proofreaded[idx+len("Proofread: "):]
		addProofreaded(cmd.UserID, cmd.ChannelID, onlyProofreaded)
	}

	response := fmt.Sprintf("*Original:* %s\n%s", cmd.Text, proofreaded)

	blocks := addBlockButtons(response)

	_, err = client.PostEphemeral(cmd.ChannelID, cmd.UserID, slack.MsgOptionBlocks(blocks.
		BlockSet...))
	if err != nil {
		return fmt.Errorf("error posting message: %w", err)
	}
	return nil
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
				slack.NewButtonBlockElement("approve", "approve", slack.NewTextBlockObject("plain_text", "Approve", false, false)),
				slack.NewButtonBlockElement("reject", "reject", slack.NewTextBlockObject("plain_text", "Reject", false, false)),
			),
		},
	}

}

func proofreadText(text string) (string, error) {
	proofread, err := llmInstance.SendRequest(text)
	if err != nil {
		log.Printf("Error proofreading text: %v", err)
		return "Error in proofreading", err
	}
	return proofread, nil
}
