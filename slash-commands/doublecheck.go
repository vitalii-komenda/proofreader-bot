package slashcommands

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/vitalii-komenda/proofreader-bot/llm"
	"github.com/vitalii-komenda/proofreader-bot/repository"
)

func HandleDoublecheck(cmd slack.SlashCommand, client *slack.Client, db repository.AccessToken, llmInstance llm.LLM) error {
	token, err := db.GetAccessToken(cmd.UserID)
	if err != nil {
		log.Printf("User token not found for user %s\n", cmd.UserID)
		return nil
	}

	client = slack.New(token)

	_, _, _, err = client.JoinConversation(cmd.ChannelID)
	if err != nil {
		// If there's an error joining, it might be because we're already in the channel
		// or it's a DM. We can proceed with posting the message.
		// log.Printf("Could not join channel %s: %v", cmd.ChannelID, err)
	}

	if cmd.Text == "" {
		_, err := client.PostEphemeral(cmd.ChannelID, cmd.UserID, slack.MsgOptionText("Please provide some text to proofread.", false))
		if err != nil {
			return fmt.Errorf("error posting message: %w", err)
		}
		return nil
	}

	proofreaded, err := proofreadText(cmd.Text, llmInstance)
	if err != nil {
		return fmt.Errorf("error proofreading text: %w", err)
	}

	if idx := strings.Index(proofreaded, "Proofread"); idx != -1 {
		onlyProofreaded := proofreaded[idx+len("Proofread: "):]
		cacheUserText(cmd.UserID, cmd.ChannelID, onlyProofreaded)
	}

	response := fmt.Sprintf("*Original:* %s\n%s", cmd.Text, proofreaded)

	blocks := addBlockButtons(response)

	fmt.Printf("Posting message to channel %s\n", cmd.ChannelID)
	_, err = client.PostEphemeral(cmd.ChannelID, cmd.UserID, slack.MsgOptionBlocks(blocks.
		BlockSet...))
	if err != nil {
		return fmt.Errorf("error posting message: %w", err)
	}
	return nil
}

func proofreadText(text string, llmInstance llm.LLM) (string, error) {
	proofread, err := llmInstance.SendRequest(text, llm.Proofreader)
	if err != nil {
		log.Printf("Error proofreading text: %v", err)
		return "Error in proofreading", err
	}
	return proofread, nil
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
