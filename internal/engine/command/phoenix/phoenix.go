package phoenix

import (
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

func NewPhoenix(wallet *wallet.Wallet, faucetAmount uint, clientMgr *client.Mgr, db repository.DB,
) Phoenix {
	return Phoenix{
		wallet:       wallet,
		clientMgr:    clientMgr,
		db:           db,
		faucetAmount: faucetAmount,
	}
}

func (pt *Phoenix) GetCommand() command.Command {
	middlewareHandler := command.NewMiddlewareHandler(&pt.db, pt.wallet)

	/*
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
			MiddlewareHandler:     pt.faucetHandler,
		}


		subCmdHealth := command.Command{
			Name:        HealthCommandName,
			Desc:        "Checking Phoenix test-network health status",
			Help:        "",
			Args:        []command.Args{},
			SubCommands: nil,
			AppIDs:      entity.AllAppIDs(),
			MiddlewareHandler:     pt.networkHealthHandler,
		}

		subCmdNodeInfo := command.Command{
			Name: NodeInfoCommandName,
			Desc: "View the information of a node running on Phoenix test-network",
			Help: "Provide your validator address on the specific node to get the validator and node info (Phoenix network)",
			Args: []command.Args{
				{
					Name:     "validator_address",
					Desc:     "Your validator address start with tpc1p...",
					Optional: false,
				},
			},
			SubCommands: nil,
			AppIDs:      entity.AllAppIDs(),
			MiddlewareHandler:     pt.nodeInfoHandler,
		}
	*/

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

	subCmdWallet := command.Command{
		Name:        WalletCommandName,
		Desc:        "Check the status of Pagu faucet address wallet on Phoenix network",
		Help:        "",
		Args:        nil,
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     pt.walletHandler,
	}

	cmdPhoenix := command.Command{
		Name:        CommandName,
		Desc:        "Phoenix Testnet tools and utils for developers",
		Help:        "",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
	}

	// cmdPhoenix.AddSubCommand(subCmdFaucet)
	// cmdPhoenix.AddSubCommand(subCmdHealth)
	// cmdPhoenix.AddSubCommand(subCmdNodeInfo)
	cmdPhoenix.AddSubCommand(subCmdWallet)
	cmdPhoenix.AddSubCommand(subCmdStatus)

	return cmdPhoenix
}
