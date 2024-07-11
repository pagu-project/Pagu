package phoenix

import (
	"fmt"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/engine/command/network"
	"github.com/pagu-project/Pagu/internal/entity"
	utils2 "github.com/pagu-project/Pagu/pkg/utils"
)

//nolint:unused // remove me after I am used
func (pt *Phoenix) nodeInfoHandler(cmd *command.Command,
	_ entity.AppID, _ string, args ...string,
) command.CommandResult {
	valAddress := args[0]

	peerInfo, err := pt.clientMgr.GetPeerInfo(valAddress)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	ip := utils2.ExtractIPFromMultiAddr(peerInfo.Address)
	geoData := utils2.GetGeoIP(pt.ctx, ip)

	nodeInfo := &network.NodeInfo{
		PeerID:     peerInfo.PeerId,
		IPAddress:  peerInfo.Address,
		Agent:      peerInfo.Agent,
		Moniker:    peerInfo.Moniker,
		Country:    geoData.CountryName,
		City:       geoData.City,
		RegionName: geoData.RegionName,
		TimeZone:   geoData.TimeZone,
		ISP:        geoData.ISP,
	}

	// here we check if the node is also a validator.
	// if its a validator , then we populate the validator data.
	// if not validator then we set everything to 0/empty .
	val, err := pt.clientMgr.GetValidatorInfo(valAddress)
	if err == nil && val != nil {
		nodeInfo.ValidatorNum = val.Validator.Number
		nodeInfo.AvailabilityScore = val.Validator.AvailabilityScore
		nodeInfo.StakeAmount = val.Validator.Stake
		nodeInfo.LastBondingHeight = val.Validator.LastBondingHeight
		nodeInfo.LastSortitionHeight = val.Validator.LastSortitionHeight
	} else {
		nodeInfo.ValidatorNum = 0
		nodeInfo.AvailabilityScore = 0
		nodeInfo.StakeAmount = 0
		nodeInfo.LastBondingHeight = 0
		nodeInfo.LastSortitionHeight = 0
	}

	var pip19Score string
	if nodeInfo.AvailabilityScore >= 0.9 {
		pip19Score = fmt.Sprintf("%v✅", nodeInfo.AvailabilityScore)
	} else {
		pip19Score = fmt.Sprintf("%v⚠️", nodeInfo.AvailabilityScore)
	}

	stakeAmountInNanoPAC := nodeInfo.StakeAmount
	stakeAmount := amount.Amount(stakeAmountInNanoPAC)

	// Format the stake amount for display.
	formattedStakeAmount := stakeAmount.Format(amount.UnitPAC)

	return cmd.SuccessfulResult("PeerID: %s\nIP Address: %s\nAgent: %s\n"+
		"Moniker: %s\nCountry: %s\nCity: %s\nRegion Name: %s\nTimeZone: %s\n"+
		"ISP: %s\n\nValidator Info🔍\nNumber: %v\nPIP-19 Score: %s\nStake: %v PAC's\n",
		nodeInfo.PeerID, nodeInfo.IPAddress, nodeInfo.Agent, nodeInfo.Moniker, nodeInfo.Country,
		nodeInfo.City, nodeInfo.RegionName, nodeInfo.TimeZone, nodeInfo.ISP,
		utils2.FormatNumber(int64(nodeInfo.ValidatorNum)),
		pip19Score, formattedStakeAmount)
}
