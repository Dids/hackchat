package main

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/Dids/hackchat/discord"
	"github.com/Dids/hackchat/hackmud"
	_ "github.com/joho/godotenv/autoload"
)

var (
	// Internal version of the application
	version string = "dev"

	// Hackmud API client
	Hackmud *hackmud.Client

	// Discord client
	Discord *discord.Client
)

func main() {
	log.Println("--- hackchat "+version, "---")

	log.Println("Starting up")
	// TODO: Does logrus close the file or do we need to close it ourselves?
	// Log to both console and file
	logFile, logFileErr := os.OpenFile("hackchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logFileErr != nil {
		panic(logFileErr)
	}
	defer logFile.Close()
	logWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logWriter)

	// Run shutdown logic no matter what
	defer shutdown()

	// Create a new API client
	Hackmud = hackmud.NewClient(os.Getenv("HACKMUD_CHAT_API_TOKEN"))

	// Enable request debugging
	Hackmud.Debug = true

	Discord = discord.NewClient(os.Getenv("DISCORD_API_TOKEN"))

	log.Info("Starting Discord client..")
	if err := Discord.Start(); err != nil {
		panic(err)
	}

	// FIXME: We need to have an event handler/delegate
	//        which controls the data flow between hackmud
	//        and Discord, otherwise this won't really work!

	// TODO: Implement all other API endpoints as well
	// Attempt to get the account data
	// accountData, err := Hackmud.GetAccountData()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("Got account data:", account)

	// TODO: This works, but should be formatted better etc.
	// if err := Discord.Send("```json\n" + accountData.ToJSON() + "\n```"); err != nil {
	// 	log.Error(err)
	// }

	// Wait for CTRL-C
	log.Info("Hackchat bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func shutdown() {
	log.Info("Preparing to shutdown..")

	// Stop Discord client
	log.Info("Stopping Discord client..")
	if err := Discord.Stop(); err != nil {
		panic(err)
	}

	log.Info("Exiting..")
}
