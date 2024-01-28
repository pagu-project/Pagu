package engine

import (
	"errors"
	"sync"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/config"
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

func NewBotEngine(cfg *config.Config) (IEngine, error) {
	cm := client.NewClientMgr()
	c, err := client.NewClient(cfg.LocalNode)
	if err != nil {
		log.Error("can't make a new local-net client", "err", err, "addr", cfg.LocalNode)
		return nil, err
	}

	cm.AddClient("local-net", c)

	// initializing logger global instance.
	log.InitGlobalLogger()

	// new subLogger for engine.
	eSl := log.NewSubLogger("engine")

	// new subLogger for store.
	sSl := log.NewSubLogger("store")

	// new subLogger for store.
	wSl := log.NewSubLogger("wallet")

	// load or create wallet.
	wallet := wallet.Open(cfg, wSl)
	if wallet == nil {
		log.Panic("wallet could not be opened, wallet is nil", "path", cfg.WalletPath)
	}

	log.Info("wallet opened successfully", "address", wallet.Address())

	// load store.
	store, err := store.NewStore(cfg, sSl)
	if err != nil {
		log.Panic("could not load store", "err", err, "path", cfg.StorePath)
	}

	log.Info("store loaded successfully", "path", cfg.StorePath)

	return newBotEngine(eSl, cm, wallet, store), nil
}

func newBotEngine(logger *log.SubLogger, cm *client.Mgr, w wallet.IWallet, s store.IStore) *BotEngine {
	return &BotEngine{
		logger: logger,
		Wallet: w,
		Cm:     cm,
		Store:  s,
	}
}

func (be *BotEngine) NetworkHealth() (*NetHealthResponse, error) {
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

func (be *BotEngine) NetworkStatus() (*NetStatus, error) {
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

func (be *BotEngine) NodeInfo(valAddress string) (*NodeInfo, error) {
	peerInfo, err := be.Cm.GetPeerInfo(valAddress)
	if err != nil {
		return nil, err
	}

	peerID, err := peer.IDFromBytes(peerInfo.PeerId)
	if err != nil {
		return nil, err
	}

	ip := utils.ExtractIPFromMultiAddr(peerInfo.Address)
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

func (be *BotEngine) ClaimerInfo(testNetValAddr string) (*store.Claimer, error) {
	be.RLock()
	defer be.RUnlock()

	claimer := be.Store.ClaimerInfo(testNetValAddr)
	if claimer == nil {
		return nil, errors.New("not found")
	}

	return claimer, nil
}

func (be *BotEngine) Claim(discordID string, testnetAddr string, mainnetAddr string) (string, error) {
	be.Lock()
	defer be.Unlock()

	be.logger.Info("new claim request", "mainnetAddr", mainnetAddr, "testnetAddr", testnetAddr, "discordID", discordID)

	if utils.AtomicToCoin(be.Wallet.Balance()) <= 500 {
		be.logger.Warn("bot wallet hasn't enough balance")
		return "", errors.New("insufficient wallet balance")
	}

	claimer := be.Store.ClaimerInfo(testnetAddr)
	if claimer == nil {
		return "", errors.New("claimer not found")
	}

	if claimer.DiscordID != discordID {
		be.logger.Warn("try to claim other's reward", "claimer", claimer.DiscordID, "discordID", discordID)
		return "", errors.New("invalid claimer")
	}

	if claimer.IsClaimed() {
		return "", errors.New("this claimer have already claimed rewards")
	}

	peerInfo, err := be.Cm.GetPeerInfoFirstVal(mainnetAddr)
	if err != nil {
		return "", err
	}

	memo := "TestNet reward claim from RoboPac"
	txID, err := be.Wallet.BondTransaction(peerInfo.ConsensusKeys[0], mainnetAddr, memo, claimer.TotalReward)
	if err != nil {
		return "", err
	}

	if txID == "" {
		return "", errors.New("can't send bond transaction")
	}

	be.logger.Info("new bond transaction sent", "txID", txID)

	err = be.Store.AddClaimTransaction(testnetAddr, txID)
	if err != nil {
		be.logger.Panic("unable to add the claim transaction",
			"error", err,
			"discordID", discordID,
			"testnetAddr", testnetAddr,
			"txID", txID,
		)

		return "", err
	}

	return txID, nil
}

func (be *BotEngine) Stop() {
	be.logger.Info("shutting bot engine down...")

	be.Cm.Stop()
}

func (be *BotEngine) Start() {
	be.logger.Info("starting the bot engine...")
}
