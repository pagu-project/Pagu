package engine

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/turboswap"
	"github.com/kehiy/RoboPac/twitter_api"
	"github.com/kehiy/RoboPac/utils"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/libp2p/go-libp2p/core/peer"
	gonanoid "github.com/matoous/go-nanoid/v2"
	putils "github.com/pactus-project/pactus/util"
)

type BotEngine struct {
	ctx    context.Context //nolint
	cancel func()

	wallet    wallet.IWallet
	store     store.IStore
	turboswap turboswap.ITurboSwap
	clientMgr *client.Mgr
	logger    *log.SubLogger

	twitterClient twitter_api.IClient

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
	store, err := store.NewStore(cfg.StorePath, sSl)
	if err != nil {
		log.Panic("could not load store", "err", err)
	}
	log.Info("store loaded successfully", "path", cfg.StorePath)

	twitterClient, err := twitter_api.NewClient(cfg.TwitterAPICfg.BearerToken, cfg.TwitterAPICfg.TwitterID)
	if err != nil {
		log.Panic("could not start twitter client", "err", err)
	}
	log.Info("twitterClient loaded successfully")

	turboswap, err := turboswap.NewTurboswap(cfg.TurboswapConfig.APIToken, cfg.TurboswapConfig.URL)
	if err != nil {
		log.Error("could not start twitter client", "err", err)
	}
	log.Info("turboswap loaded successfully")

	return newBotEngine(eSl, cm, wallet, store, twitterClient, turboswap), nil
}

func newBotEngine(logger *log.SubLogger, cm *client.Mgr, w wallet.IWallet, s store.IStore,
	twitterClient twitter_api.IClient, turboswap turboswap.ITurboSwap,
) *BotEngine {
	ctx, cancel := context.WithCancel(context.Background())

	return &BotEngine{
		ctx:           ctx,
		cancel:        cancel,
		logger:        logger,
		wallet:        w,
		clientMgr:     cm,
		store:         s,
		twitterClient: twitterClient,
		turboswap:     turboswap,
	}
}

func (be *BotEngine) NetworkHealth() (*NetHealthResponse, error) {
	lastBlockTime, lastBlockHeight := be.clientMgr.GetLastBlockTime()
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

	val, err := be.clientMgr.GetValidatorInfo(valAddress)
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

	claimer := be.store.ClaimerInfo(testNetValAddr)
	if claimer == nil {
		return nil, errors.New("not found")
	}

	return claimer, nil
}

func (be *BotEngine) Claim(discordID string, testnetAddr string, mainnetAddr string) (string, error) {
	be.Lock()
	defer be.Unlock()

	be.logger.Info("new claim request", "mainnetAddr", mainnetAddr, "testnetAddr", testnetAddr, "discordID", discordID)

	valInfo, _ := be.clientMgr.GetValidatorInfo(mainnetAddr)
	if valInfo != nil {
		return "", errors.New("this address is already a staked validator")
	}

	if utils.ChangeToCoin(be.wallet.Balance()) <= 500 {
		be.logger.Warn("bot wallet hasn't enough balance")
		return "", errors.New("insufficient wallet balance")
	}

	claimer := be.store.ClaimerInfo(testnetAddr)
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

	pubKey, err := be.clientMgr.FindPublicKey(mainnetAddr, true)
	if err != nil {
		return "", err
	}

	memo := "TestNet reward claim from RoboPac"
	txID, err := be.wallet.BondTransaction(pubKey, mainnetAddr, memo, claimer.TotalReward)
	if err != nil {
		return "", err
	}

	if txID == "" {
		return "", errors.New("can't send bond transaction")
	}

	be.logger.Info("new bond transaction sent", "txID", txID)

	err = be.store.AddClaimTransaction(testnetAddr, txID)
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
	return be.wallet.Address(), be.wallet.Balance()
}

func (be *BotEngine) ClaimStatus() (int64, int64, int64, int64) {
	return be.store.ClaimStatus()
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

	bi, err := be.clientMgr.GetBlockchainInfo()
	if err != nil {
		return 0, "", 0, nil
	}

	reward := (stake * int64(blocks)) / int64(putils.ChangeToCoin(bi.TotalPower))
	return reward, time, int64(utils.ChangeToCoin(bi.TotalPower)), nil
}

