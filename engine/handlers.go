package engine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
)

func (be *BotEngine) ClaimHandler(args []string) (*CommandResult, error) {
	if err := CheckArgs(3, args); err != nil {
		return nil, err
	}

	txHash, err := be.Claim(args[0], args[1], args[2])
	if err != nil {
		return nil, err
	}

	return &CommandResult{
		Successful: true,
		Message:    fmt.Sprintf("Reward claimed successfully✅\nYour claim transaction: https://pacscan.org/transactions/%s", txHash),
	}, nil
}

func (be *BotEngine) ClaimerInfoHandler(args []string) (*CommandResult, error) {
	if err := CheckArgs(1, args); err != nil {
		return nil, err
	}

	claimer, err := be.ClaimerInfo(args[0])
	if err != nil {
		return nil, err
	}

	return &CommandResult{
		Successful: true,
		Message: fmt.Sprintf("TestNet Address: %s\namount: %v PACs\nIsClaimed: %v\n txHash: %s",
			args[0], util.ChangeToString(claimer.TotalReward), claimer.IsClaimed(), claimer.ClaimedTxID),
	}, nil
}

func (be *BotEngine) NetworkHealthHandler(args []string) (*CommandResult, error) {
	health, err := be.NetworkHealth()
	if err != nil {
		return nil, err
	}

	var status string
	if health.HealthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return &CommandResult{
		Successful: true,
		Message: fmt.Sprintf("Network is %s\nCurrentTime: %v\nLastBlockTime: %v\nTime Diff: %v\nLast Block Height: %v",
			status, health.CurrentTime.Format("02/01/2006, 15:04:05"), health.LastBlockTime.Format("02/01/2006, 15:04:05"), health.TimeDifference, utils.FormatNumber(int64(health.LastBlockHeight))),
	}, nil
}

func (be *BotEngine) NodeInfoHandler(args []string) (*CommandResult, error) {
	if err := CheckArgs(1, args); err != nil {
		return nil, err
	}

	nodeInfo, err := be.NodeInfo(args[0])
	if err != nil {
		return nil, err
	}

	var pip19Score string
	if nodeInfo.AvailabilityScore >= 0.9 {
		pip19Score = fmt.Sprintf("%v✅", nodeInfo.AvailabilityScore)
	} else {
		pip19Score = fmt.Sprintf("%v⚠️", nodeInfo.AvailabilityScore)
	}

	result := fmt.Sprintf("PeerID: %s\nIP Address: %s\nAgent: %s\n"+
		"Moniker: %s\nCountry: %s\nCity: %s\nRegion Name: %s\nTimeZone: %s\n"+
		"ISP: %s\n\nValidator Info🔍\nNumber: %v\nPIP-19 Score: %s\nStake: %v PAC's\n",
		nodeInfo.PeerID, nodeInfo.IPAddress, nodeInfo.Agent, nodeInfo.Moniker, nodeInfo.Country,
		nodeInfo.City, nodeInfo.RegionName, nodeInfo.TimeZone, nodeInfo.ISP, utils.FormatNumber(int64(nodeInfo.ValidatorNum)),
		pip19Score, utils.FormatNumber(int64(util.ChangeToCoin(nodeInfo.StakeAmount))))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) NetworkStatusHandler(args []string) (*CommandResult, error) {
	net, err := be.NetworkStatus()
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Network Name: %s\nConnected Peers: %v\n"+
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v PAC\nTotal Committee Power: %v PAC\nCirculating Supply: %v PAC\n"+
		"\n> Note📝: This info is from one random network node. Non-blockchain data may not be consistent.",
		net.NetworkName,
		utils.FormatNumber(int64(net.ConnectedPeersCount)),
		utils.FormatNumber(int64(net.ValidatorsCount)),
		utils.FormatNumber(int64(net.TotalAccounts)),
		utils.FormatNumber(int64(net.CurrentBlockHeight)),
		utils.FormatNumber(int64(util.ChangeToCoin(net.TotalNetworkPower))),
		utils.FormatNumber(int64(util.ChangeToCoin(net.TotalCommitteePower))),
		utils.FormatNumber(int64(util.ChangeToCoin(net.CirculatingSupply))))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) BotWalletHandler(args []string) (*CommandResult, error) {
	addr, blnc := be.BotWallet()
	result := fmt.Sprintf("Address: https://pacscan.org/address/%s\nBalance: %v PAC\n", addr, utils.FormatNumber(int64(util.ChangeToCoin(blnc))))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) ClaimStatusHandler(args []string) (*CommandResult, error) {
	cs := be.ClaimStatus()
	result := fmt.Sprintf("Claimed rewards count: %v\nClaimed coins: %v PAC's\nNot-claimed rewards count: %v\nNot-claim coins: %v PAC's\n",
		cs.Claimed, util.ChangeToString(cs.ClaimedAmount), cs.NotClaimed, util.ChangeToString(cs.NotClaimedAmount))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) RewardCalcHandler(args []string) (*CommandResult, error) {
	if err := CheckArgs(2, args); err != nil {
		return nil, err
	}

	stake, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	reward, time, totalPower, err := be.RewardCalculate(int64(stake), args[1])
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator in one %s ⏰ with %v PAC total power ⚡ of committee."+
		"\n\n> Note📝: This is an estimation and the number can get changed by changes of your stake amount, total power and ...",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), time, utils.FormatNumber(totalPower))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) BoosterPaymentHandler(args []string) (*CommandResult, error) {
	if err := CheckArgs(3, args); err != nil {
		return nil, err
	}

	discordID := args[0]
	twitterName := args[1]
	valAddr := args[2]

	party, err := be.BoosterPayment(discordID, twitterName, valAddr)
	if err != nil {
		return nil, err
	}

	expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)
	msg := fmt.Sprintf("Validator `%s` registered to receive %v stake-PAC coins in total price of $%v."+
		" Visit https://nowpayments.io/payment/?iid=%v to pay it."+
		" The Discount code will expire on %v",
		party.ValAddr, party.AmountInPAC, party.TotalPrice, party.NowPaymentsInvoiceID, expiryDate.Format("2006-01-02"))

	return &CommandResult{
		Successful: true,
		Message:    msg,
	}, nil
}

