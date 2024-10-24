package main

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func handleInteraction(
	callback slack.InteractionCallback,
	client *socketmode.Client,
	request socketmode.Request,
) error {
	action := callback.ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "approve":

		handleApprove(callback, client, request)
	case "reject":
		handleReject(callback, client, request)
	default:
		return fmt.Errorf("unknown action: %s", action.ActionID)
	}
	return nil
}

func handleApprove(callback slack.InteractionCallback, client *socketmode.Client, request socketmode.Request) {
	client.Ack(request)

	response := slack.MsgOptionText("You approved the message!", false)
	_, _, err := client.Client.PostMessage(callback.Channel.ID, response)
	if err != nil {
		log.Printf("Failed to post approval message: %v", err)
	}
}

func handleReject(callback slack.InteractionCallback, client *socketmode.Client, request socketmode.Request) {
	client.Ack(request)

	log.Printf("Channel ID: %s, User ID: %s, Response URL: %s", callback.Channel.ID, callback.User.ID, callback.ResponseURL)

	_, err := client.PostEphemeral(
		callback.Channel.ID,
		callback.User.ID,
		slack.MsgOptionDeleteOriginal(callback.ResponseURL),
	)

	if err != nil {
		log.Printf("Failed to delete ephemeral message: %v", err)
	}

}
