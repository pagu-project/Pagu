package engine

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

const (
	CmdClaim                    = "claim"                      //!
	CmdClaimerInfo              = "claimer-info"               //!
	CmdNodeInfo                 = "node-info"                  //!
	CmdNetworkStatus            = "network"                    //!
	CmdNetworkHealth            = "network-health"             //!
	CmdBotWallet                = "wallet"                     //!
	CmdClaimStatus              = "claim-status"               //!
	CmdRewardCalc               = "calc-reward"                //!
	CmdTwitterCampaign          = "twitter-campaign"           //!
	CmdTwitterCampaignStatus    = "twitter-campaign-status"    //!
	CmdTwitterCampaignWhitelist = "twitter-campaign-whitelist" //!
)

// The input is always string.
//
//	The input format is like: [Command] <Arguments ...>
//
// The output is always string, but format might be JSON. ???
func (be *BotEngine) Run(input string) (string, error) {
	logger.Debug("parse command", "input", input)

	cmd, args := be.parseQuery(input)

	switch cmd {
	case CmdClaim:
		if err := CheckArgs(3, args); err != nil {
			return "", err
		}

		txHash, err := be.Claim(args[0], args[1], args[2])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Reward claimed successfully✅\nYour claim transaction: https://pacscan.org/transactions/%s", txHash), nil

	case CmdClaimerInfo:
		if err := CheckArgs(1, args); err != nil {
			return "", err
		}

		claimer, err := be.ClaimerInfo(args[0])
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("TestNet Address: %s\namount: %v PACs\nIsClaimed: %v\n txHash: %s",
			args[0], util.ChangeToString(claimer.TotalReward), claimer.IsClaimed(), claimer.ClaimedTxID), nil

	case CmdNetworkHealth:
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

	case CmdNodeInfo:
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

	case CmdNetworkStatus:
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

	case CmdBotWallet:
		addr, blnc := be.BotWallet()
		return fmt.Sprintf("Address: https://pacscan.org/address/%s\nBalance: %v PAC\n", addr, utils.FormatNumber(int64(util.ChangeToCoin(blnc)))), nil

	case CmdClaimStatus:
		claimed, claimedAmount, notClaimed, notClaimedAmount := be.ClaimStatus()
		return fmt.Sprintf("Claimed rewards count: %v\nClaimed coins: %v PAC's\nNot-claimed rewards count: %v\nNot-claim coins: %v PAC's\n",
			claimed, util.ChangeToString(claimedAmount), notClaimed, util.ChangeToString(notClaimedAmount)), nil

	case CmdRewardCalc:
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

	case CmdTwitterCampaign:
		if err := CheckArgs(3, args); err != nil {
			return "", err
		}

		discordID := args[0]
		twitterName := args[1]
		valAddr := args[2]

		party, err := be.TwitterCampaign(discordID, twitterName, valAddr)
		if err != nil {
			return "", err
		}
		expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)
		msg := fmt.Sprintf("Validator `%s` registered with Discount code `%v`."+
			" Visit https://app.turboswap.io/ to claim your discounted stake-PAC coins."+
			" The Discount code will expire on %v",
			party.ValAddr, party.DiscountCode, expiryDate.Format("2006-01-02"))
		return msg, nil

	case CmdTwitterCampaignStatus:
		if err := CheckArgs(1, args); err != nil {
			return "", err
		}

		twitterName := args[0]
		party, _, err := be.TwitterCampaignStatus(twitterName)
		if err != nil {
			return "", err
		}

		// check the status (tOdO)
		expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)
		msg := fmt.Sprintf("Validator `%s` registered with Discount code `%v`."+
			" Visit https://app.turboswap.io/ to claim your discounted stake-PAC coins."+
			" The Discount code will expire on %v",
			party.ValAddr, party.DiscountCode, expiryDate.Format("2006-01-02"))
		return msg, nil

	case CmdTwitterCampaignWhitelist:
		if err := CheckArgs(2, args); err != nil {
			return "", err
		}

		twitterName := args[0]
		authorizedDiscordID := args[1]
		err := be.TwitterCampaignWhitelist(twitterName, authorizedDiscordID)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf("Twitter `%s` whitelisted", twitterName)
		return msg, nil

	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

func (be *BotEngine) parseQuery(query string) (string, []string) {
	subs := strings.Split(query, " ")
	if len(subs) == 0 {
		return "", nil
	}

	return subs[0], subs[1:]
}
