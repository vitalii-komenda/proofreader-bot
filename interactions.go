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
		handleApprove(callback, client)
	case "reject":
		handleReject(callback, client)
	default:
		return fmt.Errorf("unknown action: %s", action.ActionID)
	}
	return nil
}

func handleApprove(callback slack.InteractionCallback, client *slack.Client) {
	response := slack.MsgOptionText("copied text", false)
	token, ok := userTokens.Load(callback.User.ID)
	if !ok {
		log.Printf("User token not found for user %s", callback.User.ID)
		return
	} else {

		fmt.Printf("User token found for user %s - %s", callback.User.ID, token)
	}
	client2 := slack.New(token.(string))

	_, _, err := client2.PostMessage(callback.Channel.ID, response, slack.MsgOptionAsUser(true))
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
