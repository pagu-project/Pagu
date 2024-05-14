package phoenix

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/engine/command"
	"github.com/pagu-project/Pagu/engine/command/network"
	"github.com/pagu-project/Pagu/utils"
	"time"
)

func (pt *Phoenix) networkHealthHandler(cmd command.Command, _ command.AppID, _ string, _ ...string) command.CommandResult {
	lastBlockTime, lastBlockHeight := pt.clientMgr.GetLastBlockTime()
	lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0).Format("02/01/2006, 15:04:05")
	currentTime := time.Now()

	timeDiff := (currentTime.Unix() - int64(lastBlockTime))

	healthStatus := true
	if timeDiff > 15 {
		healthStatus = false
	}

	var status string
	if healthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return cmd.SuccessfulResult("Network is %s\nCurrentTime: %v\nLastBlockTime: %v\nTime Diff: %v\nLast Block Height: %v",
		status, currentTime.Format("02/01/2006, 15:04:05"), lastBlockTimeFormatted, timeDiff, utils.FormatNumber(int64(lastBlockHeight)))
}

func (pt *Phoenix) networkStatusHandler(cmd command.Command, _ command.AppID, _ string, _ ...string) command.CommandResult {
	netInfo, err := pt.clientMgr.GetNetworkInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	chainInfo, err := pt.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	cs, err := pt.clientMgr.GetCirculatingSupply()
	if err != nil {
		cs = 0
	}

	// Convert int64 to float64.
	totalNetworkPower, err := amount.NewAmount(float64(chainInfo.TotalPower))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	totalCommitteePower, err := amount.NewAmount(float64(chainInfo.CommitteePower))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	circulatingSupply, err := amount.NewAmount(float64(cs))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	// Convert Amount back to int64 for struct literal.
	net := network.NetStatus{
		ValidatorsCount:     chainInfo.TotalValidators,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   int64(totalNetworkPower.ToPAC()),
		TotalCommitteePower: int64(totalCommitteePower.ToPAC()),
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   int64(circulatingSupply.ToPAC()),
	}

	return cmd.SuccessfulResult("Network Name: %s\nConnected Peers: %v\n"+
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v\nTotal Committee Power: %v\nCirculating Supply: %v\n"+
		"\n> Note📝: This info is from one random network node. Non-blockchain data may not be consistent.",
		net.NetworkName,
		utils.FormatNumber(int64(net.ConnectedPeersCount)),
		utils.FormatNumber(int64(net.ValidatorsCount)),
		utils.FormatNumber(int64(net.TotalAccounts)),
		utils.FormatNumber(int64(net.CurrentBlockHeight)),
		net.TotalNetworkPower,
		net.TotalCommitteePower,
		net.CirculatingSupply)
}
