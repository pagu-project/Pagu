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
	putils "github.com/pactus-project/pactus/util"
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

	localClient, err := client.NewClient(cfg.LocalNode)
	if err != nil {
		log.Error("can't make a new local client", "err", err, "addr", cfg.LocalNode)
		return nil, err
	}

	cm.AddClient(localClient)

	for _, nn := range cfg.NetworkNodes {
		c, err := client.NewClient(nn)
		if err != nil {
			log.Error("can't add new network node client", "err", err, "addr", nn)
		}
		cm.AddClient(c)
	}

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

	cs, err := be.Cm.GetCirculatingSupply()
	if err != nil {
		cs = 0
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
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   cs,
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

	valInfo, _ := be.Cm.GetValidatorInfo(mainnetAddr)
	if valInfo != nil {
		return "", errors.New("this address is already a staked validator")
	}

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

func (be *BotEngine) BotWallet() (string, int64) {
	return be.Wallet.Address(), be.Wallet.Balance()
}

func (be *BotEngine) ClaimStatus() (int64, int64, int64, int64) {
	return be.Store.Status()
}

func (be *BotEngine) RewardCalculate(stake int64, t string) (int64, string, int64, error) {
	if stake < 1 || stake > 1_000 {
		return 0, "", 0, errors.New("minimum of stake is 1 PAC and maximum is 1,000 PAC")
	}

	var blocks int64
	time := t
	switch t {
	case "day":
		blocks = 8640
	case "month":
		blocks = 259200
	case "year":
		blocks = 3110400
	default:
		blocks = 8640
		time = "day"
	}

	bi, err := be.Cm.GetBlockchainInfo()
	if err != nil {
		return 0, "", 0, nil
	}

	reward := (stake * int64(blocks)) / int64(putils.ChangeToCoin(bi.TotalPower))
	return reward, time, utils.AtomicToCoin(bi.TotalPower), nil
}

func (be *BotEngine) Stop() {
	be.logger.Info("shutting bot engine down...")

	be.Cm.Stop()
}

func (be *BotEngine) Start() {
	be.logger.Info("starting the bot engine...")
}
