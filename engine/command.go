package engine

import (
	"fmt"
	"slices"
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

func (be *BotEngine) Run(appID AppID, callerID string, inputs []string) (*CommandResult, error) {
	be.logger.Debug("run command", "callerID", callerID, "inputs", inputs)

	cmdName := inputs[0]
	cmd := be.commandByName(cmdName)
	if cmd == nil {
		return nil, fmt.Errorf("unknown command: %s", cmdName)
	}
	if !cmd.HasAppId(appID) {
		return nil, fmt.Errorf("unauthorized appID: %v", appID)
	}
	args := inputs[1:]
	err := cmd.CheckArgs(args)
	if err != nil {
		return nil, err
	}

	return cmd.Handler(appID, callerID, args...)
}

func (be *BotEngine) commandByName(cmdName string) *Command {
	foundIndex := slices.IndexFunc(be.Cmds, func(cmd Command) bool {
		return cmd.Name == cmdName
	})

	if foundIndex == -1 {
		return nil
	}

	return &be.Cmds[foundIndex]
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
