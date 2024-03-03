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
	Name    string
	Desc    string
	Help    string
	Args    []Args
	AppIDs  []AppID
	Handler func(source AppID, callerID string, args ...string) (*CommandResult, error)
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
