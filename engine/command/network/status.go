package network

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/engine/command"
	"github.com/pagu-project/Pagu/utils"
)

func (be *Network) networkStatusHandler(cmd command.Command, _ command.AppID, _ string, _ ...string) command.CommandResult {
	netInfo, err := be.clientMgr.GetNetworkInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	chainInfo, err := be.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	cs, err := be.clientMgr.GetCirculatingSupply()
	if err != nil {
		cs = 0
	}

	// Convert NanoPAC to PAC using the Amount type.
	totalNetworkPower := amount.Amount(chainInfo.TotalPower).ToPAC()
	totalCommitteePower := amount.Amount(chainInfo.CommitteePower).ToPAC()
	circulatingSupply := amount.Amount(cs).ToPAC()

	net := NetStatus{
		ValidatorsCount:     chainInfo.TotalValidators,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   int64(totalNetworkPower),
		TotalCommitteePower: int64(totalCommitteePower),
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   int64(circulatingSupply),
	}

	return cmd.SuccessfulResult("Network Name: %s\nConnected Peers: %v\n"+
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v PAC\nTotal Committee Power: %v PAC\nCirculating Supply: %v PAC\n"+
		"\n> Note📝: This info is from one random network node. Non-blockchain data may not be consistent.",
		net.NetworkName,
		utils.FormatNumber(int64(net.ConnectedPeersCount)),
		utils.FormatNumber(int64(net.ValidatorsCount)),
		utils.FormatNumber(int64(net.TotalAccounts)),
		utils.FormatNumber(int64(net.CurrentBlockHeight)),
		utils.FormatNumber(net.TotalNetworkPower),
		utils.FormatNumber(net.TotalCommitteePower),
		utils.FormatNumber(net.CirculatingSupply),
	)
}
