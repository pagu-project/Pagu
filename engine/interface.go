package engine

import (
	"time"

	"github.com/robopac-project/RoboPac/engine/command"
)

type IEngine interface {
	Run(appID command.AppID, callerID string, tokens []string) (*command.CommandResult, error)
	Commands() []command.Command
}

type NetHealthResponse struct {
	HealthStatus    bool
	CurrentTime     time.Time
	LastBlockTime   time.Time
	LastBlockHeight uint32
	TimeDifference  int64
}
