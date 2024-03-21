package engine

import (
	"context"
	"errors"

	"github.com/robopac-project/RoboPac/client"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/database"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/engine/command/blockchain"
	"github.com/robopac-project/RoboPac/engine/command/network"
	"github.com/robopac-project/RoboPac/log"
	"github.com/robopac-project/RoboPac/wallet"
)

type BotEngine struct {
	ctx    context.Context //nolint
	cancel context.CancelFunc

	clientMgr *client.Mgr
	rootCmd   command.Command

	blockchainCmd blockchain.Blockchain
	networkCmd    network.Network
}

func NewBotEngine(cfg *config.Config) (*BotEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cm := client.NewClientMgr(ctx)

	localClient, err := client.NewClient(cfg.LocalNode)
	if err != nil {
		cancel()
		return nil, err
	}

	cm.AddClient(localClient)

	for _, nn := range cfg.NetworkNodes {
		c, err := client.NewClient(nn)
		if err != nil {
			log.Error("can't add new network node client", "err", err, "addr", nn)
		}
		cm.AddClient(c)
	}
	cm.Start()

	var wal wallet.IWallet

	if cfg.WalletConfig.Enable {
		// load or create wallet.
		wal = wallet.Open(&cfg.WalletConfig)
		if wal == nil {
			cancel()
			return nil, errors.New("can't open the wallet")
		}

		log.Info("wallet opened successfully", "address", wal.Address())
	}

	// load database
	db, err := database.NewDB(cfg.DataBasePath)
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("database loaded successfully")

	return newBotEngine(cm, wal, db, cfg.AuthIDs, ctx, cancel), nil
}

func newBotEngine(cm *client.Mgr, _ wallet.IWallet, _ *database.DB, _ []string,
	ctx context.Context, cnl context.CancelFunc,
) *BotEngine {
	rootCmd := command.Command{
		Emoji:       "🤖",
		Name:        "robopac",
		Desc:        "RoboPAC",
		Help:        "RoboPAC Help",
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		SubCommands: []command.Command{},
	}

	netCmd := network.NewNetwork(ctx, cm)
	bcCmd := blockchain.NewBlockchain(ctx, cm)

	return &BotEngine{
		ctx:           ctx,
		cancel:        cnl,
		clientMgr:     cm,
		rootCmd:       rootCmd,
		networkCmd:    netCmd,
		blockchainCmd: bcCmd,
	}
}

func (be *BotEngine) Commands() []command.Command {
	return be.rootCmd.SubCommands
}

func (be *BotEngine) RegisterAllCommands() {
	be.rootCmd.AddSubCommand(be.blockchainCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.networkCmd.GetCommand())

	be.rootCmd.AddHelpSubCommand()
}

func (be *BotEngine) Run(appID command.AppID, callerID string, tokens []string) command.CommandResult {
	log.Debug("run command", "callerID", callerID, "inputs", tokens)

	cmd, argsIndex := be.getCommand(tokens)
	if !cmd.HasAppId(appID) {
		return cmd.FailedResult("unauthorized appID: %v", appID)
	}

	if cmd.Handler == nil {
		return cmd.HelpResult()
	}

	args := tokens[argsIndex:]
	err := cmd.CheckArgs(args)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.Handler(cmd, appID, callerID, args...)
}

func (be *BotEngine) getCommand(tokens []string) (command.Command, int) {
	index := 0
	targetCmd := be.rootCmd
	cmds := be.rootCmd.SubCommands
	for {
		if len(tokens) <= index {
			break
		}
		token := tokens[index]
		index++

		found := false
		for _, cmd := range cmds {
			if cmd.Name == token {
				targetCmd = cmd
				cmds = cmd.SubCommands
				found = true

				break
			}
		}

		if !found {
			break
		}
	}

	return targetCmd, index
}

func (be *BotEngine) NetworkStatus() (*network.NetStatus, error) {
	netInfo, err := be.clientMgr.GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	chainInfo, err := be.clientMgr.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}

	cs, err := be.clientMgr.GetCirculatingSupply()
	if err != nil {
		cs = 0
	}

	return &network.NetStatus{
		ConnectedPeersCount: netInfo.ConnectedPeersCount,
		ValidatorsCount:     chainInfo.TotalValidators,
		TotalBytesSent:      netInfo.TotalSentBytes,
		TotalBytesReceived:  netInfo.TotalReceivedBytes,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   chainInfo.TotalPower,
		TotalCommitteePower: chainInfo.CommitteePower,
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   cs,
	}, nil
}

func (be *BotEngine) Stop() {
	log.Info("shutting bot engine down...")

	be.cancel()
	be.clientMgr.Stop()
}

func (be *BotEngine) Start() {
	log.Info("starting the bot engine...")
}
