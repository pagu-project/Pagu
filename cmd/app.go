package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/wallet"
)

func main() {
	configPath := os.Args[1]
	// load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Panic("error loading configuration %v\n", err)
		return
	}

	// load or create faucet wallet
	w := wallet.Open(cfg)

	if w == nil {
		log.Panic("faucet wallet could not be opened")
		return
	}

	sl := log.NewSubLogger("engine")

	// start be engine
	be, err := discord.Start(sl, cfg, w)
	if err != nil {
		log.Panic("could not start discord bot: %v\n", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	be.Stop()
}
