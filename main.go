package main

import (
	"log"
	"os"

	"github.com/Dids/hackchat/hackmud"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create a new API client
	client := hackmud.NewClient(os.Getenv("HACKMUD_CHAT_API_TOKEN"))

	// Enable request debugging
	client.Debug = true

	// TODO: Implement
	// _ = discord.NewClient(os.Getenv("DISCORD_API_TOKEN"))

	// TODO: Implement all other API endpoints as well
	// Attempt to get the account data
	account, err := client.GetAccountData()
	if err != nil {
		panic(err)
	}
	log.Println("Got account data:", account)
}
