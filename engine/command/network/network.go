package network

import (
	"context"

	"github.com/pagu-project/Pagu/client"
	"github.com/pagu-project/Pagu/engine/command"
)

const (
	CommandName         = "network"
	NodeInfoCommandName = "node-info"
	StatusCommandName   = "status"
	HealthCommandName   = "health"
	HelpCommandName     = "help"
)

type Network struct {
	ctx       context.Context
	clientMgr *client.Mgr
}

func NewNetwork(ctx context.Context,
	clientMgr *client.Mgr,
) Network {
	return Network{
		ctx:       ctx,
		clientMgr: clientMgr,
	}
}

type NodeInfo struct {
	PeerID              string
	IPAddress           string
	Agent               string
	Moniker             string
	Country             string
	City                string
	RegionName          string
	TimeZone            string
	ISP                 string
	ValidatorNum        int32
	AvailabilityScore   float64
	StakeAmount         int64
	LastBondingHeight   uint32
	LastSortitionHeight uint32
}

type NetStatus struct {
	NetworkName         string
	ConnectedPeersCount uint32
	ValidatorsCount     int32
	TotalBytesSent      uint32
	TotalBytesReceived  uint32
	CurrentBlockHeight  uint32
	TotalNetworkPower   int64
	TotalCommitteePower int64
	TotalAccounts       int32
	CirculatingSupply   int64
}

func (n *Network) GetCommand() command.Command {
	subCmdNodeInfo := command.Command{
		Name: NodeInfoCommandName,
		Desc: "View the information of a node",
		Help: "Provide your validator address on the specific node to get the validator and node info",
		Args: []command.Args{
			{
				Name:     "validator_address",
				Desc:     "Your validator address",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     n.nodeInfoHandler,
	}

	subCmdHealth := command.Command{
		Name:        HealthCommandName,
		Desc:        "Checking network health status",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     n.networkHealthHandler,
	}

	subCmdStatus := command.Command{
		Name:        StatusCommandName,
		Desc:        "Network statistics",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
		Handler:     n.networkStatusHandler,
	}

	cmdNetwork := command.Command{
		Name:        CommandName,
		Desc:        "Network related commands",
		Help:        "",
		Args:        nil,
		AppIDs:      command.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
	}

	cmdNetwork.AddSubCommand(subCmdHealth)
	cmdNetwork.AddSubCommand(subCmdNodeInfo)
	cmdNetwork.AddSubCommand(subCmdStatus)

	return cmdNetwork
}
