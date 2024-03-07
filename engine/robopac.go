package engine

import (
	"fmt"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
)

const (
	RoboPacCommandName     = "robopac"
	WalletCommandName      = "wallet"
	RoboPacHelpCommandName = "help"
)

func (be *BotEngine) RegisterRoboPacCommands() {
	cmdWallet := Command{
		Name:    WalletCommandName,
		Desc:    "check the RoboPac wallet balance and address",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.walletHandler,
	}

	cmdHelp := Command{
		Name: RoboPacCommandName,
		Desc: "This is Help for robopac commands",
		Help: "provide the command name as parameter",
		Args: []Args{
			{
				Name:     "sub-command",
				Desc:     "the subcommand you want to see the related help of it",
				Optional: true,
			},
		},
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: nil,
		Handler:     be.robopacHelpHandler,
	}

	cmdRoboPac := Command{
		Name:        RoboPacCommandName,
		Desc:        "robopac related commands",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&cmdWallet, &cmdHelp},
		Handler:     nil,
	}

	be.Cmds = append(be.Cmds, cmdRoboPac)
}

func (be *BotEngine) walletHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	addr, blnc := be.wallet.Address(), be.wallet.Balance()

	result := fmt.Sprintf("Address: https://pacscan.org/address/%s\nBalance: %v PAC\n", addr, utils.FormatNumber(int64(util.ChangeToCoin(blnc))))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) robopacHelpHandler(source AppID, callerID string, args ...string) (*CommandResult, error) {
	if len(args) == 0 {
		return be.help(source, callerID, RoboPacCommandName)
	}
	return be.help(source, callerID, RoboPacCommandName, args[0])
}
