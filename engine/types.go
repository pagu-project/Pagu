package engine

import "time"

type NetHealthResponse struct {
	HealthStatus    bool
	CurrentTime     time.Time
	LastBlockTime   time.Time
	LastBlockHeight uint32
	TimeDifference  int64
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
