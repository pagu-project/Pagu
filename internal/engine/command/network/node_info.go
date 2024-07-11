package network

import (
	"fmt"

	"github.com/pagu-project/Pagu/internal/entity"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	utils2 "github.com/pagu-project/Pagu/pkg/utils"
)

func (n *Network) nodeInfoHandler(cmd command.Command, _ entity.AppID, _ string, args ...string) command.CommandResult {
	valAddress := args[0]

	peerInfo, err := n.clientMgr.GetPeerInfo(valAddress)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	ip := utils2.ExtractIPFromMultiAddr(peerInfo.Address)
	geoData := utils2.GetGeoIP(n.ctx, ip)

	nodeInfo := &NodeInfo{
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
	val, err := n.clientMgr.GetValidatorInfo(valAddress)
	if err == nil && val != nil {
		nodeInfo.ValidatorNum = val.Validator.Number
		nodeInfo.AvailabilityScore = val.Validator.AvailabilityScore
		// Convert NanoPAC to PAC using the Amount type and then to int64.
		stakeAmount := amount.Amount(val.Validator.Stake).ToPAC()
		nodeInfo.StakeAmount = int64(stakeAmount) // Convert float64 to int64.
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

	return cmd.SuccessfulResult("PeerID: %s\nIP Address: %s\nAgent: %s\n"+
		"Moniker: %s\nCountry: %s\nCity: %s\nRegion Name: %s\nTimeZone: %s\n"+
		"ISP: %s\n\nValidator Info🔍\nNumber: %v\nPIP-19 Score: %s\nStake: %v PAC's\n",
		nodeInfo.PeerID, nodeInfo.IPAddress, nodeInfo.Agent, nodeInfo.Moniker, nodeInfo.Country,
		nodeInfo.City, nodeInfo.RegionName, nodeInfo.TimeZone, nodeInfo.ISP,
		utils2.FormatNumber(int64(nodeInfo.ValidatorNum)),
		pip19Score, utils2.FormatNumber(nodeInfo.StakeAmount))
}
