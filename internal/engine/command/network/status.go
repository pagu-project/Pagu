package network

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	utils2 "github.com/pagu-project/Pagu/pkg/utils"
)

func (n *Network) networkStatusHandler(cmd *command.Command,
	_ entity.AppID, _ string, _ map[string]any,
) command.CommandResult {
	netInfo, err := n.clientMgr.GetNetworkInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	chainInfo, err := n.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	cs, err := n.clientMgr.GetCirculatingSupply()
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
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v PAC\n"+
		"Total Committee Power: %v PAC\nCirculating Supply: %v PAC\n"+
		"\n> Note📝: This info is from one random network node. Non-calculator data may not be consistent.",
		net.NetworkName,
		utils2.FormatNumber(int64(net.ConnectedPeersCount)),
		utils2.FormatNumber(int64(net.ValidatorsCount)),
		utils2.FormatNumber(int64(net.TotalAccounts)),
		utils2.FormatNumber(int64(net.CurrentBlockHeight)),
		utils2.FormatNumber(net.TotalNetworkPower),
		utils2.FormatNumber(net.TotalCommitteePower),
		utils2.FormatNumber(net.CirculatingSupply),
	)
}
