package engine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
)

func ClaimHandler(be *BotEngine, args []string) (string, error) {
	if err := CheckArgs(3, args); err != nil {
		return "", err
	}

	txHash, err := be.Claim(args[0], args[1], args[2])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Reward claimed successfully✅\nYour claim transaction: https://pacscan.org/transactions/%s", txHash), nil
}

func ClaimerInfoHandler(be *BotEngine, args []string) (string, error) {
	if err := CheckArgs(1, args); err != nil {
		return "", err
	}

	claimer, err := be.ClaimerInfo(args[0])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("TestNet Address: %s\namount: %v PACs\nIsClaimed: %v\n txHash: %s",
		args[0], util.ChangeToString(claimer.TotalReward), claimer.IsClaimed(), claimer.ClaimedTxID), nil
}

func NetworkHealthHandler(be *BotEngine, args []string) (string, error) {
	health, err := be.NetworkHealth()
	if err != nil {
		return "", err
	}

	var status string
	if health.HealthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return fmt.Sprintf("Network is %s\nCurrentTime: %v\nLastBlockTime: %v\nTime Diff: %v\nLast Block Height: %v",
		status, health.CurrentTime.Format("02/01/2006, 15:04:05"), health.LastBlockTime.Format("02/01/2006, 15:04:05"), health.TimeDifference, utils.FormatNumber(int64(health.LastBlockHeight))), nil
}

func NodeInfoHandler(be *BotEngine, args []string) (string, error) {
	if err := CheckArgs(1, args); err != nil {
		return "", err
	}

	nodeInfo, err := be.NodeInfo(args[0])
	if err != nil {
		return "", err
	}

	var pip19Score string
	if nodeInfo.AvailabilityScore >= 0.9 {
		pip19Score = fmt.Sprintf("%v✅", nodeInfo.AvailabilityScore)
	} else {
		pip19Score = fmt.Sprintf("%v⚠️", nodeInfo.AvailabilityScore)
	}

	return fmt.Sprintf("PeerID: %s\nIP Address: %s\nAgent: %s\n"+
		"Moniker: %s\nCountry: %s\nCity: %s\nRegion Name: %s\nTimeZone: %s\n"+
		"ISP: %s\n\nValidator Info🔍\nNumber: %v\nPIP-19 Score: %s\nStake: %v PAC's\n",
		nodeInfo.PeerID, nodeInfo.IPAddress, nodeInfo.Agent, nodeInfo.Moniker, nodeInfo.Country,
		nodeInfo.City, nodeInfo.RegionName, nodeInfo.TimeZone, nodeInfo.ISP, utils.FormatNumber(int64(nodeInfo.ValidatorNum)),
		pip19Score, utils.FormatNumber(int64(util.ChangeToCoin(nodeInfo.StakeAmount)))), nil
}

func NetworkStatusHandler(be *BotEngine, args []string) (string, error) {
	net, err := be.NetworkStatus()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Network Name: %s\nConnected Peers: %v\n"+
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v PAC\nTotal Committee Power: %v PAC\nCirculating Supply: %v PAC\n"+
		"\n> Note📝: This info is from one random network node. Non-blockchain data may not be consistent.",
		net.NetworkName,
		utils.FormatNumber(int64(net.ConnectedPeersCount)),
		utils.FormatNumber(int64(net.ValidatorsCount)),
		utils.FormatNumber(int64(net.TotalAccounts)),
		utils.FormatNumber(int64(net.CurrentBlockHeight)),
		utils.FormatNumber(int64(util.ChangeToCoin(net.TotalNetworkPower))),
		utils.FormatNumber(int64(util.ChangeToCoin(net.TotalCommitteePower))),
		utils.FormatNumber(int64(util.ChangeToCoin(net.CirculatingSupply))),
	), nil
}

func BotWalletHandler(be *BotEngine, args []string) (string, error) {
	addr, blnc := be.BotWallet()
	return fmt.Sprintf("Address: https://pacscan.org/address/%s\nBalance: %v PAC\n", addr, utils.FormatNumber(int64(util.ChangeToCoin(blnc)))), nil
}

func ClaimStatusHandler(be *BotEngine, args []string) (string, error) {
	cs := be.ClaimStatus()
	return fmt.Sprintf("Claimed rewards count: %v\nClaimed coins: %v PAC's\nNot-claimed rewards count: %v\nNot-claim coins: %v PAC's\n",
		cs.Claimed, util.ChangeToString(cs.ClaimedAmount), cs.NotClaimed, util.ChangeToString(cs.NotClaimedAmount)), nil
}

func RewardCalcHandler(be *BotEngine, args []string) (string, error) {
	if err := CheckArgs(2, args); err != nil {
		return "", err
	}

	stake, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}

	reward, time, totalPower, err := be.RewardCalculate(int64(stake), args[1])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator in one %s ⏰ with %v PAC total power ⚡ of committee."+
		"\n\n> Note📝: This is an estimation and the number can get changed by changes of your stake amount, total power and ...",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), time, utils.FormatNumber(totalPower)), nil
}

func BoosterPaymentHandler(be *BotEngine, args []string) (string, error) {
	if err := CheckArgs(3, args); err != nil {
		return "", err
	}

	discordID := args[0]
	twitterName := args[1]
	valAddr := args[2]

	party, err := be.BoosterPayment(discordID, twitterName, valAddr)
	if err != nil {
		return "", err
	}
	expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)
	msg := fmt.Sprintf("Validator `%s` registered to receive %v stake-PAC coins in total price of $%v."+
		" Visit https://nowpayments.io/payment/?iid=%v to pay it."+
		" The Discount code will expire on %v",
		party.ValAddr, party.AmountInPAC, party.TotalPrice, party.NowPaymentsInvoiceID, expiryDate.Format("2006-01-02"))
	return msg, nil
}

func BoosterClaimHandler(be *BotEngine, args []string) (string, error) {
	if err := CheckArgs(1, args); err != nil {
		return "", err
	}

	twitterName := args[0]
	party, err := be.BoosterClaim(twitterName)
	if err != nil {
		return "", err
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
	return msg, nil
}

func BoosterWhitelistHandler(be *BotEngine, args []string) (string, error) {
	if err := CheckArgs(2, args); err != nil {
		return "", err
	}

	twitterName := args[0]
	authorizedDiscordID := args[1]
	err := be.BoosterWhitelist(twitterName, authorizedDiscordID)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("Twitter `%s` whitelisted", twitterName)
	return msg, nil
}

func BoosterStatusHandler(be *BotEngine, args []string) (string, error) {
	bs := be.BoosterStatus()

	be.logger.Info("USDT Amount", "amount", bs.Usdt)

	return fmt.Sprintf("Total Coins: %v PAC\nTotal Packages: %v\nClaimed Packages: %v\nUnClaimed Packages: %v\nPayment Done: %v\nPayment Waiting: %v\nWhite Listed: %v\n",
		bs.Pac, bs.AllPkgs, bs.ClaimedPkgs, bs.UnClaimedPkgs, bs.PaymentDone, bs.PaymentWaiting, bs.Whitelists), nil
}

// default command handler.
func DefaultCommandHandler(be *BotEngine, args []string) (string, error) {
	return "", fmt.Errorf("unknown command: %s", args[0])
}
