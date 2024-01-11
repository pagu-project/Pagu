package discord

import (
	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
)

type BotEngine struct {
	Wallet wallet.IWallet
	Store  store.Store
	Cm     *client.Mgr
	Cfg    *config.Config
	logger *log.SubLogger
}

func Start(logger *log.SubLogger, cfg *config.Config, w wallet.IWallet) (*BotEngine, error) {
	cm := client.NewClientMgr()

	for _, rn := range cfg.RPCNodes {
		c, err := client.NewClient(rn)
		if err != nil {
			logger.Error("can't make new client.", "RPC Node address", rn)
			continue
		}
		logger.Info("connecting to RPC Node", "addr", rn)
		cm.AddClient(rn, c)
	}

	return &BotEngine{
		logger: logger,
		Wallet: w,
		Cfg:    cfg,
		Cm:     cm,
	}, nil
}

func (be *BotEngine) Stop() {
	be.logger.Info("shutting bot engine down...")

	be.Cm.Close()
}
