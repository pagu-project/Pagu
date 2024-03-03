package engine

import (
	"context"
	"errors"
	"sync"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/database"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/nowpayments"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/twitter_api"
	"github.com/kehiy/RoboPac/wallet"
)

type BotEngine struct {
	ctx    context.Context //nolint
	cancel func()

	wallet        wallet.IWallet
	db            *database.DB
	nowpayments   nowpayments.INowpayment
	clientMgr     *client.Mgr
	logger        *log.SubLogger
	twitterClient twitter_api.IClient

	AuthIDs []string
	Cmds    []Command

	store        store.IStore //!
	sync.RWMutex              //! remove this.
}

func NewBotEngine(cfg *config.Config) (IEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cm := client.NewClientMgr(ctx)

	localClient, err := client.NewClient(cfg.LocalNode)
	if err != nil {
		cancel()
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
	cm.Start()

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
		cancel()
		return nil, errors.New("can't open the wallet")
	}

	log.Info("wallet opened successfully", "address", wallet.Address())

	// load store.
	store, err := store.NewStore(cfg.StorePath, sSl)
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("store loaded successfully", "path", cfg.StorePath)

	// twitter
	twitterClient, err := twitter_api.NewClient(cfg.TwitterAPICfg.BearerToken, cfg.TwitterAPICfg.TwitterID)
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("twitterClient loaded successfully")

	// load database
	db, err := database.NewDB(cfg.DataBasePath)
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("database loaded successfully")

	nowpayments, err := nowpayments.NewNowPayments(&cfg.NowPaymentsConfig)
	if err != nil {
		log.Error("could not start twitter client", "err", err)
	}
	log.Info("nowpayments loaded successfully")

	return newBotEngine(eSl, cm, wallet, store, db, twitterClient, nowpayments, cfg.AuthIDs, ctx, cancel), nil
}

func newBotEngine(logger *log.SubLogger, cm *client.Mgr, w wallet.IWallet, s store.IStore, db *database.DB,
	twitterClient twitter_api.IClient, nowpayments nowpayments.INowpayment, authIDs []string,
	ctx context.Context, cnl context.CancelFunc,
) *BotEngine {
	return &BotEngine{
		ctx:           ctx,
		cancel:        cnl,
		logger:        logger,
		wallet:        w,
		clientMgr:     cm,
		store:         s,
		db:            db,
		twitterClient: twitterClient,
		nowpayments:   nowpayments,
		AuthIDs:       authIDs,
	}
}

func (be *BotEngine) Stop() {
	be.logger.Info("shutting bot engine down...")

	be.cancel()
	be.clientMgr.Stop()
}

func (be *BotEngine) Start() {
	be.logger.Info("starting the bot engine...")
}
