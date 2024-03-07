package engine

import (
	"fmt"
	"slices"
	"strings"
)

type AppID int

const (
	AppIdCLI     AppID = 1
	AppIdDiscord AppID = 2
)

type Args struct {
	Name     string
	Desc     string
	Optional bool
}

type Command struct {
	Name        string
	Desc        string
	Help        string
	Args        []Args //! should be nil for commands.
	AppIDs      []AppID
	SubCommands []*Command
	Handler     func(source AppID, callerID string, args ...string) (*CommandResult, error)
}

func (be *BotEngine) RegisterAllCommands() {
	be.RegisterBlockchainCommands()
	be.RegisterBoosterCommands()
	be.RegisterNetworkCommands()
	be.RegisterP2PMarketCommands()
	be.RegisterRoboPacCommands()
	be.RegisterTestNetRewardsCommands()
}

func (be *BotEngine) Run(appID AppID, callerID string, inputs []string) (*CommandResult, error) {
	be.logger.Debug("run command", "callerID", callerID, "inputs", inputs)

	cmdName := inputs[0]
	subCmdName := inputs[1]
	cmd := be.subCommandByName(cmdName, subCmdName)
	if cmd == nil {
		return nil, fmt.Errorf("unknown command: %s", cmdName)
	}
	if !cmd.HasAppId(appID) {
		return nil, fmt.Errorf("unauthorized appID: %v", appID)
	}

	args := inputs[2:]
	err := cmd.CheckArgs(args)
	if err != nil {
		return nil, err
	}

	return cmd.Handler(appID, callerID, args...)
}

func (be *BotEngine) subCommandByName(cmdName, subCmdName string) *Command {
	cmdIndex := slices.IndexFunc(be.Cmds, func(cmd Command) bool {
		return cmd.Name == cmdName
	})

	if cmdIndex == -1 {
		return nil
	}

	sCmdIndex := slices.IndexFunc(be.Cmds[cmdIndex].SubCommands, func(cmd *Command) bool {
		return cmd.Name == subCmdName
	})

	if sCmdIndex == -1 {
		return nil
	}

	return be.Cmds[cmdIndex].SubCommands[sCmdIndex]
}

func (be *BotEngine) commandByName(cmdName string) *Command {
	cmdIndex := slices.IndexFunc(be.Cmds, func(cmd Command) bool {
		return cmd.Name == cmdName
	})

	if cmdIndex == -1 {
		return nil
	}

	return &be.Cmds[cmdIndex]
}

func (be *BotEngine) Commands() []Command {
	return be.Cmds
}

type CommandResult struct {
	Message    string
	Successful bool
}

func MakeSuccessfulResult(message string, a ...interface{}) *CommandResult {
	return &CommandResult{
		Message:    fmt.Sprintf(message, a...),
		Successful: true,
	}
}

func MakeFailedResult(message string, a ...interface{}) *CommandResult {
	return &CommandResult{
		Message:    fmt.Sprintf(message, a...),
		Successful: false,
	}
}

func (cmd *Command) CheckArgs(input []string) error {
	minArg := len(cmd.Args)
	maxArg := len(cmd.Args)

	for _, arg := range cmd.Args {
		if arg.Optional {
			minArg--
		}
	}

	if len(input) < minArg || len(input) > maxArg {
		return fmt.Errorf("incorrect number of arguments, expected %d but got %d", minArg, len(input))
	}

	return nil
}

func (cmd *Command) HasAppId(appID AppID) bool {
	return slices.Contains(cmd.AppIDs, appID)
}

func (be *BotEngine) help(_ AppID, _ string, args ...string) (*CommandResult, error) {
	var helpMsg strings.Builder

	cmdName := args[0]

	if args[1] == "" {
		cmd := be.commandByName(cmdName)
		if cmd != nil {
			helpMsg.WriteString(fmt.Sprintf("Help of %s Command\nDesc: %s\nHelp: %s\n\nSubCommands:", cmd.Name, cmd.Desc, cmd.Help))
			for i, sc := range cmd.SubCommands {
				helpMsg.WriteString(fmt.Sprintf("%v-%s", i, sc.Name))
			}

			return MakeSuccessfulResult(helpMsg.String()), nil
		}
		return MakeFailedResult("can't find the command: %s", cmdName), nil
	}

	subCmdName := args[1]

	subCmd := be.subCommandByName(cmdName, subCmdName)
	if subCmd == nil {
		return MakeFailedResult("can't find the sub command command: %s", subCmdName), nil
	}

	helpMsg.WriteString(fmt.Sprintf("Help of %s Command\nDesc: %s\nHelp: %s\n\nArgs:", subCmd.Name, subCmd.Desc, subCmd.Help))
	for i, a := range subCmd.Args {
		helpMsg.WriteString(fmt.Sprintf("%v-%s\nDesc: %s\n optional: %v", i, a.Name, a.Desc, a.Optional))
	}

	return MakeSuccessfulResult(helpMsg.String()), nil
}
