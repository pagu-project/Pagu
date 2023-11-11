package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"pactus-bot/config"
	"pactus-bot/discord"
	"pactus-bot/wallet"
)

func main() {
	configPath := os.Args[1]
	// load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("error loading configuration %v\n", err)
		return
	}

	// load or create faucet wallet
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

	///start discord bot
	bot, err := discord.Start(cfg, w, ss)
	if err != nil {
		log.Printf("could not start discord bot: %v\n", err)
		return
	}

	// Wait here until CTRL-C or other terms signal is received.
	log.Println("Pactus Universal Robot is now running...!")
	log.Printf("The faucet address is: %v\n", cfg.FaucetAddress)
	log.Printf("The maximum faucet amount is : %.4f\n", cfg.FaucetAmount)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	_ = bot.Stop()
}
