package phoenix

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/engine/command/network"
	"github.com/pagu-project/Pagu/internal/entity"
	utils2 "github.com/pagu-project/Pagu/pkg/utils"
)

func (pt *Phoenix) networkStatusHandler(cmd *command.Command,
	_ entity.AppID, _ string, _ map[string]any,
) command.CommandResult {
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
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v\n"+
		"Total Committee Power: %v\nCirculating Supply: %v\n"+
		"\n> Note📝: This info is from one random network node. Non-calculator data may not be consistent.",
		net.NetworkName,
		utils2.FormatNumber(int64(net.ConnectedPeersCount)),
		utils2.FormatNumber(int64(net.ValidatorsCount)),
		utils2.FormatNumber(int64(net.TotalAccounts)),
		utils2.FormatNumber(int64(net.CurrentBlockHeight)),
		net.TotalNetworkPower,
		net.TotalCommitteePower,
		net.CirculatingSupply)
}
