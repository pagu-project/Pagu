package engine

import (
	"context"
	"errors"
	"time"

	"github.com/pagu-project/Pagu/internal/engine/command/voucher"

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
	voucherCmd    voucher.Voucher
	marketCmd     market.Market
}

type IEngine interface {
	Run(appID entity.AppID, callerID string, tokens []string) (*command.CommandResult, error)
	Commands() []command.Command
}

func NewBotEngine(cfg *config.Config) (*BotEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	db, err := repository.NewDB(cfg.Database.URL)
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("database loaded successfully")

	cm := client2.NewClientMgr(ctx)

	if len(cfg.LocalNode) > 0 {
		localClient, err := client2.NewClient(cfg.LocalNode)
		if err != nil {
			cancel()
			return nil, err
		}

		cm.AddClient(localClient)
	}

	for _, nn := range cfg.NetworkNodes {
		c, err := client2.NewClient(nn)
		if err != nil {
			log.Error("can't add new network node client", "err", err, "addr", nn)
		}
		cm.AddClient(c)
	}

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

	return newBotEngine(ctx, cancel, db, cm, wal, cfg.Phoenix.FaucetAmount, cfg.BotName), nil
}

func newBotEngine(ctx context.Context, cnl context.CancelFunc, db *repository.DB, cm *client2.Mgr, wallet *wallet.Wallet, phoenixFaucetAmount uint, botName string) *BotEngine {
	rootCmd := command.Command{
		Emoji:       "🤖",
		Name:        "pagu",
		Help:        "Root Command",
		AppIDs:      entity.AllAppIDs(),
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
	ptCmd := phoenixtestnet.NewPhoenix(wallet, phoenixFaucetAmount, cm, *db)
	zealyCmd := zealy.NewZealy(db, wallet)
	voucherCmd := voucher.NewVoucher(db, wallet, cm)
	marketCmd := market.NewMarket(cm, priceCache)

	return &BotEngine{
		ctx:              ctx,
		cancel:           cnl,
		clientMgr:        cm,
		rootCmd:          rootCmd,
		networkCmd:       netCmd,
		blockchainCmd:    bcCmd,
		phoenixCmd:       ptCmd,
		phoenixClientMgr: cm,
		zealyCmd:         zealyCmd,
		voucherCmd:       voucherCmd,
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
	be.rootCmd.AddSubCommand(be.voucherCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.marketCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.phoenixCmd.GetCommand())

	be.rootCmd.AddHelpSubCommand()
}

func (be *BotEngine) Run(appID entity.AppID, callerID string, tokens []string) command.CommandResult {
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

	for _, middlewareFunc := range cmd.Middlewares {
		if err = middlewareFunc(&cmd, appID, callerID, args...); err != nil {
			log.Error(err.Error())
			return cmd.ErrorResult(errors.New("command is not available. please try again later"))
		}
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