func (be *BotEngine) TwitterCampaign(discordID, twitterName, valAddr string) (*store.TwitterParty, error) {
	be.Lock()
	defer be.Unlock()

	existingParty := be.store.FindTwitterParty(twitterName)
	if existingParty != nil {
		return existingParty, nil
	}

	valInfo, _ := be.clientMgr.GetValidatorInfo(valAddr)
	if valInfo != nil {
		return nil, errors.New("this address is already a staked validator")
	}

	pubKey, err := be.clientMgr.FindPublicKey(valAddr, false)
	if err != nil {
		return nil, err
	}

	userInfo, err := be.twitterClient.UserInfo(be.ctx, twitterName)
	if err != nil {
		return nil, err
	}
	if !userInfo.IsVerified {
		if !be.store.IsWhitelisted(userInfo.TwitterID) {
			threeYearsAgo := time.Now().AddDate(-3, 0, 0)
			if userInfo.CreatedAt.After(threeYearsAgo) {
				return nil, errors.New("the Twitter account is less than 3 years old." +
					" To whitelist your Twitter click here: https://forms.gle/fMaN1xtE322RBEYX8")
			}

			if userInfo.Followers < 200 {
				return nil, errors.New("the Twitter account has less tha 200 followers." +
					" To whitelist your Twitter click here: https://forms.gle/fMaN1xtE322RBEYX8")
			}
		}
	}

	tweetInfo, err := be.twitterClient.RetweetSearch(be.ctx, discordID, twitterName)
	if err != nil {
		return nil, err
	}

	discountCode, err := gonanoid.Generate("0123456789", 6)
	if err != nil {
		return nil, err
	}

	totalPrice := 50
	amountInPAC := 150
	if userInfo.Followers > 1000 {
		amountInPAC = 200
	}

	party := &store.TwitterParty{
		TwitterID:    userInfo.TwitterID,
		TwitterName:  userInfo.TwitterName,
		RetweetID:    tweetInfo.ID,
		ValAddr:      valAddr,
		ValPubKey:    pubKey,
		TotalPrice:   totalPrice,
		AmountInPAC:  amountInPAC,
		DiscountCode: discountCode,
		DiscordID:    discordID,
		CreatedAt:    time.Now().Unix(),
	}

	err = be.turboswap.SendDiscountCode(be.ctx, party)
	if err != nil {
		return nil, err
	}

	err = be.store.AddTwitterParty(party)
	if err != nil {
		return nil, err
	}

	return party, nil
}

func (be *BotEngine) TwitterCampaignStatus(twitterName string) (*store.TwitterParty, *turboswap.DiscountStatus, error) {
	party := be.store.FindTwitterParty(twitterName)
	if party == nil {
		return nil, nil, fmt.Errorf("no discount code generated for this Twitter account: `%v`", twitterName)
	}
	status, err := be.turboswap.GetStatus(be.ctx, party)
	if err != nil {
		return nil, nil, err
	}

	return party, status, nil
}

func (be *BotEngine) TwitterCampaignWhitelist(twitterName string, authorizedDiscordID string) error {
	authorizedIDs := []string{}

	if !slices.Contains(authorizedIDs, authorizedDiscordID) {
		return fmt.Errorf("unauthorize person")
	}

	foundParty := be.store.FindTwitterParty(twitterName)
	if foundParty != nil {
		return fmt.Errorf("the Twitter `%v` already registered for the campagna. Discount code is %v",
			foundParty.TwitterName, foundParty.DiscountCode)
	}

	userInfo, err := be.twitterClient.UserInfo(be.ctx, twitterName)
	if err != nil {
		return err
	}

	return be.store.WhitelistTwitterAccount(userInfo.TwitterID, userInfo.TwitterName, authorizedDiscordID)
}

func (be *BotEngine) Stop() {
	be.logger.Info("shutting bot engine down...")

	be.cancel()
	be.clientMgr.Stop()
}

func (be *BotEngine) Start() {
	be.logger.Info("starting the bot engine...")
}
