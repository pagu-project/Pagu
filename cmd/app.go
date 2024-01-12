package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
)

func main() {
	configPath := os.Args[1]

	// load configuration.
	config, err := config.Load(configPath)
	if err != nil {
		log.Panic("error loading configuration %v\n", err)
	}

	// load or create wallet.
	wallet := wallet.Open(config)
	if wallet == nil {
		log.Panic("wallet could not be opened, wallet is nil", "path", config.WalletPath)
	}

	log.Info("wallet opened successfully", "address", wallet.Address())

	// new subLogger for engine.
	sl := log.NewSubLogger("engine")

	// load store
	store, err := store.LoadStore(config)
	if err != nil {
		log.Panic("could not load store", "err", err, "path", config.StorePath)
	}

	// start botEngine engine.
	botEngine, err := engine.Start(sl, config, wallet, store)
	if err != nil {
		log.Panic("could not start discord bot", "err", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan

	// gracefully shutdown the bot.
	botEngine.Stop()
}
