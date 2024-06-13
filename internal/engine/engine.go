package engine

import (
	"context"
	"time"

	"github.com/pagu-project/Pagu/internal/engine/command/market"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/job"
	"github.com/pagu-project/Pagu/pkg/cache"

	"github.com/pagu-project/Pagu/internal/repository"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/engine/command/calculator"
	"github.com/pagu-project/Pagu/internal/engine/command/network"
	phoenixtestnet "github.com/pagu-project/Pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/Pagu/internal/engine/command/zealy"
	client2 "github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/log"
	"github.com/pagu-project/Pagu/pkg/wallet"

	"github.com/pagu-project/Pagu/config"
)

type BotEngine struct {
	ctx    context.Context //nolint
	cancel context.CancelFunc

	clientMgr        *client2.Mgr
	phoenixClientMgr *client2.Mgr
	rootCmd          command.Command

	blockchainCmd calculator.Calculator
	networkCmd    network.Network
	phoenixCmd    phoenixtestnet.Phoenix
	zealyCmd      zealy.Zealy
	marketCmd     market.Market
}

func NewBotEngine(cfg *config.Config) (*BotEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// ? adding main network client manager.
	cm := client2.NewClientMgr(ctx)

	localClient, err := client2.NewClient(cfg.LocalNode)
	if err != nil {
		cancel()
		return nil, err
	}

	cm.AddClient(localClient)

	for _, nn := range cfg.NetworkNodes {
		c, err := client2.NewClient(nn)
		if err != nil {
			log.Error("can't add new network node client", "err", err, "addr", nn)
		}
		cm.AddClient(c)
	}

	// ? adding phoenix test network client manager.
	phoenixCm := client2.NewClientMgr(ctx)
	for _, tnn := range cfg.Phoenix.NetworkNodes {
		c, err := client2.NewClient(tnn)
		if err != nil {
			log.Error("can't add new network node client", "err", err, "addr", tnn)
		}

		phoenixCm.AddClient(c)
	}

	// ? opening wallet if it's enabled.
	var wal *wallet.Wallet
	if cfg.Wallet.Enable {
		// load or create wallet.
		wal = wallet.Open(&cfg.Wallet)
		if wal == nil {
			cancel()
			return nil, WalletError{
				Reason: "can't open mainnet wallet",
			}
		}

		log.Info("wallet opened successfully", "address", wal.Address())
	}

	// ? opening testnet (Phoenix) wallet if it's enabled.
	var phoenixWal *wallet.Wallet
	if cfg.TestNetWallet.Enable {
		// load or create wallet.
		wal = wallet.Open(&cfg.TestNetWallet)
		if wal == nil {
			cancel()
			return nil, WalletError{
				Reason: "can't open testnet wallet",
			}
		}

		log.Info("testnet wallet opened successfully", "address", wal.Address())
	}

	// ? loading database.
	db, err := repository.NewDB(cfg.Database.URL)
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("database loaded successfully")

	return newBotEngine(cm, phoenixCm, wal, phoenixWal, db, cfg.AuthIDs, ctx, cancel), nil
}

func newBotEngine(cm, ptcm *client2.Mgr, wallet *wallet.Wallet, phoenixWal *wallet.Wallet, db *repository.DB, _ []string,
	ctx context.Context, cnl context.CancelFunc,
) *BotEngine {
	rootCmd := command.Command{
		Emoji:       "🤖",
		Name:        "pagu",
		Desc:        "Root Command",
		Help:        "Pagu Help Command",
		AppIDs:      command.AllAppIDs(),
		SubCommands: make([]command.Command, 3),
	}

	// price caching job
	priceCache := cache.NewBasic[string, entity.Price](0 * time.Second)
	priceJob := job.NewPrice(priceCache)
	priceJobSched := job.NewScheduler()
	priceJobSched.Submit(priceJob)
	go priceJobSched.Run()

	netCmd := network.NewNetwork(ctx, cm)
	bcCmd := calculator.NewCalculator(cm)
	ptCmd := phoenixtestnet.NewPhoenix(phoenixWal, ptcm, *db)
	zCmd := zealy.NewZealy(db, wallet)
	marketCmd := market.NewMarket(cm, priceCache)

	return &BotEngine{
		ctx:              ctx,
		cancel:           cnl,
		clientMgr:        cm,
		rootCmd:          rootCmd,
		networkCmd:       netCmd,
		blockchainCmd:    bcCmd,
		phoenixCmd:       ptCmd,
		phoenixClientMgr: ptcm,
		zealyCmd:         zCmd,
		marketCmd:        marketCmd,
	}
}

func (be *BotEngine) Commands() []command.Command {
	return be.rootCmd.SubCommands
}

func (be *BotEngine) RegisterAllCommands() {
	be.rootCmd.AddSubCommand(be.blockchainCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.networkCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.zealyCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.marketCmd.GetCommand())
	// be.rootCmd.AddSubCommand(be.phoenixCmd.GetCommand()) // TODO: FIX WALLET ISSUE

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

	if len(targetCmd.Args) != 0 && index != 0 {
		return targetCmd, index - 1 //! TODO: FIX ME IN THE MAIN LOGIC
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
	log.Info("Stopping the Bot Engine")

	be.cancel()
	be.clientMgr.Stop()
	be.phoenixClientMgr.Stop()
}

func (be *BotEngine) Start() {
	log.Info("Starting the Bot Engine")

	be.clientMgr.Start()
	be.phoenixClientMgr.Start()
}
