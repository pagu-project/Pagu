package network

import (
	"context"

	"github.com/pagu-project/Pagu/internal/entity"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/pkg/client"
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

func NewNetwork(ctx context.Context, clientMgr *client.Mgr) Network {
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
	TotalBytesSent      int64
	TotalBytesReceived  int64
	CurrentBlockHeight  uint32
	TotalNetworkPower   int64
	TotalCommitteePower int64
	TotalAccounts       int32
	CirculatingSupply   int64
}

func (n *Network) GetCommand() command.Command {
	subCmdNodeInfo := command.Command{
		Name: NodeInfoCommandName,
		Help: "View the information of a node",
		Args: []command.Args{
			{
				Name:     "validator_address",
				Desc:     "Your validator address",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.nodeInfoHandler,
	}

	subCmdHealth := command.Command{
		Name:        HealthCommandName,
		Help:        "Checking network health status",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.networkHealthHandler,
	}

	subCmdStatus := command.Command{
		Name:        StatusCommandName,
		Help:        "Network statistics",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.networkStatusHandler,
	}

	cmdNetwork := command.Command{
		Name:        CommandName,
		Help:        "Network related commands",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskAll,
	}

	cmdNetwork.AddSubCommand(subCmdHealth)
	cmdNetwork.AddSubCommand(subCmdNodeInfo)
	cmdNetwork.AddSubCommand(subCmdStatus)

	return cmdNetwork
}
