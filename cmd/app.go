package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pactus-bot/config"
	"pactus-bot/discord"
	"pactus-bot/wallet"
	"syscall"
)

func main() {

	//load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("error loading configuration %v\n", err)
		return
	}

	//load or create faucet wallet
	w := wallet.Open(cfg)

	if w == nil {
		log.Println("faucet wallet could not be opened")
		return
	}
	// load list of validators already received faucet
	ss, err := discord.LoadData(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	//start discord bot
	bot, err := discord.Start(cfg, w, ss)
	if err != nil {
		log.Printf("could not start discord bot: %v\n", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Pactus Universal Robot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	bot.Stop()
}
