package testnetreward

import (
	"context"
	"fmt"
	"sync"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/engine/command"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/utils"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/pactus-project/pactus/util"
)

const (
	TestNetRewardCommandName     = "testnet-reward"
	ClaimCommandName             = "claim"
	ClaimerInfoCommandName       = "claimer-info"
	ClaimStatusCommandName       = "status"
	TestNetRewardHelpCommandName = "help"
)

type TestNetReward struct {
	sync.RWMutex //! remove this.

	ctx       context.Context
	AdminIDs  []string
	store     store.IStore
	wallet    wallet.IWallet
	clientMgr *client.Mgr
}

func NewTestNetReward(ctx context.Context,
	adminIDs []string,
	store store.IStore,
	wallet wallet.IWallet,
	clientMgr *client.Mgr,
) *TestNetReward {
	return &TestNetReward{
		ctx:       ctx,
		AdminIDs:  adminIDs,
		store:     store,
		wallet:    wallet,
		clientMgr: clientMgr,
	}
}

func (tr *TestNetReward) GetCommand() *command.Command {
	subCmdClaim := command.Command{
		Name: ClaimCommandName,
		Desc: "Claim your testnet rewards",
		Help: "Provide your mainnet validator address and testnet validator address which was eligible",
		Args: []command.Args{
			{
				Name:     "mainnet-address",
				Desc:     "Your main-net (validator) address like: pc1p...",
				Optional: false,
			},
			{
				Name:     "testnet-address",
				Desc:     "Your test-net (validator) address like: tpc1p...",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler:     tr.claimHandler,
	}

	subCmdClaimerInfo := command.Command{
		Name: ClaimerInfoCommandName,
		Desc: "Check your claim status with your testnet validator address",
		Help: "",
		Args: []command.Args{
			{
				Name:     "testnet-address",
				Desc:     "Your test-net (validator) address like: tpc1p...",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler:     tr.claimerInfoHandler,
	}

	subCmdClaimStatus := command.Command{
		Name:        ClaimStatusCommandName,
		Desc:        "The status of testnet rewards claiming",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler:     tr.claimStatusHandler,
	}

	cmdTestNetReward := command.Command{
		Name:        TestNetRewardCommandName,
		Desc:        "Claiming your testnet earned rewards",
		Help:        "",
		Args:        nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		SubCommands: []*command.Command{&subCmdClaim, &subCmdClaimStatus, &subCmdClaimerInfo},
		Handler:     nil,
	}

	cmdTestNetReward.AddSubCommand(&subCmdClaim)
	cmdTestNetReward.AddSubCommand(&subCmdClaimStatus)
	cmdTestNetReward.AddSubCommand(&subCmdClaimerInfo)

	cmdTestNetReward.AddHelpSubCommand()

	return &cmdTestNetReward
}

func (tr *TestNetReward) claimerInfoHandler(cmd *command.Command, _ command.AppID, _ string, args ...string) *command.CommandResult {
	tr.RLock()
	defer tr.RUnlock()

	testNetValAddr := args[0]

	claimer := tr.store.ClaimerInfo(testNetValAddr)
	if claimer == nil {
		return &command.CommandResult{
			Error:      "not found",
			Successful: false,
		}
	}

	return &command.CommandResult{
		Successful: true,
		Message: fmt.Sprintf("TestNet Address: %s\namount: %v PACs\nIsClaimed: %v\n txHash: %s",
			args[0], util.ChangeToString(claimer.TotalReward), claimer.IsClaimed(), claimer.ClaimedTxID),
	}
}

func (tr *TestNetReward) claimHandler(cmd *command.Command, _ command.AppID, callerID string, args ...string) *command.CommandResult {
	tr.Lock()
	defer tr.Unlock()

	mainnetAddr := args[0]
	testnetAddr := args[1]

	log.Info("new claim request", "mainnetAddr", mainnetAddr, "testnetAddr", testnetAddr, "discordID", callerID)

	valInfo, _ := tr.clientMgr.GetValidatorInfo(mainnetAddr)
	if valInfo != nil {
		return &command.CommandResult{
			Error:      "this address is already a staked validator",
			Successful: false,
		}
	}

	if utils.ChangeToCoin(tr.wallet.Balance()) <= 500 {
		log.Warn("bot wallet hasn't enough balance")
		return &command.CommandResult{
			Error:      "insufficient wallet balance",
			Successful: false,
		}
	}

	claimer := tr.store.ClaimerInfo(testnetAddr)
	if claimer == nil {
		return &command.CommandResult{
			Error:      "claimer not found",
			Successful: false,
		}
	}

	if claimer.DiscordID != callerID {
		log.Warn("try to claim other's reward", "claimer", claimer.DiscordID, "discordID", callerID)
		return &command.CommandResult{
			Error:      "invalid claimer",
			Successful: false,
		}
	}

	if claimer.IsClaimed() {
		return &command.CommandResult{
			Error:      "this claimer have already claimed rewards",
			Successful: false,
		}
	}

	pubKey, err := tr.clientMgr.FindPublicKey(mainnetAddr, true)
	if err != nil {
		return &command.CommandResult{
			Error:      err.Error(),
			Successful: false,
		}
	}

	memo := "TestNet reward claim from RoboPac"
	txID, err := tr.wallet.BondTransaction(pubKey, mainnetAddr, memo, claimer.TotalReward)
	if err != nil {
		return &command.CommandResult{
			Error:      err.Error(),
			Successful: false,
		}
	}

	if txID == "" {
		return &command.CommandResult{
			Error:      "can't send bond transaction",
			Successful: false,
		}
	}

	log.Info("new bond transaction sent", "txID", txID)

	err = tr.store.AddClaimTransaction(testnetAddr, txID)
	if err != nil {
		log.Panic("unable to add the claim transaction",
			"error", err,
			"discordID", callerID,
			"testnetAddr", testnetAddr,
			"txID", txID,
		)

		return &command.CommandResult{
			Error:      err.Error(),
			Successful: false,
		}
	}

	return &command.CommandResult{
		Successful: true,
		Message:    fmt.Sprintf("Reward claimed successfully✅\nYour claim transaction: https://pacscan.org/transactions/%s", txID),
	}
}

func (tr *TestNetReward) claimStatusHandler(cmd *command.Command, _ command.AppID, _ string, _ ...string) *command.CommandResult {
	cs := tr.store.ClaimStatus()

	result := fmt.Sprintf("Claimed rewards count: %v\nClaimed coins: %v PAC's\nNot-claimed rewards count: %v\nNot-claim coins: %v PAC's\n",
		cs.Claimed, util.ChangeToString(cs.ClaimedAmount), cs.NotClaimed, util.ChangeToString(cs.NotClaimedAmount))

	return &command.CommandResult{
		Successful: true,
		Message:    result,
	}
}
