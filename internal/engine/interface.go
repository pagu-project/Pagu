package engine

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

type IEngine interface {
	Run(appID entity.AppID, callerID string, tokens []string) (*command.CommandResult, error)
	Commands() []command.Command
}
