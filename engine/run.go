package engine

import (
	"fmt"
	"strings"
)

// CommandHandler is a function type for handling commands.
type CommandHandler func(be *BotEngine, args []string) (*CommandResult, error)

const (
	CmdClaim            = "claim"             //!
	CmdClaimerInfo      = "claimer-info"      //!
	CmdNodeInfo         = "node-info"         //!
	CmdNetworkStatus    = "network"           //!
	CmdNetworkHealth    = "network-health"    //!
	CmdBotWallet        = "wallet"            //!
	CmdClaimStatus      = "claim-status"      //!
	CmdRewardCalc       = "calc-reward"       //!
	CmdBoosterPayment   = "booster-payment"   //!
	CmdBoosterClaim     = "booster-claim"     //!
	CmdBoosterWhitelist = "booster-whitelist" //!
	CmdBoosterStatus    = "booster-status"    //!
	CmdDefault          = "default"           //!
)

// CommandHandlers is a map of command names to their corresponding handlers.
var CommandHandlers = map[string]CommandHandler{
	CmdClaim:            (*BotEngine).ClaimHandler,
	CmdClaimerInfo:      (*BotEngine).ClaimerInfoHandler,
	CmdNetworkHealth:    (*BotEngine).NetworkHealthHandler,
	CmdNodeInfo:         (*BotEngine).NodeInfoHandler,
	CmdNetworkStatus:    (*BotEngine).NetworkStatusHandler,
	CmdBotWallet:        (*BotEngine).BotWalletHandler,
	CmdClaimStatus:      (*BotEngine).ClaimStatusHandler,
	CmdRewardCalc:       (*BotEngine).RewardCalcHandler,
	CmdBoosterPayment:   (*BotEngine).BoosterPaymentHandler,
	CmdBoosterClaim:     (*BotEngine).BoosterClaimHandler,
	CmdBoosterWhitelist: (*BotEngine).BoosterWhitelistHandler,
	CmdBoosterStatus:    (*BotEngine).BoosterStatusHandler,
	CmdDefault:          (*BotEngine).DefaultCommandHandler,
}

// The input is always string.
//
//	The input format is like: [Command] <Arguments ...>
//
// The output is always string, but format might be JSON. ???

func (be *BotEngine) Run(appID AppID, callerID string, inputs []string) (*CommandResult, error) {
	cmd, args := be.parseQuery(strings.Join(inputs, " "))
	handler, found := CommandHandlers[cmd]
	if !found {
		// handler unknown commands.
		return &CommandResult{
			Successful: false,
			Message:    fmt.Sprintf("unknown command: %s", args[0]),
		}, nil
	}
	return handler(be, args)
}

func (be *BotEngine) parseQuery(query string) (string, []string) {
	subs := strings.Split(query, " ")
	if len(subs) == 0 {
		return "", nil
	}

	return subs[0], subs[1:]
}
