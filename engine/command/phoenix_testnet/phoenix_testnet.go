package phoenixtestnet

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util"
	"github.com/robopac-project/RoboPac/client"
	"github.com/robopac-project/RoboPac/database"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/engine/command/network"
	"github.com/robopac-project/RoboPac/utils"
	"github.com/robopac-project/RoboPac/wallet"
)

const (
	PhoenixTestnetCommandName  = "phoenix"
	PhoenixFaucetCommandName   = "faucet"
	PhoenixWalletCommandName   = "wallet"
	PhoenixStatusCommandName   = "status"
	PhoenixHealthCommandName   = "health"
	PhoenixNodeInfoCommandName = "node-info"
	PhoenixHelpCommandName     = "help"
)

type PhoenixTestnet struct {
	wallet    wallet.IWallet
	db        database.DB
	clientMgr *client.Mgr
}

func NewPhoenixTestnet(wallet wallet.IWallet,
	clientMgr *client.Mgr, db database.DB,
) PhoenixTestnet {
	return PhoenixTestnet{
		wallet:    wallet,
		clientMgr: clientMgr,
		db:        db,
	}
}

func (pt *PhoenixTestnet) GetCommand() command.Command {
	subCmdFaucet := command.Command{
		Name: PhoenixFaucetCommandName,
		Desc: "Get 5 tPAC Coins on Phoenix Testnet for Testing your code or project",
		Help: "There is a limit that you can only get faucets 1 time per day with each user ID and address",
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "your testnet address, example: tpc1z....",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.faucetHandler,
	}

	subCmdWallet := command.Command{
		Name:        PhoenixWalletCommandName,
		Desc:        "Check the status of RoboPac faucet address wallet on Phoenix network",
		Help:        "",
		Args:        nil,
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.walletHandler,
	}

	subCmdHealth := command.Command{
		Name:        PhoenixHealthCommandName,
		Desc:        "Checking Phoenix test-network health status",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord, command.AppIdgRPC, command.AppIdTelegram},
		Handler:     pt.networkHealthHandler,
	}

	subCmdStatus := command.Command{
		Name:        PhoenixStatusCommandName,
		Desc:        "Phoenix test-network statistics",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.networkStatusHandler,
	}

	subCmdNodeInfo := command.Command{
		Name: PhoenixNodeInfoCommandName,
		Desc: "View the information of a node running on Phoenix test-network",
		Help: "Provide your validator address on the specific node to get the validator and node info (Phoenix network)",
		Args: []command.Args{
			{
				Name:     "validator-address",
				Desc:     "Your validator address start with tpc1p...",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     pt.nodeInfoHandler,
	}

	cmdPhoenixTestnet := command.Command{
		Name:        PhoenixTestnetCommandName,
		Desc:        "Phoenix Testnet tools and utils for developers",
		Help:        "",
		Args:        nil,
		AppIDs:      command.AllAppIDs(),
		SubCommands: make([]command.Command, 2),
		Handler:     nil,
	}

	cmdPhoenixTestnet.AddSubCommand(subCmdFaucet)
	cmdPhoenixTestnet.AddSubCommand(subCmdWallet)
	cmdPhoenixTestnet.AddSubCommand(subCmdHealth)
	cmdPhoenixTestnet.AddSubCommand(subCmdStatus)
	cmdPhoenixTestnet.AddSubCommand(subCmdNodeInfo)

	return cmdPhoenixTestnet
}

func (pt *PhoenixTestnet) faucetHandler(cmd command.Command, _ command.AppID, callerID string, args ...string) command.CommandResult {
	if !pt.db.HasUser(callerID) {
		if err := pt.db.AddUser(
			&database.User{
				ID: callerID,
			},
		); err != nil {
			return cmd.ErrorResult(err)
		}
	}

	if !pt.db.CanGetFaucet(callerID) {
		return cmd.FailedResult("Uh, you used your share of faucets today!")
	}

	if pt.wallet.Balance() < 5 {
		return cmd.FailedResult("RoboPac Phoenix wallet is empty, please contact the team!")
	}

	toAddr := args[0]
	txID, err := pt.wallet.TransferTransaction(toAddr, "Phoenix Testnet RoboPac Faucet", 5) //! define me on config?
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = pt.db.AddFaucet(&database.Faucet{
		Address:         toAddr,
		Amount:          5,
		TransactionHash: txID,
		UserID:          callerID,
	}); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("You got %d tPAC in %s address on Phoenix Testnet!", 5, toAddr)
}

func (pt *PhoenixTestnet) walletHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	return cmd.SuccessfulResult("RoboPac Phoenix Address: %s\nBalance: %d", pt.wallet.Address(), pt.wallet.Balance())
}

func (pt *PhoenixTestnet) networkHealthHandler(cmd command.Command, _ command.AppID, _ string, _ ...string) command.CommandResult {
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

func (pt *PhoenixTestnet) networkStatusHandler(cmd command.Command, _ command.AppID, _ string, _ ...string) command.CommandResult {
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

	net := network.NetStatus{
		ValidatorsCount:     chainInfo.TotalValidators,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   chainInfo.TotalPower,
		TotalCommitteePower: chainInfo.CommitteePower,
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   cs,
	}

	return cmd.SuccessfulResult("Network Name: %s\nConnected Peers: %v\n"+
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
}

func (pt *PhoenixTestnet) nodeInfoHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	valAddress := args[0]

	peerInfo, err := pt.clientMgr.GetPeerInfo(valAddress)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	peerID, err := peer.IDFromBytes(peerInfo.PeerId)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	ip := utils.ExtractIPFromMultiAddr(peerInfo.Address)
	geoData := utils.GetGeoIP(ip)

	nodeInfo := &network.NodeInfo{
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

	return cmd.SuccessfulResult("PeerID: %s\nIP Address: %s\nAgent: %s\n"+
		"Moniker: %s\nCountry: %s\nCity: %s\nRegion Name: %s\nTimeZone: %s\n"+
		"ISP: %s\n\nValidator Info🔍\nNumber: %v\nPIP-19 Score: %s\nStake: %v PAC's\n",
		nodeInfo.PeerID, nodeInfo.IPAddress, nodeInfo.Agent, nodeInfo.Moniker, nodeInfo.Country,
		nodeInfo.City, nodeInfo.RegionName, nodeInfo.TimeZone, nodeInfo.ISP, utils.FormatNumber(int64(nodeInfo.ValidatorNum)),
		pip19Score, utils.FormatNumber(int64(util.ChangeToCoin(nodeInfo.StakeAmount))))
}
