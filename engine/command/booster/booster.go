package booster

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/engine/command"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/nowpayments"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/twitter_api"
	"github.com/kehiy/RoboPac/utils"
	"github.com/kehiy/RoboPac/wallet"
	gonanoid "github.com/matoous/go-nanoid"
)

const (
	boosterCommandName          = "booster"
	boosterPaymentCommandName   = "payment"
	boosterClaimCommandName     = "claim"
	boosterWhitelistCommandName = "whitelist"
	boosterStatusCommandName    = "status"
)

func boosterPrice(allPackages int) int {
	switch {
	case allPackages < 100:
		return 30
	case allPackages < 200:
		return 40
	case allPackages < 300:
		return 50
	default:
		return 100
	}
}

type Booster struct {
	sync.RWMutex //! remove this.

	ctx           context.Context
	AdminIDs      []string
	store         store.IStore
	wallet        wallet.IWallet
	nowpayments   nowpayments.INowpayment
	clientMgr     *client.Mgr
	twitterClient twitter_api.IClient
}

func NewBooster(ctx context.Context,
	adminIDs []string,
	store store.IStore,
	wallet wallet.IWallet,
	nowpayments nowpayments.INowpayment,
	clientMgr *client.Mgr,
	twitterClient twitter_api.IClient,
) *Booster {
	return &Booster{
		ctx:           ctx,
		AdminIDs:      adminIDs,
		store:         store,
		wallet:        wallet,
		nowpayments:   nowpayments,
		clientMgr:     clientMgr,
		twitterClient: twitterClient,
	}
}

func (booster *Booster) GetCommand() *command.Command {
	subCmdBoosterPayment := &command.Command{
		Name: boosterPaymentCommandName,
		Desc: "Make a payment link for booster program",
		Help: "Provide your twitter username and mainnet validator address",
		Args: []command.Args{
			{
				Name:     "twitter-name",
				Desc:     "Your twitter user name without @",
				Optional: false,
			},
			{
				Name:     "validator-address",
				Desc:     "Your validator address to be registered",
				Optional: false,
			},
		},
		AppIDs:  []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler: booster.boosterPaymentHandler,
	}

	subCmdBoosterClaim := &command.Command{
		Name: boosterClaimCommandName,
		Desc: "Claim your booster program stake",
		Help: "You have to do the booster payment first, then try to claim it",
		Args: []command.Args{
			{
				Name:     "twitter-name",
				Desc:     "Your twitter user name without @",
				Optional: false,
			},
		},
		AppIDs:  []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler: booster.boosterClaimHandler,
	}

	subCmdBoosterWhitelist := &command.Command{
		Name: boosterWhitelistCommandName,
		Desc: "Whitelist a user for the booster program",
		Help: "This sub-command is **admin only**",
		Args: []command.Args{
			{
				Name:     "twitter-name",
				Desc:     "Your twitter user name without @",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler:     booster.boosterWhitelistHandler,
	}

	subCmdBoosterStatus := &command.Command{
		Name:        boosterStatusCommandName,
		Desc:        "Status of booster program claims and ...",
		Help:        "",
		Args:        []command.Args{},
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		SubCommands: nil,
		Handler:     booster.boosterStatusHandler,
	}

	cmdBooster := &command.Command{
		Name:    boosterCommandName,
		Emoji:   "✨",
		Color:   "#50C878",
		Desc:    "Pactus Validator Booster Program",
		Help:    "",
		Args:    nil,
		AppIDs:  []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler: nil,
	}

	cmdBooster.AddSubCommand(subCmdBoosterClaim)
	cmdBooster.AddSubCommand(subCmdBoosterPayment)
	cmdBooster.AddSubCommand(subCmdBoosterStatus)
	cmdBooster.AddSubCommand(subCmdBoosterWhitelist)

	return cmdBooster
}

func (booster *Booster) boosterPaymentHandler(cmd *command.Command, _ command.AppID, callerID string, args ...string) *command.CommandResult {
	booster.Lock()
	defer booster.Unlock()

	boosterStatus := booster.store.BoosterStatus()
	if boosterStatus.AllPkgs > 500 {
		return cmd.FailedResult("program is finished")
	}

	twitterName := args[0]
	valAddr := args[1]

	existingParty := booster.store.FindTwitterParty(twitterName)
	if existingParty != nil {
		if existingParty.TransactionID != "" {
			return cmd.FailedResult("transaction is processed before: https://pacscan.org/transactions/%v", existingParty.TransactionID)
		} else {
			return cmd.FailedResult("Payment is not done")
		}
	}

	valInfo, _ := booster.clientMgr.GetValidatorInfo(valAddr)
	if valInfo != nil {
		return cmd.FailedResult("this address is already a staked validator")
	}

	pubKey, err := booster.clientMgr.FindPublicKey(valAddr, false)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	userInfo, err := booster.twitterClient.UserInfo(booster.ctx, twitterName)
	if err != nil {
		return cmd.ErrorResult(err)
	}
	if !userInfo.IsVerified {
		if !booster.store.IsWhitelisted(userInfo.TwitterID) {
			threeYearsAgo := time.Now().AddDate(-3, 0, 0)
			if userInfo.CreatedAt.After(threeYearsAgo) {
				return cmd.FailedResult("the Twitter account is less than 3 years old." +
					" To whitelist your Twitter click here: https://forms.gle/fMaN1xtE322RBEYX8")
			}

			if userInfo.Followers < 200 {
				return cmd.FailedResult("the Twitter account has less than 200 followers." +
					" To whitelist your Twitter click here: https://forms.gle/fMaN1xtE322RBEYX8")
			}
		}
	}

	tweetInfo, err := booster.twitterClient.RetweetSearch(booster.ctx, callerID, twitterName)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	discountCode, err := gonanoid.Generate("0123456789", 8)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	totalPrice := boosterPrice(boosterStatus.AllPkgs)
	amountInPAC := int64(150)
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
		DiscordID:    callerID,
		CreatedAt:    time.Now().Unix(),
	}

	err = booster.nowpayments.CreatePayment(party)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	err = booster.store.SaveTwitterParty(party)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)

	result := fmt.Sprintf("Validator `%s` registered to receive %v stake-PAC coins in total price of $%v."+
		" Visit https://nowpayments.io/payment/?iid=%v to pay it."+
		" The Discount code will expire on %v",
		party.ValAddr, party.AmountInPAC, party.TotalPrice, party.NowPaymentsInvoiceID, expiryDate.Format("2006-01-02"))

	return cmd.SuccessfulResult(result)
}

