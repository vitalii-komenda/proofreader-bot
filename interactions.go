package main

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
)

func handleInteraction(
	callback slack.InteractionCallback,
	client *slack.Client,

) error {
	action := callback.ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "approve":
		handleApprove(callback)
	case "reject":
		handleReject(callback, client)
	default:
		return fmt.Errorf("unknown action: %s", action.ActionID)
	}
	return nil
}

func handleApprove(callback slack.InteractionCallback) {

	text, ok := getProofreaded(callback.User.ID, callback.Channel.ID)
	if !ok {
		log.Printf("Failed to get proofreaded text")
		return
	}

	response := slack.MsgOptionText(text, false)
	token, err := getAccessToken(callback.User.ID)

	if err != nil {
		log.Printf("User token not found for user %s\n", callback.User.ID)
		return
	} else {

		fmt.Printf("User token found for user %s\n", callback.User.ID)
	}
	userClient := slack.New(token)

	_, _, err = userClient.PostMessage(callback.Channel.ID, response, slack.MsgOptionAsUser(true))
	if err != nil {
		log.Printf("Failed to post approval message: %v", err)
	}

}

func handleReject(callback slack.InteractionCallback, client *slack.Client) {
	_, err := client.PostEphemeral(
		callback.Channel.ID,
		callback.User.ID,
		slack.MsgOptionDeleteOriginal(callback.ResponseURL),
	)

	if err != nil {
		log.Printf("Failed to delete ephemeral message: %v", err)
	}
}
