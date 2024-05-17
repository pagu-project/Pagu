package zealy

import (
	"github.com/pagu-project/Pagu/database"
	"github.com/pagu-project/Pagu/engine/command"
	"github.com/pagu-project/Pagu/wallet"
)

const (
	CommandName              = "zealy"
	ClaimCommandName         = "claim"
	StatusCommandName        = "status"
	ImportWinnersCommandName = "import-winners"
	HelpCommandName          = "help"
)

type Zealy struct {
	db     *database.DB
	wallet *wallet.Wallet
}

func NewZealy(
	db *database.DB, wallet *wallet.Wallet,
) Zealy {
	return Zealy{
		db:     db,
		wallet: wallet,
	}
}

func (z *Zealy) GetCommand() command.Command {
	subCmdClaim := command.Command{
		Name: ClaimCommandName,
		Desc: "Claim your Zealy Reward",
		Help: "",
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "Your Pactus address",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     z.claimHandler,
	}

	subCmdStatus := command.Command{
		Name:        StatusCommandName,
		Desc:        "Status of Zealy reward claims",
		Help:        "",
		Args:        nil,
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     z.statusHandler,
	}

	cmdZealy := command.Command{
		Name:        CommandName,
		Desc:        "Zealy Commands",
		Help:        "",
		Args:        nil,
		AppIDs:      command.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
	}

	cmdZealy.AddSubCommand(subCmdClaim)
	cmdZealy.AddSubCommand(subCmdStatus)

	// only accessible from cli
	subCmdImportWinners := command.Command{
		Name: ImportWinnersCommandName,
		Desc: "Import Zealy winners using csv file",
		Help: "",
		Args: []command.Args{
			{
				Name:     "path",
				Desc:     "CSV file path",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI},
		Handler:     z.importWinnersHandler,
	}

	cmdZealy.AddSubCommand(subCmdImportWinners)
	return cmdZealy
}

func (z *Zealy) claimHandler(cmd command.Command, _ command.AppID, callerID string, args ...string) command.CommandResult {
	user, err := z.db.GetZealyUser(callerID)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if len(user.TxHash) > 0 {
		return cmd.FailedResult("You already claimed your reward: https://pacviewer.com/transaction/%s",
			user.TxHash)
	}

	address := args[0]
	txHash, err := z.wallet.TransferTransaction(address, "PaGu Zealy reward distribution", int64(user.Amount))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = z.db.UpdateZealyUser(callerID, txHash); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Zealy reward claimed successfully: https://pacviewer.com/transaction/%s",
		txHash)
}

func (z *Zealy) statusHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	allUsers, err := z.db.GetAllZealyUser()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	total := 0
	totalClaimed := 0
	totalNotClaimed := 0

	totalAmount := 0
	totalClaimedAmount := 0
	totalNotClaimedAmount := 0

	for _, u := range allUsers {
		total++
		totalAmount += int(u.Amount)

		if len(u.TxHash) > 0 {
			totalClaimed++
			totalClaimedAmount += int(u.Amount)
		} else {
			totalNotClaimed++
			totalNotClaimedAmount += int(u.Amount)
		}
	}

	return cmd.SuccessfulResult("Total Users: %v\nTotal Claims: %v\nTotal not remained claims: %v\nTotal Coins: %v PAC\n"+
		"Total claimed coins: %v PAC\nTotal not claimed coins: %v PAC\n", total, totalClaimed, totalNotClaimed,
		totalAmount, totalClaimedAmount, totalNotClaimedAmount,
	)
}
