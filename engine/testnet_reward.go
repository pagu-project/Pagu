package engine

import (
	"errors"
	"fmt"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
)

const (
	TestNetRewardCommandName     = "testnet-reward"
	ClaimCommandName             = "claim"
	ClaimerInfoCommandName       = "claimer-info"
	ClaimStatusCommandName       = "claims-status"
	TestNetRewardHelpCommandName = "help"
)

func (be *BotEngine) RegisterTestNetRewardsCommands() {
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

	cmdHelp := Command{
		Name: RoboPacCommandName,
		Desc: "This is Help for testnet rewards commands",
		Help: "provide the command name as parameter",
		Args: []Args{
			{
				Name:     "sub-command",
				Desc:     "the subcommand you want to see the related help of it",
				Optional: true,
			},
		},
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: nil,
		Handler:     be.testnetRewardHelpHandler,
	}

	cmdTestNetReward := Command{
		Name:        TestNetRewardCommandName,
		Desc:        "claiming your testnet earned rewards",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&cmdClaim, &cmdClaimStatus, &cmdClaimerInfo, &cmdHelp},
		Handler:     nil,
	}

	be.Cmds = append(be.Cmds, cmdTestNetReward)
}

func (be *BotEngine) claimerInfoHandler(_ AppID, _ string, args ...string) (*CommandResult, error) {
	be.RLock()
	defer be.RUnlock()

	testNetValAddr := args[0]

	claimer := be.store.ClaimerInfo(testNetValAddr)
	if claimer == nil {
		return nil, errors.New("not found")
	}

	return &CommandResult{
		Successful: true,
		Message: fmt.Sprintf("TestNet Address: %s\namount: %v PACs\nIsClaimed: %v\n txHash: %s",
			args[0], util.ChangeToString(claimer.TotalReward), claimer.IsClaimed(), claimer.ClaimedTxID),
	}, nil
}

func (be *BotEngine) claimHandler(_ AppID, callerID string, args ...string) (*CommandResult, error) {
	be.Lock()
	defer be.Unlock()

	mainnetAddr := args[0]
	testnetAddr := args[1]

	be.logger.Info("new claim request", "mainnetAddr", mainnetAddr, "testnetAddr", testnetAddr, "discordID", callerID)

	valInfo, _ := be.clientMgr.GetValidatorInfo(mainnetAddr)
	if valInfo != nil {
		return nil, errors.New("this address is already a staked validator")
	}

	if utils.ChangeToCoin(be.wallet.Balance()) <= 500 {
		be.logger.Warn("bot wallet hasn't enough balance")
		return nil, errors.New("insufficient wallet balance")
	}

	claimer := be.store.ClaimerInfo(testnetAddr)
	if claimer == nil {
		return nil, errors.New("claimer not found")
	}

	if claimer.DiscordID != callerID {
		be.logger.Warn("try to claim other's reward", "claimer", claimer.DiscordID, "discordID", callerID)
		return nil, errors.New("invalid claimer")
	}

	if claimer.IsClaimed() {
		return nil, errors.New("this claimer have already claimed rewards")
	}

	pubKey, err := be.clientMgr.FindPublicKey(mainnetAddr, true)
	if err != nil {
		return nil, err
	}

	memo := "TestNet reward claim from RoboPac"
	txID, err := be.wallet.BondTransaction(pubKey, mainnetAddr, memo, claimer.TotalReward)
	if err != nil {
		return nil, err
	}

	if txID == "" {
		return nil, errors.New("can't send bond transaction")
	}

	be.logger.Info("new bond transaction sent", "txID", txID)

	err = be.store.AddClaimTransaction(testnetAddr, txID)
	if err != nil {
		be.logger.Panic("unable to add the claim transaction",
			"error", err,
			"discordID", callerID,
			"testnetAddr", testnetAddr,
			"txID", txID,
		)

		return nil, err
	}

	return &CommandResult{
		Successful: true,
		Message:    fmt.Sprintf("Reward claimed successfully✅\nYour claim transaction: https://pacscan.org/transactions/%s", txID),
	}, nil
}

func (be *BotEngine) claimStatusHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	cs := be.store.ClaimStatus()

	result := fmt.Sprintf("Claimed rewards count: %v\nClaimed coins: %v PAC's\nNot-claimed rewards count: %v\nNot-claim coins: %v PAC's\n",
		cs.Claimed, util.ChangeToString(cs.ClaimedAmount), cs.NotClaimed, util.ChangeToString(cs.NotClaimedAmount))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) testnetRewardHelpHandler(source AppID, callerID string, args ...string) (*CommandResult, error) {
	if len(args) == 0 {
		return be.help(source, callerID, TestNetRewardCommandName)
	}
	return be.help(source, callerID, TestNetRewardCommandName, args[0])
}
