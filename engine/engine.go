package engine

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/utils"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/libp2p/go-libp2p/core/peer"
)

type BotEngine struct {
	Wallet wallet.IWallet
	Store  store.IStore
	Cm     *client.Mgr
	logger *log.SubLogger

	sync.RWMutex
}

func NewBotEngine(logger *log.SubLogger, cm *client.Mgr, w wallet.IWallet, s store.IStore) (Engine, error) {
	return &BotEngine{
		logger: logger,
		Wallet: w,
		Cm:     cm,
		Store:  s,
	}, nil
}

func (be *BotEngine) NetworkHealth(_ []string) (*NetHealthResponse, error) {
	lastBlockTime, lastBlockHeight := be.Cm.GetLastBlockTime()
	lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0)
	currentTime := time.Now()

	timeDiff := (currentTime.Unix() - int64(lastBlockTime))

	healthStatus := true
	if timeDiff > 15 {
		healthStatus = false
	}

	return &NetHealthResponse{
		HealthStatus:    healthStatus,
		CurrentTime:     currentTime,
		LastBlockTime:   lastBlockTimeFormatted,
		LastBlockHeight: lastBlockHeight,
		TimeDifference:  timeDiff,
	}, nil
}

func (be *BotEngine) NetworkStatus(_ []string) (*NetStatus, error) {
	netInfo, err := be.Cm.GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	chainInfo, err := be.Cm.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}

	return &NetStatus{
		ConnectedPeersCount: netInfo.ConnectedPeersCount,
		ValidatorsCount:     chainInfo.TotalValidators,
		TotalBytesSent:      netInfo.TotalSentBytes,
		TotalBytesReceived:  netInfo.TotalReceivedBytes,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   chainInfo.TotalPower,
		TotalCommitteePower: chainInfo.CommitteePower,
		NetworkName:         netInfo.NetworkName,
	}, nil
}

func (be *BotEngine) MyInfo([]string) (string, error) {
	be.RLock()
	defer be.RUnlock()

	return "not implemented", nil
}

func (be *BotEngine) Withdraw([]string) (string, error) {
	be.Lock()
	defer be.Unlock()

	return "not implemented", nil
}

func (be *BotEngine) NodeInfo(tokens []string) (*NodeInfo, error) {
	if len(tokens) != 1 {
		return nil, errors.New("missing argument: validator address")
	}

	valAddress := tokens[0]

	peerInfo, _, err := be.Cm.GetPeerInfo(valAddress)
	if err != nil {
		return nil, err
	}

	peerID, err := peer.IDFromBytes(peerInfo.PeerId)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.Split(peerInfo.Address, "/")[2], "/")
	ip := parts[0]
	geoData := utils.GetGeoIP(ip)

	val, err := be.Cm.GetValidatorInfo(valAddress)
	if err != nil {
		return nil, err
	}

	return &NodeInfo{
		PeerID:              peerID.String(),
		IPAddress:           peerInfo.Address,
		Agent:               peerInfo.Agent,
		Moniker:             peerInfo.Moniker,
		Country:             geoData.CountryName,
		City:                geoData.City,
		RegionName:          geoData.RegionName,
		TimeZone:            geoData.TimeZone,
		ISP:                 geoData.ISP,
		ValidatorNum:        val.Validator.Number,
		AvailabilityScore:   val.Validator.AvailabilityScore,
		StakeAmount:         val.Validator.Stake,
		LastBondingHeight:   val.Validator.LastBondingHeight,
		LastSortitionHeight: val.Validator.LastSortitionHeight,
	}, nil
}

func (be *BotEngine) Stop() {
	be.logger.Info("shutting bot engine down...")
}

func (be *BotEngine) Start() {
	be.logger.Info("starting the bot engine...")
}
