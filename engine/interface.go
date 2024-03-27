package engine

import (
	"github.com/robopac-project/RoboPac/engine/command"
)

type IEngine interface {
	Run(appID command.AppID, callerID string, tokens []string) (*command.CommandResult, error)
	Commands() []command.Command
}
