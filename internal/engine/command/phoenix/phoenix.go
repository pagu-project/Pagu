package phoenix

import (
	"context"
	"fmt"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/amount"
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
	ctx          context.Context
	wallet       wallet.IWallet
	db           repository.Database
	clientMgr    client.Manager
	faucetAmount amount.Amount
}

func NewPhoenix(ctx context.Context, wlt wallet.IWallet, faucetAmount amount.Amount,
	clientMgr client.Manager, db repository.Database,
) *Phoenix {
	return &Phoenix{
		ctx:          ctx,
		wallet:       wlt,
		clientMgr:    clientMgr,
		db:           db,
		faucetAmount: faucetAmount,
	}
}

func (pt *Phoenix) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(pt.db, pt.wallet)

	subCmdStatus := &command.Command{
		Name:        StatusCommandName,
		Help:        "Phoenix Testnet statistics",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: nil,
		Handler:     pt.networkStatusHandler,
		TargetFlag:  command.TargetMaskTest,
	}

	subCmdFaucet := &command.Command{
		Name: FaucetCommandName,
		Help: fmt.Sprintf("Get %f tPAC Coins on Phoenix Testnet for Testing your code or project", pt.faucetAmount.ToPAC()),
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "your testnet address [example: tpc1z...]",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: []command.MiddlewareFunc{middlewareHandler.WalletBalance},
		Handler:     pt.faucetHandler,
		TargetFlag:  command.TargetMaskTest,
	}

	cmdPhoenix := &command.Command{
		Name:        CommandName,
		Help:        "Phoenix Testnet tools and utils for developers",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskTest,
	}

	cmdPhoenix.AddSubCommand(subCmdFaucet)
	cmdPhoenix.AddSubCommand(subCmdStatus)

	return cmdPhoenix
}
