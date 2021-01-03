package discord

import (
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

// TODO: Build some sort of command parser with a configurable prefix,
//       so we can easily delegate command messages to some kind of command handler

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Detect if we're dealing with a command
	isCommand := m.Content[0] == '!'

	// Detect if the message is from an owner
	isOwner := m.Author.ID == os.Getenv("HACKMUD_OWNER_ID")

	// Handle commands
	if isCommand {
		// Don't run any commands if not the owner
		if !isOwner {
			log.Warn("User '" + m.Author.Username + "' is not allowed to run '" + m.Content + "' command")
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, request denied")
			return
		}

		log.Info("User '" + m.Author.Username + "' ran '" + m.Content + "' command")

		// TODO: Refactor this so we can easily detect unknown/unavailable commands
		// TODO: Actually measure some kind of latency with this (latency from Discord API?)
		// Handle ping/pong commands
		if m.Content == "!ping" {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> ***PONG!***")
		} else if m.Content == "!pong" {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> ***PING!***")
			// Handle inspire command
		} else if m.Content == "!inspire" {
			resp, err := http.Get("https://inspirobot.me/api?generate=true")
			if err != nil {
				log.Error(err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error(err)
				return
			}
			s.ChannelMessageSend(m.ChannelID, string(body))
		} else {
			// TODO: Use the following regex to validate commands (a-z and lowercase only, otherwise invalid)
			//       ^!([a-z]+)$
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> _Unknown command: "+m.Content+"_")
		}
	} else {
		// TODO: Handle two-way chat bridge logic here
		log.Info("[CHAT] " + m.Author.Username + ": " + m.Content)
	}
}
