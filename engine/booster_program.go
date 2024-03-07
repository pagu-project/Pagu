package engine

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/utils"
	gonanoid "github.com/matoous/go-nanoid"
)

const (
	BoosterCommandName          = "booster"
	BoosterPaymentCommandName   = "payment"
	BoosterClaimCommandName     = "claim"
	BoosterWhitelistCommandName = "whitelist"
	BoosterStatusCommandName    = "status"
	BoosterHelpCommandName      = "help"
)

func (be *BotEngine) RegisterBoosterCommands() {
	subCmdBoosterPayment := Command{
		Name: BoosterPaymentCommandName,
		Desc: "Make a payment link for booster program",
		Help: "Provide your twitter username and mainnet validator address",
		Args: []Args{
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
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.boosterPaymentHandler,
	}

	subCmdBoosterClaim := Command{
		Name: BoosterClaimCommandName,
		Desc: "Claim your booster program stakes",
		Help: "You have to do the booster payment first, then try to claim it",
		Args: []Args{
			{
				Name:     "twitter-name",
				Desc:     "Your twitter user name without @",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.boosterClaimHandler,
	}

	subCmdBoosterWhitelist := Command{
		Name: BoosterWhitelistCommandName,
		Desc: "Whitelist a user for the booster program",
		Help: "This sub-command is **admin only**",
		Args: []Args{
			{
				Name:     "twitter-name",
				Desc:     "Your twitter user name without @",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		Handler:     be.boosterWhitelistHandler,
	}

	subCmdBoosterStatus := Command{
		Name:        BoosterStatusCommandName,
		Desc:        "Status of booster program claims and ...",
		Help:        "",
		Args:        []Args{},
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: nil,
		Handler:     be.boosterStatusHandler,
	}

	subCmdHelp := Command{
		Name: BoosterHelpCommandName,
		Desc: "Help for booster program commands",
		Help: "Provide the command name as parameter",
		Args: []Args{
			{
				Name:     "sub-command",
				Desc:     "The subcommand you want to see the related help of it",
				Optional: true,
			},
		},
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: nil,
		Handler:     be.boosterHelpHandler,
	}

	cmdBooster := Command{
		Name:        BoosterCommandName,
		Desc:        "Pactus validator booster program",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&subCmdBoosterClaim, &subCmdBoosterPayment, &subCmdBoosterStatus, &subCmdBoosterWhitelist, &subCmdHelp},
		Handler:     nil,
	}

	be.Cmds = append(be.Cmds, cmdBooster)
}

func (be *BotEngine) boosterPaymentHandler(_ AppID, callerID string, args ...string) (*CommandResult, error) {
	be.Lock()
	defer be.Unlock()

	boosterStatus := be.store.BoosterStatus()
	if boosterStatus.AllPkgs > 500 {
		return nil, errors.New("program is finished")
	}

	twitterName := args[0]
	valAddr := args[1]

	existingParty := be.store.FindTwitterParty(twitterName)
	if existingParty != nil {
		if existingParty.TransactionID != "" {
			return nil, fmt.Errorf("transaction is processed before: https://pacscan.org/transactions/%v", existingParty.TransactionID)
		} else {
			return nil, errors.New("")
		}
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
				return nil, errors.New("the Twitter account has less than 200 followers." +
					" To whitelist your Twitter click here: https://forms.gle/fMaN1xtE322RBEYX8")
			}
		}
	}

	tweetInfo, err := be.twitterClient.RetweetSearch(be.ctx, callerID, twitterName)
	if err != nil {
		return nil, err
	}

	discountCode, err := gonanoid.Generate("0123456789", 8)
	if err != nil {
		return nil, err
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

	err = be.nowpayments.CreatePayment(party)
	if err != nil {
		return nil, err
	}

	err = be.store.SaveTwitterParty(party)
	if err != nil {
		return nil, err
	}

	expiryDate := time.Unix(party.CreatedAt, 0).AddDate(0, 0, 7)

	result := fmt.Sprintf("Validator `%s` registered to receive %v stake-PAC coins in total price of $%v."+
		" Visit https://nowpayments.io/payment/?iid=%v to pay it."+
		" The Discount code will expire on %v",
		party.ValAddr, party.AmountInPAC, party.TotalPrice, party.NowPaymentsInvoiceID, expiryDate.Format("2006-01-02"))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) boosterClaimHandler(_ AppID, _ string, args ...string) (*CommandResult, error) {
	be.Lock()
	defer be.Unlock()

	twitterName := args[0]

	party := be.store.FindTwitterParty(twitterName)
	if party == nil {
		return nil, fmt.Errorf("no discount code generated for this Twitter account: `%v`", twitterName)
	}
	err := be.nowpayments.UpdatePayment(party)
	if err != nil {
		return nil, err
	}

	if party.NowPaymentsFinished {
		if party.TransactionID == "" {
			be.logger.Info("sending bond transaction", "receiver", party.ValAddr, "amount", party.AmountInPAC)
			memo := "Booster Program"
			txID, err := be.wallet.BondTransaction(party.ValPubKey, party.ValAddr, memo, utils.CoinToChange(float64(party.AmountInPAC)))
			if err != nil {
				return nil, err
			}

			if txID == "" {
				return nil, errors.New("can't send bond transaction")
			}

			party.TransactionID = txID

			err = be.store.SaveTwitterParty(party)
			if err != nil {
				return nil, err
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

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) boosterWhitelistHandler(_ AppID, callerID string, args ...string) (*CommandResult, error) {
	if !slices.Contains(be.AuthIDs, callerID) {
		return nil, fmt.Errorf("unauthorized person")
	}

	twitterName := args[0]

	foundParty := be.store.FindTwitterParty(twitterName)
	if foundParty != nil {
		return nil, fmt.Errorf("the Twitter `%v` already registered for the campaign. Discount code is %v",
			foundParty.TwitterName, foundParty.DiscountCode)
	}

	userInfo, err := be.twitterClient.UserInfo(be.ctx, twitterName)
	if err != nil {
		return nil, err
	}

	if err = be.store.WhitelistTwitterAccount(userInfo.TwitterID,
		userInfo.TwitterName, callerID); err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Twitter `%s` whitelisted", twitterName)

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) boosterStatusHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	bs := be.store.BoosterStatus()

	result := fmt.Sprintf("Total Coins: %v PAC\nTotal Packages: %v\nClaimed Packages: %v\nUnClaimed Packages: %v\nPayment Done: %v\nPayment Waiting: %v\nWhite Listed: %v\n",
		bs.Pac, bs.AllPkgs, bs.ClaimedPkgs, bs.UnClaimedPkgs, bs.PaymentDone, bs.PaymentWaiting, bs.Whitelists)

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) boosterHelpHandler(source AppID, callerID string, args ...string) (*CommandResult, error) {
	if len(args) == 0 {
		return be.help(source, callerID, BoosterCommandName)
	}
	return be.help(source, callerID, BoosterCommandName, args[0])
}
