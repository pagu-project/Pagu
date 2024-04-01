package phoenixtestnet

import (
	"github.com/robopac-project/RoboPac/client"
	"github.com/robopac-project/RoboPac/database"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/wallet"
)

const (
	PhoenixTestnetCommandName = "phoenix"
	FaucetCommandName         = "faucet"
	WalletCommandName         = "wallet"
	StatusCommandName         = "status"
	HealthCommandName         = "health"
	BlockChainHelpCommandName = "help"
)

type PhoenixTestnet struct {
	wallet    wallet.IWallet
	db        database.DB
	clientMgr *client.Mgr
}

func NewPhoenixTestnet(wallet wallet.IWallet,
	clientMgr *client.Mgr, db database.DB,
) PhoenixTestnet {
	return PhoenixTestnet{
		wallet:    wallet,
		clientMgr: clientMgr,
		db:        db,
	}
}

func (pt *PhoenixTestnet) GetCommand() command.Command {
	subCmdFaucet := command.Command{
		Name: FaucetCommandName,
		Desc: "Get 5 tPAC Coins on Phoenix Testnet for Testing your code or project",
		Help: "There is a limit that you can only get faucets 1 time per day with each user ID and address",
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "your testnet address, example: tpc1z....",
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

	cmdPhoenixTestnet := command.Command{
		Name:        PhoenixTestnetCommandName,
		Desc:        "Phoenix Testnet tools and utils for developers",
		Help:        "",
		Args:        nil,
		AppIDs:      command.AllAppIDs(),
		SubCommands: make([]command.Command, 2),
		Handler:     nil,
	}

	cmdPhoenixTestnet.AddSubCommand(subCmdFaucet)
	cmdPhoenixTestnet.AddSubCommand(subCmdWallet)

	return cmdPhoenixTestnet
}

func (pt *PhoenixTestnet) faucetHandler(cmd command.Command, _ command.AppID, callerID string, args ...string) command.CommandResult {
	if !pt.db.HasUser(callerID) {
		if err := pt.db.AddUser(
			&database.User{
				ID: callerID,
			},
		); err != nil {
			return cmd.ErrorResult(err)
		}
	}

	if !pt.db.CanGetFaucet(callerID) {
		return cmd.FailedResult("Uh, you used your share of faucets today!")
	}

	if pt.wallet.Balance() < 5 {
		return cmd.FailedResult("RoboPac Phoenix wallet is empty, please contact the team!")
	}

	toAddr := args[0]
	txID, err := pt.wallet.TransferTransaction(toAddr, "Phoenix Testnet RoboPac Faucet", 5) //! define me on config?
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = pt.db.AddFaucet(&database.Faucet{
		Address:         toAddr,
		Amount:          5,
		TransactionHash: txID,
		UserID:          callerID,
	}); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("You got %d tPAC in %s address on Phoenix Testnet!", 5, toAddr)
}

func (pt *PhoenixTestnet) walletHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	return cmd.SuccessfulResult("RoboPac Phoenix Address: %s\nBalance: %d", pt.wallet.Address(), pt.wallet.Balance())
}