func (be *BotEngine) BoosterClaimHandler(args []string) (*CommandResult, error) {
	if err := CheckArgs(1, args); err != nil {
		return nil, err
	}

	twitterName := args[0]
	party, err := be.BoosterClaim(twitterName)
	if err != nil {
		return nil, err
	}

	var msg string
	if party.NowPaymentsFinished {
		msg = fmt.Sprintf("Validator `%s` received %v stake-PAC coins."+
			" Transaction: https://pacscan.org/transactions/%v.",
			party.ValAddr, party.AmountInPAC, party.TransactionID)
	} else {
		expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)
		msg = fmt.Sprintf("Validator `%s` registered to receive %v stake-PAC coins in total price of $%v."+
			" Visit https://nowpayments.io/payment/?iid=%v and pay the total amount."+
			" The Discount code will expire on %v",
			party.ValAddr, party.AmountInPAC, party.TotalPrice, party.NowPaymentsInvoiceID, expiryDate.Format("2006-01-02"))
	}

	return &CommandResult{
		Successful: true,
		Message:    msg,
	}, nil
}

func (be *BotEngine) BoosterWhitelistHandler(args []string) (*CommandResult, error) {
	if err := CheckArgs(2, args); err != nil {
		return nil, err
	}

	twitterName := args[0]
	authorizedDiscordID := args[1]
	err := be.BoosterWhitelist(twitterName, authorizedDiscordID)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Twitter `%s` whitelisted", twitterName)

	return &CommandResult{
		Successful: true,
		Message:    msg,
	}, nil
}

func (be *BotEngine) BoosterStatusHandler(args []string) (*CommandResult, error) {
	bs := be.BoosterStatus()

	be.logger.Info("USDT Amount", "amount", bs.Usdt)

	message := fmt.Sprintf("Total Coins: %v PAC\nTotal Packages: %v\nClaimed Packages: %v\nUnClaimed Packages: %v\nPayment Done: %v\nPayment Waiting: %v\nWhite Listed: %v\n",
		bs.Pac, bs.AllPkgs, bs.ClaimedPkgs, bs.UnClaimedPkgs, bs.PaymentDone, bs.PaymentWaiting, bs.Whitelists)

	return &CommandResult{
		Successful: true,
		Message:    message,
	}, nil
}

// default command handler.
func (be *BotEngine) DefaultCommandHandler(args []string) (*CommandResult, error) {
	result := &CommandResult{
		Successful: false,
		Message:    fmt.Sprintf("unknown command: %s", args[0]),
	}

	return result, nil
}
