package phoenix

import (
	"github.com/pagu-project/Pagu/database"
	"github.com/pagu-project/Pagu/internal/engine/command"
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
	wallet    *wallet.Wallet
	db        database.DB
	clientMgr *client.Mgr
}

func NewPhoenix(wallet *wallet.Wallet,
	clientMgr *client.Mgr, db database.DB,
) Phoenix {
	return Phoenix{
		wallet:    wallet,
		clientMgr: clientMgr,
		db:        db,
	}
}

func (pt *Phoenix) GetCommand() command.Command {
	subCmdFaucet := command.Command{
		Name: FaucetCommandName,
		Desc: "Get 5 tPAC Coins on Phoenix Testnet for Testing your code or project",
		Help: "There is a limit that you can only get faucets 1 time per day with each user ID and address",
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "your testnet address [example: tpc1z...]",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.faucetHandler,
	}

	subCmdWallet := command.Command{
		Name:        WalletCommandName,
		Desc:        "Check the status of RoboPac faucet address wallet on Phoenix network",
		Help:        "",
		Args:        nil,
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.walletHandler,
	}

	subCmdHealth := command.Command{
		Name:        HealthCommandName,
		Desc:        "Checking Phoenix test-network health status",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.networkHealthHandler,
	}

	subCmdStatus := command.Command{
		Name:        StatusCommandName,
		Desc:        "Phoenix test-network statistics",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.networkStatusHandler,
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
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.nodeInfoHandler,
	}

	cmdPhoenix := command.Command{
		Name:        CommandName,
		Desc:        "Phoenix Testnet tools and utils for developers",
		Help:        "",
		Args:        nil,
		AppIDs:      command.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
	}

	cmdPhoenix.AddSubCommand(subCmdFaucet)
	cmdPhoenix.AddSubCommand(subCmdWallet)
	cmdPhoenix.AddSubCommand(subCmdHealth)
	cmdPhoenix.AddSubCommand(subCmdStatus)
	cmdPhoenix.AddSubCommand(subCmdNodeInfo)

	return cmdPhoenix
}