func (booster *Booster) boosterClaimHandler(cmd *command.Command, _ command.AppID, _ string, args ...string) *command.CommandResult {
	booster.Lock()
	defer booster.Unlock()

	twitterName := args[0]

	party := booster.store.FindTwitterParty(twitterName)
	if party == nil {
		return cmd.FailedResult("no discount code generated for this Twitter account: `%v`", twitterName)
	}
	err := booster.nowpayments.UpdatePayment(party)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if party.NowPaymentsFinished {
		if party.TransactionID == "" {
			log.Info("sending bond transaction", "receiver", party.ValAddr, "amount", party.AmountInPAC)
			memo := "Booster Program"
			txID, err := booster.wallet.BondTransaction(party.ValPubKey, party.ValAddr, memo, utils.CoinToChange(float64(party.AmountInPAC)))
			if err != nil {
				return cmd.ErrorResult(err)
			}

			if txID == "" {
				return cmd.FailedResult("can't send bond transaction")
			}

			party.TransactionID = txID

			err = booster.store.SaveTwitterParty(party)
			if err != nil {
				return cmd.ErrorResult(err)
			}
		}
	}

	var result string
	if party.NowPaymentsFinished {
		result = fmt.Sprintf("Validator `%s` received %v stake-PAC coins."+
			" Transaction: https://pacscan.org/transactions/%v.",
			party.ValAddr, party.AmountInPAC, party.TransactionID)
	} else {
		expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)
		result = fmt.Sprintf("Validator `%s` registered to receive %v stake-PAC coins in total price of $%v."+
			" Visit https://nowpayments.io/payment/?iid=%v and pay the total amount."+
			" The Discount code will expire on %v",
			party.ValAddr, party.AmountInPAC, party.TotalPrice, party.NowPaymentsInvoiceID, expiryDate.Format("2006-01-02"))
	}

	return cmd.SuccessfulResult(result)
}

func (booster *Booster) boosterWhitelistHandler(cmd *command.Command, _ command.AppID, callerID string, args ...string) *command.CommandResult {
	if !slices.Contains(booster.AdminIDs, callerID) {
		return cmd.FailedResult("unauthorized person")
	}

	twitterName := args[0]

	foundParty := booster.store.FindTwitterParty(twitterName)
	if foundParty != nil {
		return cmd.FailedResult("the Twitter `%v` already registered for the campaign. Discount code is %v",
			foundParty.TwitterName, foundParty.DiscountCode)
	}

	userInfo, err := booster.twitterClient.UserInfo(booster.ctx, twitterName)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = booster.store.WhitelistTwitterAccount(userInfo.TwitterID,
		userInfo.TwitterName, callerID); err != nil {
		return cmd.ErrorResult(err)
	}

	result := fmt.Sprintf("Twitter `%s` whitelisted", twitterName)

	return cmd.SuccessfulResult(result)
}

func (booster *Booster) boosterStatusHandler(cmd *command.Command, _ command.AppID, _ string, _ ...string) *command.CommandResult {
	bs := booster.store.BoosterStatus()

	result := fmt.Sprintf("Total Coins: %v PAC\nTotal Packages: %v\nClaimed Packages: %v\nUnClaimed Packages: %v\nPayment Done: %v\nPayment Waiting: %v\nWhite Listed: %v\n",
		bs.Pac, bs.AllPkgs, bs.ClaimedPkgs, bs.UnClaimedPkgs, bs.PaymentDone, bs.PaymentWaiting, bs.Whitelists)

	return cmd.SuccessfulResult(result)
}
