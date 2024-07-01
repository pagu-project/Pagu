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
}

func NewPhoenix(wallet *wallet.Wallet, faucetAmount uint, clientMgr *client.Mgr, db repository.DB) Phoenix {
	return Phoenix{
		wallet:       wallet,
		clientMgr:    clientMgr,
		db:           db,
		faucetAmount: faucetAmount,
	}
}

func (pt *Phoenix) GetCommand() command.Command {
	middlewareHandler := command.NewMiddlewareHandler(&pt.db, pt.wallet)

	subCmdStatus := command.Command{
		Name:        StatusCommandName,
		Help:        "Phoenix test-network statistics",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: []command.MiddlewareFunc{middlewareHandler.CreateUser},
		Handler:     pt.networkStatusHandler,
	}

	subCmdFaucet := command.Command{
		Name: FaucetCommandName,
		Help: fmt.Sprintf("Get %d tPAC Coins on Phoenix Testnet for Testing your code or project", pt.faucetAmount),
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
		Help:        "Phoenix Testnet tools and utils for developers",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskTest,
	}

	cmdPhoenix.AddSubCommand(subCmdFaucet)
	cmdPhoenix.AddSubCommand(subCmdStatus)

	return cmdPhoenix
}
