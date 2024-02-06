package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pactus-project/pactus/util"
)

const (
	CmdClaim         = "claim"          //!
	CmdClaimerInfo   = "claimer-info"   //!
	CmdNodeInfo      = "node-info"      //!
	CmdNetworkStatus = "network"        //!
	CmdNetworkHealth = "network-health" //!
	CmdBotWallet     = "wallet"         //!
	CmdClaimStatus   = "claim-status"   //!
	CmdRewardCalc    = "calc-reward"    //!
)

// The input is always string.
//
//	The input format is like: [Command] <Arguments ...>
//
// The output is always string, but format might be JSON. ???
func (be *BotEngine) Run(input string) (string, error) {
	cmd, args := be.parseQuery(input)

	switch cmd {
	case CmdClaim:
		if len(args) != 3 {
			return "", fmt.Errorf("expected to have 3 arguments, but it received %d", len(args))
		}

		txHash, err := be.Claim(args[0], args[1], args[2])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Reward claimed successfully✅\nYour claim transaction: https://pacscan.org/transactions/%s", txHash), nil

	case CmdClaimerInfo:
		if len(args) != 1 {
			return "", fmt.Errorf("expected to have 1 arguments, but it received %d", len(args))
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
			status, health.CurrentTime.Format("02/01/2006, 15:04:05"), health.LastBlockTime.Format("02/01/2006, 15:04:05"), health.TimeDifference, health.LastBlockHeight), nil

	case CmdNodeInfo:
		if len(args) != 1 {
			return "", fmt.Errorf("expected to have 1 arguments, but it received %d", len(args))
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
			"ISP: %s\n\nValidator Info🔍\nNumber: %v\nPIP19-Score: %s\nStake: %v PAC's\n",
			nodeInfo.PeerID, nodeInfo.IPAddress, nodeInfo.Agent, nodeInfo.Moniker, nodeInfo.Country,
			nodeInfo.City, nodeInfo.RegionName, nodeInfo.TimeZone, nodeInfo.ISP, nodeInfo.ValidatorNum,
			pip19Score, util.ChangeToString(nodeInfo.StakeAmount)), nil

	case CmdNetworkStatus:
		net, err := be.NetworkStatus()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Network Name: %s\nConnected Peers: %v\n"+
			"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v PAC\nTotal Committee Power: %v PAC\nCirculating Supply: %v PAC\n"+
			"\n> Note📝: This info is from one random network node. Non-blockchain data may not be consistent.",
			net.NetworkName, net.ConnectedPeersCount, net.ValidatorsCount, net.TotalAccounts, net.CurrentBlockHeight, util.ChangeToString(net.TotalNetworkPower),
			util.ChangeToString(net.TotalCommitteePower), util.ChangeToString(net.CirculatingSupply)), nil

	case CmdBotWallet:
		addr, blnc := be.BotWallet()
		return fmt.Sprintf("Address: https://pacscan.org/address/%s\nBalance: %v\n", addr, util.ChangeToString(blnc)), nil

	case CmdClaimStatus:
		claimed, claimedAmount, notClaimed, notClaimedAmount := be.ClaimStatus()
		return fmt.Sprintf("Claimed rewards count: %v\nClaimed coins: %v PAC's\nNot-claimed rewards count: %v\nNot-claim coins: %v PAC's\n",
			claimed, util.ChangeToString(claimedAmount), notClaimed, util.ChangeToString(notClaimedAmount)), nil

	case CmdRewardCalc:
		if len(args) != 2 {
			return "", fmt.Errorf("expected to have 2 arguments, but it received %d", len(args))
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
			reward, stake, time, totalPower), nil

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
