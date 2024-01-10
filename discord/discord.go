package discord

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/libp2p/go-libp2p/core/peer"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Bot struct {
	discordSession *discordgo.Session
	faucetWallet   *wallet.Wallet
	cfg            *config.Config
	store          *SafeStore
	referralStore  *ReferralStore

	cm *client.Mgr
}

// guildID: "795592769300987944"

func Start(cfg *config.Config, w *wallet.Wallet, ss *SafeStore, rs *ReferralStore) (*Bot, error) {
	cm := client.NewClientMgr()

	for _, s := range cfg.Servers {
		c, err := client.NewClient(s)
		if err != nil {
			log.Printf("unable to create client at: %s. err: %s", s, err)
		} else {
			log.Printf("adding client at: %s", s)
			cm.AddClient(s, c)
		}
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Printf("error creating Discord session: %v", err)
		return nil, err
	}
	bot := &Bot{cfg: cfg, discordSession: dg, faucetWallet: w, store: ss, cm: cm, referralStore: rs}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(bot.messageHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection: %v", err)
		return nil, err
	}

	return bot, nil
}

func (b *Bot) Stop() error {
	return b.discordSession.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
// nolint.
func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	p := message.NewPrinter(language.English)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore messages which is not from accepted channels
	if m.GuildID != "795592769300987944" {
		return
	}

	if strings.ToLower(m.Content) == "help" {
		help(s, m)
		return
	}

	if strings.ToLower(m.Content) == "network" {
		msg := b.networkInfo()
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	if strings.ToLower(m.Content) == "address" {
		msg := fmt.Sprintf("Faucet address is %v", b.cfg.FaucetAddress)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	// If the message is "balance" reply with "available faucet balance"
	if strings.ToLower(m.Content) == "balance" {
		balance := b.faucetWallet.GetBalance()
		v, d := b.store.GetDistribution()
		msg := p.Sprintf("Available faucet balance is %.4f tPAC's🪙\n", balance.Available)
		msg += p.Sprintf("A total of %.4f tPAC's has been distributed to %d validators.\n", d, v)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	if strings.ToLower(m.Content) == "health" {
		lastBlockTime, LastBlockHeight := b.cm.GetLastBlockTime()
		lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0)
		currentTime := time.Now()

		timeDiff := (currentTime.Unix() - int64(lastBlockTime))
		if timeDiff > 15 {
			msg := p.Sprintf("Network is **unhealthy❌**\nLast block time⛓️: %v\nCurrent time🕧: %v\nTime Difference: %v seconds\nLast block height⛓️: %v\nDifference is more than 15 seconds.",
				lastBlockTimeFormatted.Format("02/01/2006, 15:04:05"), currentTime.Format("02/01/2006, 15:04:05"), timeDiff, LastBlockHeight)
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		msg := p.Sprintf("Network is **healthy✅**\nLast block time⛓️: %v\nCurrent time🕧: %v\nTime Difference: %v seconds\nLast block height⛓️: %v\nDifference is less than 15 seconds.",
			lastBlockTimeFormatted.Format("02/01/2006, 15:04:05"), currentTime.Format("02/01/2006, 15:04:05"), timeDiff, LastBlockHeight)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	if strings.Contains(strings.ToLower(m.Content), "peer-info") {
		trimmedPrefix := strings.TrimPrefix(strings.ToLower(m.Content), "peer-info")
		trimmedAddress := strings.Trim(trimmedPrefix, " ")

		peerInfo, _, err := b.cm.GetPeerInfo(trimmedAddress)
		if err != nil {
			msg := p.Sprintf("An error occurred %v\n", err)
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		peerID, err := peer.IDFromBytes(peerInfo.PeerId)
		if err != nil {
			msg := p.Sprintf("An error occurred %v\n", err)
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		parts := strings.Split(strings.Split(peerInfo.Address, "/")[2], "/")
		ip := parts[0]
		fmt.Println(ip)
		geoData := getGeoIP(ip)

		val, err := b.cm.GetValidatorInfo(trimmedAddress)
		if err != nil {
			msg := p.Sprintf("An error occurred %v\n", err)
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		score := ""
		if val.Validator.AvailabilityScore < 0.8 {
			score = fmt.Sprintf("%v⚠️⚠️", val.Validator.AvailabilityScore)
		}

		if val.Validator.AvailabilityScore >= 0.8 {
			score = fmt.Sprintf("%v🟢✅", val.Validator.AvailabilityScore)
		}

		msg := p.Sprintf("--------------------Peer Info--------------------\n")
		msg += p.Sprintf("Peer ID: %v\n", peerID)
		msg += p.Sprintf("IP address: %v\n", peerInfo.Address)
		msg += p.Sprintf("Agent: %v\n", peerInfo.Agent)
		msg += p.Sprintf("Moniker: %v\n", peerInfo.Moniker)
		msg += p.Sprintf("Country: %v\n", geoData.CountryName)
		msg += p.Sprintf("City: %v\n", geoData.City)
		msg += p.Sprintf("Region Name: %v\n", geoData.RegionName)
		msg += p.Sprintf("TimeZone: %v\n", geoData.TimeZone)
		msg += p.Sprintf("ISP: %v\n", geoData.ISP)
		msg += p.Sprintf("--------------------Validator Info----------------\n")
		msg += p.Sprintf("Number: %v\n", val.Validator.Number)
		msg += p.Sprintf("**Availability score: %v\n**", score)
		msg += p.Sprintf("Stake amount: %v tPAC's\n", util.ChangeToCoin(val.Validator.Stake))
		msg += p.Sprintf("Last bonding height: %v\n", val.Validator.LastBondingHeight)
		msg += p.Sprintf("Last sortition height: %v\n", val.Validator.LastSortitionHeight)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}
	if strings.Contains(strings.ToLower(m.Content), "synced") {
		trimmedPrefix := strings.TrimPrefix(strings.ToLower(m.Content), "peer-info")
		trimmedAddress := strings.Trim(trimmedPrefix, " ")

		peerInfo, _, err := b.cm.GetPeerInfo(trimmedAddress)
		if err != nil {
			msg := p.Sprintf("An error occurred %v\n", err)
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		notSyncedMsg := "this peer is not synced with network, gRPC is disabled or doesn't have public IP address."
		syncedMsg := "**this peer is synced with network**"

		msg := notSyncedMsg
		c, err := client.NewClient(strings.Split(peerInfo.Address, "/")[2] + ":50052")
		if err != nil {
			msg = notSyncedMsg
		}
		lastBlockTime, _, err := c.LastBlockTime()
		if err != nil {
			msg = notSyncedMsg
		}
		currentTime := time.Now().Unix()

		if (uint32(currentTime) - lastBlockTime) < 15 {
			msg = syncedMsg
		}

		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	if strings.ToLower(m.Content) == "my-referral" {
		referrals := b.referralStore.GetAllReferrals()
		for _, r := range referrals {
			if r.DiscordID == m.Author.ID {
				msg := fmt.Sprintf("Your referral information👥:\nPoints: %v (%v tPACs)\nCode: ```%v```\n", r.Points, (r.Points * 10), r.ReferralCode)
				_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
				return
			}
		}

		referralCode, err := gonanoid.Generate("0123456789", 6)
		if err != nil {
			msg := "can't generate referral code, please try again later."
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		err = b.referralStore.NewReferral(m.Author.ID, m.Author.Username, referralCode)
		if err != nil {
			msg := "can't generate referral code, please try again later."
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		msg := fmt.Sprintf("Your referral information👥:\nPoints: %v (%v tPAC's)\nCode: ```%v```\n", 0, 0, referralCode)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	if strings.Contains(strings.ToLower(m.Content), "faucet-referral") || strings.Contains(strings.ToLower(m.Content), "faucet") {
		msg := "Hi, faucet and referral campaign for testnet-2 is closed now!\nCheck this post:\nhttps://ptb.discord.com/channels/795592769300987944/811878389304655892/1192018704855220234"
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	// if m.Content == "pip19-score-average" {
	// 	totalValidators := float64(0)
	// 	sumScores := float64(0)

	// 	for i := 0; i < 2065; i++ {
	// 		val, err := b.cm.GetValidatorInfoByNumber(int32(i))
	// 		if err != nil {
	// 			continue
	// 		}
	// 		totalValidators += 1
	// 		sumScores += val.Validator.AvailabilityScore
	// 		fmt.Println(i)
	// 	}

	// 	msg := fmt.Sprintf("average of %v validators pip19 score is %f", totalValidators, sumScores/totalValidators)
	// 	_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
	// 	return
	// }

	if m.Content == "pip19-report" {
		t := time.Now()

		totalActiveValidators := 0
		scoresSum := float64(0)
		notActiveNodes := 0

		results := []Result{}

		info, err := b.cm.GetNetworkInfo()
		if err != nil {
			msg := "error getting network info"
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		fmt.Printf("total peers: %v\n", info.ConnectedPeersCount)

		for i, p := range info.ConnectedPeers {
			fmt.Printf("new peer %v\n", i)
			r := Result{}
			r.Agent = p.Agent
			r.RemoteAddress = p.Address
			if p.Height < 682_000 {
				fmt.Printf("new peer %v is not active\n", i)
				notActiveNodes += 1
				continue
			}
			for iv, v := range p.ConsensusKeys {
				fmt.Printf("new validator %v\n", iv)
				val, err := b.cm.GetValidatorInfo(v)
				if err != nil {
					continue
				}
				r.PIP19Score = val.Validator.AvailabilityScore
				r.ValidatorAddress = v

				results = append(results, r)
				totalActiveValidators += 1
				scoresSum += val.Validator.AvailabilityScore
			}
		}

		data, err := json.Marshal(results)
		if err != nil {
			msg := "error saving report"
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		if err = os.WriteFile("pip19Report.json", data, 0o600); err != nil {
			msg := "error saving report"
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		msg := fmt.Sprintf("Time: %v\nSum Scores: %v\nNot Active Nodes: %v\nActive Validators: %v\nTotal Nodes: %v\n",
			t.Format("04:05"), scoresSum, notActiveNodes, totalActiveValidators, info.ConnectedPeersCount)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}
}

// help sends a message detailing how to use the bot discord-client side
// nolint.
func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Pactus Universal Robot",
		URL:   "https://pactus.org",
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://pactus.org",
			IconURL: s.State.User.AvatarURL(""),
			Name:    s.State.User.Username,
		},
		Description: "RoboPac is a robot that provides support and information about the Pactus Blockchain.\n" +
			"To see the faucet account balance, simply type: `balance`\n" +
			"To see the faucet address, simply type: `address`\n" +
			"To get network information, simply type: `network`\n" +
			"To get network health status, simply type: `health`\n" +
			"To get your node syncing status, simply type: `synced [validator address]`\n" +
			"To get peer information, simply type: `peer-info [validator address]`\n" +
			"To get your referral information, simply type: `my-referral`\n" +
			"To request faucet for test network *with referral code*: simply type `faucet-referral [validator address] [referral code]`\nreferral faucet will get 100 tPAC's\n" +
			"To request faucet for test network: simply type `faucet [validator address]`. normal faucet will get 60 tPAC's",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Example of requesting `faucet` ",
				Value: "faucet tpc1pxl333elgnrdtk0kjpjdvky44yu62x0cwupnpjl",
			},
		},
	})
}

func (b *Bot) networkInfo() string {
	msg := "Pactus is truly decentralized Proof of Stake Blockchain.⛓️"
	nodes, err := b.cm.GetNetworkInfo()
	if err != nil {
		log.Printf("error establishing connection")
		return msg
	}
	msg += "\n📊 The current statistics are:\n"
	// msg += fmt.Sprintf("🕧Node started at: %v\n", time.UnixMilli(nodes.sta*1000).Format("02/01/2006, 15:04:05"))
	msg += fmt.Sprintf("⬆️Total bytes sent: %v\n", uint32(nodes.TotalSentBytes))
	msg += fmt.Sprintf("⬇️Total bytes received: %v\n", uint32(nodes.TotalReceivedBytes))
	msg += fmt.Sprintf("👾Number of connected peers: %v\n", nodes.ConnectedPeersCount)
	// check block height
	blockchainInfo, err := b.cm.GetBlockchainInfo()
	if err != nil {
		log.Printf("error current block height")
		return msg
	}
	msg += fmt.Sprintf("⛓️Block height: %v\n", blockchainInfo.LastBlockHeight)
	msg += fmt.Sprintf("🦾Total power: %.4f PACs\n", util.ChangeToCoin(blockchainInfo.TotalPower))
	msg += fmt.Sprintf("🦾Total committee power: %.4f PACs\n", util.ChangeToCoin(blockchainInfo.CommitteePower))
	msg += fmt.Sprintf("✔️Total validators: %v\n", blockchainInfo.TotalValidators)
	return msg
}

func (b *Bot) GetPeerInfo(address string) (*pactus.PeerInfo, error) {
	_, err := crypto.AddressFromString(address)
	if err != nil {
		log.Printf("invalid address")

		return nil, err
	}

	_, err = b.cm.IsValidator(address)
	if err != nil {
		return nil, err
	}

	peerInfo, _, err := b.cm.GetPeerInfoFirstVal(address)
	if err != nil {
		return nil, err
	}
	return peerInfo, nil
}
