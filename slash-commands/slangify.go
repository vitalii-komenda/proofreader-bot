package slashcommands

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/vitalii-komenda/proofreader-bot/llm"
	"github.com/vitalii-komenda/proofreader-bot/repository"
)

func HandleSlangify(cmd slack.SlashCommand, client *slack.Client, db repository.AccessToken, llmInstance llm.LLM) error {
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
		_, err := client.PostEphemeral(cmd.ChannelID, cmd.UserID, slack.MsgOptionText("Please provide some text to slangify.", false))
		if err != nil {
			return fmt.Errorf("error posting message: %w", err)
		}
		return nil
	}

	slangified, err := slangifyText(cmd.Text, llmInstance)
	if err != nil {
		return fmt.Errorf("error slangidy text: %w", err)
	}

	if idx := strings.Index(slangified, "Lowkey"); idx != -1 {
		onlySlangified := slangified[idx+len("Lowkey: "):]
		cacheUserText(cmd.UserID, cmd.ChannelID, onlySlangified)
	}

	response := fmt.Sprintf("*Original:* %s\n%s", cmd.Text, slangified)

	blocks := addBlockButtons(response)

	fmt.Printf("Posting message to channel %s\n", cmd.ChannelID)
	_, err = client.PostEphemeral(cmd.ChannelID, cmd.UserID, slack.MsgOptionBlocks(blocks.
		BlockSet...))
	if err != nil {
		return fmt.Errorf("error posting message: %w", err)
	}
	return nil
}

func slangifyText(text string, llmInstance llm.LLM) (string, error) {
	proofread, err := llmInstance.SendRequest(text, llm.Slang)
	if err != nil {
		log.Printf("Error slangify text: %v", err)
		return "Error in slangifieing", err
	}
	return proofread, nil
}
