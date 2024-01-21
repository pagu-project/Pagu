package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
)

func main() {
	// load configuration.
	config, err := config.Load()
	if err != nil {
		log.Panic("error loading configuration %v\n", err)
	}

	// starting client manager for RPC.
	cm := client.NewClientMgr()

	for _, rn := range config.RPCNodes {
		c, err := client.NewClient(rn)
		if err != nil {
			log.Error("can't make new client.", "RPC Node address", rn)
			continue
		}
		log.Info("connecting to RPC Node", "addr", rn)
		cm.AddClient(rn, c)
	}

	// initializing logger global instance.
	log.InitGlobalLogger()

	// new subLogger for engine.
	eSl := log.NewSubLogger("engine")

	// new subLogger for store.
	sSl := log.NewSubLogger("store")

	// new subLogger for store.
	wSl := log.NewSubLogger("wallet")

	// load or create wallet.
	wallet := wallet.Open(config, wSl)
	if wallet == nil {
		log.Panic("wallet could not be opened, wallet is nil", "path", config.WalletPath)
	}

	log.Info("wallet opened successfully", "address", wallet.Address())

	// load store.
	store, err := store.LoadStore(config, sSl)
	if err != nil {
		log.Panic("could not load store", "err", err, "path", config.StorePath)
	}

	log.Info("store loaded successfully", "path", config.StorePath)

	// starting botEngine.
	botEngine, err := engine.NewBotEngine(eSl, cm, wallet, store)
	if err != nil {
		log.Panic("could not start discord bot", "err", err)
	}
	botEngine.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan

	// gracefully shutdown the bot.
	botEngine.Stop()
	cm.Close()
}
