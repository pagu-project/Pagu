package engine

import (
	"fmt"
	"slices"

	"github.com/kehiy/RoboPac/log"
)

const (
	CmdClaim            = "claim"
	CmdClaimerInfo      = "claimer-info"
	CmdNodeInfo         = "node-info"
	CmdNetworkStatus    = "network"
	CmdNetworkHealth    = "network-health"
	CmdBotWallet        = "wallet"
	CmdClaimStatus      = "claim-status"
	CmdRewardCalc       = "calc-reward"
	CmdBoosterPayment   = "booster-payment"
	CmdBoosterClaim     = "booster-claim"
	CmdBoosterWhitelist = "booster-whitelist"
	CmdBoosterStatus    = "booster-status"
	CmdDefault          = "default"
)

func (be *BotEngine) RegisterCommands() {
}

func (be *BotEngine) Commands() []Command {
	return be.Cmds
}

func (be *BotEngine) Run(appID AppID, callerID string, inputs []string) (*CommandResult, error) {
	log.Debug("run command", "callerID", callerID, "inputs", inputs)

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
