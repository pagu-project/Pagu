package engine

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/kehiy/RoboPac/database"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/utils"
	"github.com/libp2p/go-libp2p/core/peer"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
)

func (be *BotEngine) networkHealthHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	lastBlockTime, lastBlockHeight := be.clientMgr.GetLastBlockTime()
	lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0).Format("02/01/2006, 15:04:05")
	currentTime := time.Now()

	timeDiff := (currentTime.Unix() - int64(lastBlockTime))

	healthStatus := true
	if timeDiff > 15 {
		healthStatus = false
	}

	var status string
	if healthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return &CommandResult{
		Successful: true,
		Message: fmt.Sprintf("Network is %s\nCurrentTime: %v\nLastBlockTime: %v\nTime Diff: %v\nLast Block Height: %v",
			status, currentTime.Format("02/01/2006, 15:04:05"), lastBlockTimeFormatted, timeDiff, utils.FormatNumber(int64(lastBlockHeight))),
	}, nil
}

func (be *BotEngine) networkStatusHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
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

	net := NetStatus{
		ValidatorsCount:     chainInfo.TotalValidators,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   chainInfo.TotalPower,
		TotalCommitteePower: chainInfo.CommitteePower,
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   cs,
	}

	result := fmt.Sprintf("Network Name: %s\nConnected Peers: %v\n"+
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v PAC\nTotal Committee Power: %v PAC\nCirculating Supply: %v PAC\n"+
		"\n> Note📝: This info is from one random network node. Non-blockchain data may not be consistent.",
		net.NetworkName,
		utils.FormatNumber(int64(net.ConnectedPeersCount)),
		utils.FormatNumber(int64(net.ValidatorsCount)),
		utils.FormatNumber(int64(net.TotalAccounts)),
		utils.FormatNumber(int64(net.CurrentBlockHeight)),
		utils.FormatNumber(int64(util.ChangeToCoin(net.TotalNetworkPower))),
		utils.FormatNumber(int64(util.ChangeToCoin(net.TotalCommitteePower))),
		utils.FormatNumber(int64(util.ChangeToCoin(net.CirculatingSupply))))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) nodeInfoHandler(_ AppID, _ string, args ...string) (*CommandResult, error) {
	valAddress := args[0]

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

	nodeInfo := &NodeInfo{
		PeerID:     peerID.String(),
		IPAddress:  peerInfo.Address,
		Agent:      peerInfo.Agent,
		Moniker:    peerInfo.Moniker,
		Country:    geoData.CountryName,
		City:       geoData.City,
		RegionName: geoData.RegionName,
		TimeZone:   geoData.TimeZone,
		ISP:        geoData.ISP,
	}

	// here we check if the node is also a validator.
	// if its a validator , then we populate the validator data.
	// if not validator then we set everything to 0/empty .
	val, err := be.clientMgr.GetValidatorInfo(valAddress)
	if err == nil && val != nil {
		nodeInfo.ValidatorNum = val.Validator.Number
		nodeInfo.AvailabilityScore = val.Validator.AvailabilityScore
		nodeInfo.StakeAmount = val.Validator.Stake
		nodeInfo.LastBondingHeight = val.Validator.LastBondingHeight
		nodeInfo.LastSortitionHeight = val.Validator.LastSortitionHeight
	} else {
		nodeInfo.ValidatorNum = 0
		nodeInfo.AvailabilityScore = 0
		nodeInfo.StakeAmount = 0
		nodeInfo.LastBondingHeight = 0
		nodeInfo.LastSortitionHeight = 0
	}

	var pip19Score string
	if nodeInfo.AvailabilityScore >= 0.9 {
		pip19Score = fmt.Sprintf("%v✅", nodeInfo.AvailabilityScore)
	} else {
		pip19Score = fmt.Sprintf("%v⚠️", nodeInfo.AvailabilityScore)
	}

	result := fmt.Sprintf("PeerID: %s\nIP Address: %s\nAgent: %s\n"+
		"Moniker: %s\nCountry: %s\nCity: %s\nRegion Name: %s\nTimeZone: %s\n"+
		"ISP: %s\n\nValidator Info🔍\nNumber: %v\nPIP-19 Score: %s\nStake: %v PAC's\n",
		nodeInfo.PeerID, nodeInfo.IPAddress, nodeInfo.Agent, nodeInfo.Moniker, nodeInfo.Country,
		nodeInfo.City, nodeInfo.RegionName, nodeInfo.TimeZone, nodeInfo.ISP, utils.FormatNumber(int64(nodeInfo.ValidatorNum)),
		pip19Score, utils.FormatNumber(int64(util.ChangeToCoin(nodeInfo.StakeAmount))))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) claimerInfoHandler(_ AppID, _ string, args ...string) (*CommandResult, error) {
	be.RLock()
	defer be.RUnlock()

	testNetValAddr := args[0]

	claimer := be.store.ClaimerInfo(testNetValAddr)
	if claimer == nil {
		return nil, errors.New("not found")
	}

	return &CommandResult{
		Successful: true,
		Message: fmt.Sprintf("TestNet Address: %s\namount: %v PACs\nIsClaimed: %v\n txHash: %s",
			args[0], util.ChangeToString(claimer.TotalReward), claimer.IsClaimed(), claimer.ClaimedTxID),
	}, nil
}

