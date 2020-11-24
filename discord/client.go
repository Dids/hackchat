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
	discordClient, err := discordgo.New(fmt.Sprintf("Bot %s", apiToken))
	if err != nil {
		panic(err)
	}

	// TODO: Setup error handlers and more handlers in general?

	// Register the messageCreate func as a callback for MessageCreate events.
	discordClient.AddHandler(messageCreate)

	// TODO: We definitely want more intents, but do we really need to specify them manually now?!
	// In this example, we only care about receiving message events.
	discordClient.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	return &Client{
		apiToken:      apiToken,
		DiscordClient: discordClient,
		Debug:         false,
	}
}

func (c *Client) Start() error {
	return c.DiscordClient.Open()
}

func (c *Client) Stop() error {
	return c.DiscordClient.Close()
}
