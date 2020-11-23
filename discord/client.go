package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Client wrapper for Discord
type Client struct {
	apiToken      string
	DiscordClient *discordgo.Session
	Debug         bool
}

// NewClient returns a new Client wrapper for Discord
func NewClient(apiToken string) *Client {
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", apiToken))
	if err != nil {
		panic(err)
	}

	// TODO: Setup error handlers and stuff

	return &Client{
		apiToken:      apiToken,
		DiscordClient: discord,
		Debug:         false,
	}
}
