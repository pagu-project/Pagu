package engine

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
)

const (
	RoboPacCommandName = "robopac"
	HelpCommandName    = "help"
	WalletCommandName  = "wallet"
)

func (be *BotEngine) RegisterRoboPacCommands() {
	cmdHelp := Command{
		Name:    HelpCommandName,
		Desc:    "This is Help!",
		Help:    "",
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.help,
		Args: []Args{
			{Name: "command", Desc: "help", Optional: true},
		},
	}

	cmdWallet := Command{
		Name:    WalletCommandName,
		Desc:    "check the RoboPac wallet balance and address",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.walletHandler,
	}

	cmdRoboPac := Command{
		Name:        RoboPacCommandName,
		Desc:        "robopac related commands",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&cmdHelp, &cmdWallet},
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

func (be *BotEngine) help(source AppID, _ string, args ...string) (*CommandResult, error) {
	helpStr := ""
	if len(args) > 0 {
		cmdName := args[0]
		cmd := be.commandByName(cmdName)
		if cmd == nil {
			return nil, fmt.Errorf("unknown command: %s", cmdName)
		}

		argsStr := ""
		for _, arg := range cmd.Args {
			argsStr += fmt.Sprintf("<%v> ", arg.Name)
		}
		argsStr = argsStr[:len(argsStr)-1]

		helpStr += cmd.Desc
		helpStr += fmt.Sprintf("%v\nUsage: `%v %v`", cmd.Help, cmd.Name, argsStr)
	} else {
		helpStr += "List of available commands:\n"
		for _, cmd := range be.Cmds {
			if !slices.Contains(cmd.AppIDs, source) {
				continue
			}

			padding := 12 - len(cmd.Name)
			helpStr += fmt.Sprintf("`%s`:%s%v\n", cmd.Name, strings.Repeat(" ", padding), cmd.Desc)
		}
	}
	return MakeSuccessfulResult(helpStr), nil
}
