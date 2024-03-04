package engine

import (
	"fmt"
	"time"

	"github.com/kehiy/RoboPac/utils"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util"
)

const (
	NetworkCommandName       = "network"
	NodeInfoCommandName      = "node-info"
	NetworkStatusCommandName = "status"
	NetworkHealthCommandName = "health"
)

func (be *BotEngine) RegisterNetworkCommands() {
	sCmdNodeInfo := Command{
		Name: NodeInfoCommandName,
		Desc: "check the information of a node by providing it's validator address",
		Help: "",
		Args: []Args{
			{
				Name:     "validator-address",
				Desc:     "your validator address",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.nodeInfoHandler,
	}

	sCmdHealth := Command{
		Name:    NetworkHealthCommandName,
		Desc:    "checking network health status",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.networkHealthHandler,
	}

	sCmdStatus := Command{
		Name:    NetworkStatusCommandName,
		Desc:    "network statistics",
		Help:    "",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.networkStatusHandler,
	}

	cmdNetwork := Command{
		Name:        NetworkCommandName,
		Desc:        "network related commands",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&sCmdHealth, &sCmdStatus, &sCmdNodeInfo},
		Handler:     nil,
	}

	be.Cmds = append(be.Cmds, cmdNetwork)
}

func (be *BotEngine) networkHealthHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	lastBlockTime, lastBlockHeight := be.clientMgr.GetLastBlockTime()
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

	return &CommandResult{
		Successful: true,
		Message: fmt.Sprintf("Network is %s\nCurrentTime: %v\nLastBlockTime: %v\nTime Diff: %v\nLast Block Height: %v",
			status, currentTime.Format("02/01/2006, 15:04:05"), lastBlockTimeFormatted, timeDiff, utils.FormatNumber(int64(lastBlockHeight))),
	}, nil
}

func (be *BotEngine) networkStatusHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	netInfo, err := be.clientMgr.GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	chainInfo, err := be.clientMgr.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}

	cs, err := be.clientMgr.GetCirculatingSupply()
	if err != nil {
		cs = 0
	}

	net := NetStatus{
		ValidatorsCount:     chainInfo.TotalValidators,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   chainInfo.TotalPower,
		TotalCommitteePower: chainInfo.CommitteePower,
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   cs,
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

func (be *BotEngine) nodeInfoHandler(_ AppID, _ string, args ...string) (*CommandResult, error) {
	valAddress := args[0]

	peerInfo, err := be.clientMgr.GetPeerInfo(valAddress)
	if err != nil {
		return nil, err
	}

	peerID, err := peer.IDFromBytes(peerInfo.PeerId)
	if err != nil {
		return nil, err
	}

	ip := utils.ExtractIPFromMultiAddr(peerInfo.Address)
	geoData := utils.GetGeoIP(ip)

	nodeInfo := &NodeInfo{
		PeerID:     peerID.String(),
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
	val, err := be.clientMgr.GetValidatorInfo(valAddress)
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
