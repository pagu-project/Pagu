package engine

import (
	"fmt"
	"slices"

	"github.com/kehiy/RoboPac/log"
)

const (
	ClaimCommandName       = "claim"
	ClaimerInfoCommandName = "claimer-info"
	ClaimStatusCommandName = "claim-status"

	NodeInfoCommandName      = "node-info"
	NetworkStatusCommandName = "network"
	NetworkHealthCommandName = "network-health"

	HelpCommandName       = "help"
	WalletCommandName     = "wallet"
	CalcRewardCommandName = "calc-reward"

	BoosterPaymentCommandName   = "booster-payment"
	BoosterClaimCommandName     = "booster-claim"
	BoosterWhitelistCommandName = "booster-whitelist"
	BoosterStatusCommandName    = "booster-status"

	DepositAddressCommandName = "deposit-address"
)

func (be *BotEngine) RegisterCommands() {
	cmdClaim := Command{
		Name: ClaimCommandName,
		Desc: "claim your test-net rewards",
		Help: "",
		Args: []Args{
			{
				Name:     "mainnet-address",
				Desc:     "your main-net (validator) address like: pc1p...",
				Optional: false,
			},
			{
				Name:     "testnet-address",
				Desc:     "your test-net (validator) address like: tpc1p...",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.claimHandler,
	}

	cmdClaimerInfo := Command{
		Name: ClaimerInfoCommandName,
		Desc: "check your claim status with your testnet validator address",
		Help: "",
		Args: []Args{
			{
				Name:     "testnet-address",
				Desc:     "your test-net (validator) address like: tpc1p...",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.claimerInfoHandler,
	}

	cmdClaimStatus := Command{
		Name:    ClaimStatusCommandName,
		Desc:    "check the status of testnet rewards claiming",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.claimStatusHandler,
	}

	cmdNodeInfo := Command{
		Name: NodeInfoCommandName,
		Desc: "check the information of a node by providing it's validator address",
		Help: "",
		Args: []Args{
			{
				Name:     "validator-address",
				Desc:     "your validator address",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.nodeInfoHandler,
	}

	cmdNetworkHealth := Command{
		Name:    NetworkHealthCommandName,
		Desc:    "checking network health status",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.networkHealthHandler,
	}

	cmdNetworkStatus := Command{
		Name:    NetworkStatusCommandName,
		Desc:    "network statistics",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.networkStatusHandler,
	}

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

	cmdCalcReward := Command{
		Name: CalcRewardCommandName,
		Desc: "claculate how much PAC coins you will earn with your validator stakes",
		Help: "",
		Args: []Args{
			{
				Name:     "stake-amount",
				Desc:     "amount of stake in your validator (1-1000)",
				Optional: false,
			},
			{
				Name:     "time-interval",
				Desc:     "after one: day | month | year",
				Optional: true,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.calcRewardHandler,
	}

	cmdBoosterPayment := Command{
		Name: BoosterPaymentCommandName,
		Desc: "make a payment link for booster program",
		Help: "",
		Args: []Args{
			{
				Name:     "twitter-name",
				Desc:     "your twitter user name without @",
				Optional: false,
			},
			{
				Name:     "validator-address",
				Desc:     "your validator address to be registered",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.boosterPaymentHandler,
	}

	cmdBoosterClaim := Command{
		Name: BoosterClaimCommandName,
		Desc: "claim your booster program stakes",
		Help: "",
		Args: []Args{
			{
				Name:     "twitter-name",
				Desc:     "your twitter user name without @",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.boosterClaimHandler,
	}

	cmdBoosterWhitelist := Command{
		Name: BoosterWhitelistCommandName,
		Desc: "whitelist a user for booster program (admin only)",
		Help: "",
		Args: []Args{
			{
				Name:     "twitter-name",
				Desc:     "your twitter user name without @",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.boosterWhitelistHandler,
	}

	cmdBoosterStatus := Command{
		Name:    BoosterStatusCommandName,
		Desc:    "status of booster program claims and ...",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.boosterStatusHandler,
	}

	cmdDepositAddress := Command{
		Name:    DepositAddressCommandName,
		Desc:    "create a deposit address for P2P offer",
		Help:    "it will show your address if you already have an deposit address",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.depositAddressHandler,
	}

	//! test-net reward commands
	be.Cmds = append(be.Cmds, cmdClaim)
	be.Cmds = append(be.Cmds, cmdClaimerInfo)
	be.Cmds = append(be.Cmds, cmdClaimStatus)

	//! network info commands
	be.Cmds = append(be.Cmds, cmdNodeInfo)
	be.Cmds = append(be.Cmds, cmdNetworkHealth)
	be.Cmds = append(be.Cmds, cmdNetworkStatus)

	//! bot info and util commands
	be.Cmds = append(be.Cmds, cmdHelp)
	be.Cmds = append(be.Cmds, cmdWallet)
	be.Cmds = append(be.Cmds, cmdCalcReward)

	//! booster program commands
	be.Cmds = append(be.Cmds, cmdBoosterPayment)
	be.Cmds = append(be.Cmds, cmdBoosterClaim)
	be.Cmds = append(be.Cmds, cmdBoosterWhitelist)
	be.Cmds = append(be.Cmds, cmdBoosterStatus)

	//! P2P offer commands
	be.Cmds = append(be.Cmds, cmdDepositAddress)
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
