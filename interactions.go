package main

import (
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/vitalii-komenda/proofreader-bot/llm"
	slashcommands "github.com/vitalii-komenda/proofreader-bot/slash-commands"
)

func handleInteraction(
	callback slack.InteractionCallback,
	client *slack.Client,
) error {
	action := callback.ActionCallback.BlockActions[0]
	switch action.ActionID {
	case "approve", "approve-slangify":
		handleApprove(callback, action.ActionID)
	case "reject":
		handleReject(callback, client)
	case "rephrase":
		handleRephrase(callback, client)
	default:
		return fmt.Errorf("unknown action: %s", action.ActionID)
	}
	return nil
}

func rephraseText(text string, llmInstance llm.LLM) (string, error) {
	proofread, err := llmInstance.SendRequest(text, llm.Rephrase)
	if err != nil {
		log.Printf("Error rephrasing text: %v", err)
		return "Error in rephrasing", err
	}
	return proofread, nil
}

func handleRephrase(callback slack.InteractionCallback, client *slack.Client) {
	text, ok := slashcommands.GetUserText(callback.User.ID, callback.Channel.ID, string(llm.Slang)+"original")
	if !ok {
		log.Printf("Failed to get slangified text")
		return
	}

	rephrased, err := rephraseText(text, llmInstance)
	if err != nil {
		log.Printf("Failed to rephrase text: %v", err)
		return
	}

	slashcommands.CacheUserText(callback.User.ID, callback.Channel.ID, string(llm.Slang), slashcommands.SeparateProposed(rephrased))

	token, err := db.GetAccessToken(callback.User.ID)
	if err != nil {
		log.Printf("User token not found for user %s\n", callback.User.ID)
		return
	}

	userClient := slack.New(token)

	response := fmt.Sprintf("*Original:* %s\n%s", text, rephrased)

	blocks := slashcommands.AddSendDelRephraseButtons(response)

	_, err = userClient.PostEphemeral(callback.Channel.ID, callback.User.ID, slack.MsgOptionBlocks(blocks.
		BlockSet...), slack.MsgOptionReplaceOriginal(callback.ResponseURL))
	if err != nil {
		log.Printf("error posting message: %w", err)
	}
}

func handleApprove(callback slack.InteractionCallback, actionID string) {
	var text string
	var ok bool

	if actionID == "approve-slangify" {
		text, ok = slashcommands.GetUserText(callback.User.ID, callback.Channel.ID, string(llm.Slang))
	} else {
		text, ok = slashcommands.GetUserText(callback.User.ID, callback.Channel.ID, string(llm.Proofread))
	}
	if !ok {
		log.Printf("Failed to get text %s", actionID)
		return
	}

	response := slack.MsgOptionText(text, false)
	token, err := db.GetAccessToken(callback.User.ID)

	if err != nil {
		log.Printf("User token not found for user %s\n", callback.User.ID)
		return
	} else {

		fmt.Printf("User token found for user %s\n", callback.User.ID)
	}
	userClient := slack.New(token)

	_, err = userClient.PostEphemeral(
		callback.Channel.ID,
		callback.User.ID,
		slack.MsgOptionDeleteOriginal(callback.ResponseURL),
	)

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