func (be *BotEngine) claimHandler(_ AppID, callerID string, args ...string) (*CommandResult, error) {
	be.Lock()
	defer be.Unlock()

	mainnetAddr := args[0]
	testnetAddr := args[1]

	be.logger.Info("new claim request", "mainnetAddr", mainnetAddr, "testnetAddr", testnetAddr, "discordID", callerID)

	valInfo, _ := be.clientMgr.GetValidatorInfo(mainnetAddr)
	if valInfo != nil {
		return nil, errors.New("this address is already a staked validator")
	}

	if utils.ChangeToCoin(be.wallet.Balance()) <= 500 {
		be.logger.Warn("bot wallet hasn't enough balance")
		return nil, errors.New("insufficient wallet balance")
	}

	claimer := be.store.ClaimerInfo(testnetAddr)
	if claimer == nil {
		return nil, errors.New("claimer not found")
	}

	if claimer.DiscordID != callerID {
		be.logger.Warn("try to claim other's reward", "claimer", claimer.DiscordID, "discordID", callerID)
		return nil, errors.New("invalid claimer")
	}

	if claimer.IsClaimed() {
		return nil, errors.New("this claimer have already claimed rewards")
	}

	pubKey, err := be.clientMgr.FindPublicKey(mainnetAddr, true)
	if err != nil {
		return nil, err
	}

	memo := "TestNet reward claim from RoboPac"
	txID, err := be.wallet.BondTransaction(pubKey, mainnetAddr, memo, claimer.TotalReward)
	if err != nil {
		return nil, err
	}

	if txID == "" {
		return nil, errors.New("can't send bond transaction")
	}

	be.logger.Info("new bond transaction sent", "txID", txID)

	err = be.store.AddClaimTransaction(testnetAddr, txID)
	if err != nil {
		be.logger.Panic("unable to add the claim transaction",
			"error", err,
			"discordID", callerID,
			"testnetAddr", testnetAddr,
			"txID", txID,
		)

		return nil, err
	}

	return &CommandResult{
		Successful: true,
		Message:    fmt.Sprintf("Reward claimed successfully✅\nYour claim transaction: https://pacscan.org/transactions/%s", txID),
	}, nil
}

func (be *BotEngine) walletHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	addr, blnc := be.wallet.Address(), be.wallet.Balance()

	result := fmt.Sprintf("Address: https://pacscan.org/address/%s\nBalance: %v PAC\n", addr, utils.FormatNumber(int64(util.ChangeToCoin(blnc))))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) claimStatusHandler(_ AppID, _ string, _ ...string) (*CommandResult, error) {
	cs := be.store.ClaimStatus()

	result := fmt.Sprintf("Claimed rewards count: %v\nClaimed coins: %v PAC's\nNot-claimed rewards count: %v\nNot-claim coins: %v PAC's\n",
		cs.Claimed, util.ChangeToString(cs.ClaimedAmount), cs.NotClaimed, util.ChangeToString(cs.NotClaimedAmount))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}

func (be *BotEngine) calcRewardHandler(_ AppID, _ string, args ...string) (*CommandResult, error) {
	stake, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	time := args[1]

	if stake < 1 || stake > 1_000 {
		return nil, errors.New("minimum of stake is 1 PAC and maximum is 1,000 PAC")
	}

	var blocks int
	switch time {
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
		return nil, err
	}

	reward := int64(stake*blocks) / int64(util.ChangeToCoin(bi.TotalPower))

	result := fmt.Sprintf("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator in one %s ⏰ with %v PAC total power ⚡ of committee."+
		"\n\n> Note📝: This is an estimation and the number can get changed by changes of your stake amount, total power and ...",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), time, utils.FormatNumber(bi.TotalPower))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
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
			logger.Info("sending bond transaction", "receiver", party.ValAddr, "amount", party.AmountInPAC)
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

func (be *BotEngine) depositAddressHandler(_ AppID, callerID string, _ ...string) (*CommandResult, error) {
	u, err := be.db.GetUser(callerID)
	if err == nil {
		return MakeSuccessfulResult(
			"You already have a deposit address: %s", u.DepositAddress,
		), nil
	}

	addr, err := be.wallet.NewAddress(fmt.Sprintf("deposit address for %s", callerID))
	if err != nil {
		return MakeFailedResult(
			"can't make a new address: %v", err,
		), nil
	}

	err = be.db.AddUser(
		&database.DiscordUser{
			DiscordID:      callerID,
			DepositAddress: addr,
		},
	)
	if err != nil {
		return MakeFailedResult(
			"can't add discord user to database: %v", err,
		), nil
	}

	return MakeSuccessfulResult(
		"Deposit address crated for you successfully: %s", addr,
	), nil
}

func (be *BotEngine) help(source AppID, _ string, args ...string) (*CommandResult, error) {
	helpStr := ""
	if len(args) > 0 {
		cmdName := args[0]
		cmd := be.commandByName(cmdName)
		if cmd == nil {
			return nil, fmt.Errorf("unknown command: %s", cmdName)
		}

		argsStr := ""
		for _, arg := range cmd.Args {
			argsStr += fmt.Sprintf("<%v> ", arg.Name)
		}
		argsStr = argsStr[:len(argsStr)-1]

		helpStr += cmd.Desc
		helpStr += fmt.Sprintf("%v\nUsage: `%v %v`", cmd.Help, cmd.Name, argsStr)
	} else {
		helpStr += "List of available commands:\n"
		for _, cmd := range be.Cmds {
			if !slices.Contains(cmd.AppIDs, source) {
				continue
			}

			padding := 12 - len(cmd.Name)
			helpStr += fmt.Sprintf("`%s`:%s%v\n", cmd.Name, strings.Repeat(" ", padding), cmd.Desc)
		}
	}
	return MakeSuccessfulResult(helpStr), nil
}
