package phoenix

import (
	"fmt"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

const (
	CommandName         = "phoenix"
	FaucetCommandName   = "faucet"
	WalletCommandName   = "wallet"
	StatusCommandName   = "status"
	HealthCommandName   = "health"
	NodeInfoCommandName = "node-info"
	HelpCommandName     = "help"
)

type Phoenix struct {
	wallet       *wallet.Wallet
	db           repository.DB
	clientMgr    *client.Mgr
	faucetAmount uint
	targetMask   int
}

func NewPhoenix(wallet *wallet.Wallet, faucetAmount uint, clientMgr *client.Mgr, db repository.DB, target int) Phoenix {
	return Phoenix{
		wallet:       wallet,
		clientMgr:    clientMgr,
		db:           db,
		faucetAmount: faucetAmount,
		targetMask:   target,
	}
}

func (pt *Phoenix) GetCommand() command.Command {
	middlewareHandler := command.NewMiddlewareHandler(&pt.db, pt.wallet)

	subCmdStatus := command.Command{
		Name:        StatusCommandName,
		Desc:        "Phoenix test-network statistics",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: []command.MiddlewareFunc{middlewareHandler.CreateUser},
		Handler:     pt.networkStatusHandler,
	}

	subCmdFaucet := command.Command{
		Name: FaucetCommandName,
		Desc: fmt.Sprintf("Get %d tPAC Coins on Phoenix Testnet for Testing your code or project", pt.faucetAmount),
		Help: "There is a limit that you can only get faucets 1 time per day with each user ID and address",
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "your testnet address [example: tpc1z...]",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: []command.MiddlewareFunc{middlewareHandler.CreateUser, middlewareHandler.WalletBalance},
		Handler:     pt.faucetHandler,
	}

	cmdPhoenix := command.Command{
		Name:        CommandName,
		Desc:        "Phoenix Testnet tools and utils for developers",
		Help:        "",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
		TargetMask:  pt.targetMask,
	}

	cmdPhoenix.AddSubCommand(subCmdFaucet)
	cmdPhoenix.AddSubCommand(subCmdStatus)

	return cmdPhoenix
}
