package engine

import (
	"strings"
)

// CommandHandler is a function type for handling commands.
type CommandHandler func(be *BotEngine, args []string) (string, error)

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
	CmdClaim:            ClaimHandler,
	CmdClaimerInfo:      ClaimerInfoHandler,
	CmdNetworkHealth:    NetworkHealthHandler,
	CmdNodeInfo:         NodeInfoHandler,
	CmdNetworkStatus:    NetworkStatusHandler,
	CmdBotWallet:        BotWalletHandler,
	CmdClaimStatus:      ClaimStatusHandler,
	CmdRewardCalc:       RewardCalcHandler,
	CmdBoosterPayment:   BoosterPaymentHandler,
	CmdBoosterClaim:     BoosterClaimHandler,
	CmdBoosterWhitelist: BoosterWhitelistHandler,
	CmdBoosterStatus:    BoosterStatusHandler,
	CmdDefault:          DefaultCommandHandler,
}

// The input is always string.
//
//	The input format is like: [Command] <Arguments ...>
//
// The output is always string, but format might be JSON. ???
func (be *BotEngine) Run(input string) (string, error) {
	cmd, args := be.parseQuery(input)
	handler, found := CommandHandlers[cmd]
	if !found {
		// handler unknown commands.
		return DefaultCommandHandler(be, args)
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
