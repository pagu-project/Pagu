package discord

import (
	"fmt"
	"log"
	"strings"
	"time"

	"pactus-bot/client"
	"pactus-bot/config"
	"pactus-bot/wallet"

	"github.com/bwmarrin/discordgo"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Bot struct {
	discordSession *discordgo.Session
	faucetWallet   *wallet.Wallet
	cfg            *config.Config
	store          *SafeStore
}

func Start(cfg *config.Config, w *wallet.Wallet, ss *SafeStore) (*Bot, error) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Printf("error creating Discord session: %v", err)
		return nil, err
	}
	bot := &Bot{cfg: cfg, discordSession: dg, faucetWallet: w, store: ss} // TODO: remove this hard coded id

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
func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	p := message.NewPrinter(language.English)
	// log.Printf("received message: %v\n", m.Content)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "help" {
		help(s, m)
		return
	}

	if m.Content == "network" {
		msg := b.networkInfo()
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}
	if m.Content == "address" {
		msg := fmt.Sprintf("Faucet address is: %v", b.cfg.FaucetAddress)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}
	// If the message is "balance" reply with "available faucet balance"
	if m.Content == "balance" {
		balance := b.faucetWallet.GetBalance()
		v, d := b.store.GetDistribution()
		msg := p.Sprintf("Available faucet balance is %.4f PACs\n", balance.Available)
		msg += p.Sprintf("A total of %.4f PACs has been distributed to %d validators\n", d, v)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
		return
	}

	if strings.Contains(m.Content, "faucet") {
		trimmedPrixix := strings.TrimPrefix(m.Content, "faucet")
		// faucet message must contain address/pubkey
		trimmedAddress := strings.Trim(trimmedPrixix, " ")
		peerID, pubKey, isValid, msg := b.validateInfo(trimmedAddress, m.Author.ID)

		if !isValid {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			return
		}

		if pubKey != "" {
			// check available balance
			balance := b.faucetWallet.GetBalance()
			if balance.Available < b.cfg.FaucetAmount {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "Insufficient faucet balance. Try again later.", m.Reference())
				return
			}

			// send faucet
			txHash := b.faucetWallet.BondTransaction(pubKey, trimmedAddress, b.cfg.FaucetAmount)
			if txHash != "" {
				err := b.store.SetData(peerID, trimmedAddress, m.Author.Username, m.Author.ID, b.cfg.FaucetAmount)
				if err != nil {
					log.Printf("error saving faucet information: %v\n", err)
				}
				msg := p.Sprintf("Faucet ( %.4f PACs) is staked on your node successfully!", b.cfg.FaucetAmount)
				_, _ = s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
			}
		}
	}
}

// help sends a message detailing how to use the bot discord-client side
// nolint
func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Pactus Universal Robot",
		URL:   "https://pactus.org",
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://pactus.org",
			IconURL: s.State.User.AvatarURL(""),
			Name:    s.State.User.Username,
		},
		Description: "Pactus Universal Robot is a robot that provides support and information about the Pactus Blockchain.\n" +
			"To see the faucet account balance, simply type: `balance`\n" +
			"To see the faucet address, simply type: `address`\n" +
			"To get network information, simply type: `network`\n" +
			"To request faucet for test network: simply post `faucet [validator address]`.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Example of requesting `faucet` ",
				Value: "faucet tpc1pxl333elgnrdtk0kjpjdvky44yu62x0cwupnpjl",
			},
		},
	})
}

func (b *Bot) validateInfo(address, discordID string) (string, string, bool, string) {
	_, err := crypto.AddressFromString(address)
	if err != nil {
		log.Printf("invalid address")
		return "", "", false, "Pactus Universal Robot is unable to handle your request." +
			" If you are requesting testing faucet, supply the valid address."
	}

	// check if the user is existing
	v, exists := b.store.FindDiscordID(discordID)
	if exists {
		return "", "", false, "Sorry. You already received faucet using this address: " + v.ValidatorAddress
	}

	cl, err := client.NewClient(b.cfg.Server)
	if err != nil {
		log.Printf("error establishing connection")
		return "", "", false, "The bot cannot establish connection to the blochain network. Try again later."
	}
	defer cl.Close()

	// check if the address exists in the list of validators
	isValidator := cl.IsValidator(address)
	if isValidator {
		return "", "", false, "Sorry. Your address is in the list of active validators. You do not need faucet gain."
	}

	peerInfo, pub, err := cl.GetPeerInfo(address)
	if err != nil || pub == nil {
		log.Printf("error getting peer info")
		return "", "", false, "Your node information could not obtained." +
			" Make sure your node is fully synced before requesting the faucet."
	}

	// check if the validator has already been given the faucet
	peerID, err := peer.IDFromBytes(peerInfo.PeerId)
	if err != nil || peerID.String() == "" {
		log.Printf("error getting peer id")
		return "", "", false, "Your node information could not obtained." +
			" Make sure your node is fully synced before requesting the faucet."
	}
	v, exists = b.store.GetData(peerID.String())
	if exists || v != nil {
		return "", "", false, "Sorry. You already received faucet using this address: " + v.ValidatorAddress
	}

	// check block height
	// height, err := cl.GetBlockchainHeight()
	// if err != nil {
	// 	log.Printf("error current block height")
	// 	return "", "", false, "The bot cannot establish connection to the blochain network. Try again later."
	// }
	// if (height - peerInfo.Height) > 1080 {
	//	msg := fmt.Sprintf("Your node is not fully synchronised. It is is behind by %v blocks." +
	//		" Make sure that your node is fully synchronised before requesting faucet.", (height - peerInfo.Height))

	// 	log.Printf("peer %s with address %v is not well synced: ", peerInfo.PeerId, address)
	// 	return "", "", false, msg
	// }
	return peerID.String(), pub.String(), true, ""
}

func (b *Bot) networkInfo() string {
	msg := "Pactus is truly decentralised proof of stake blockchain."
	cl, err := client.NewClient(b.cfg.Server)
	if err != nil {
		log.Printf("error establishing connection")
		return msg
	}
	defer cl.Close()
	nodes, err := cl.GetNetworkInfo()
	if err != nil {
		log.Printf("error establishing connection")
		return msg
	}
	msg += "\nThe following are the currentl statistics:\n"
	msg += fmt.Sprintf("Network started at : %v\n", time.UnixMilli(nodes.StartedAt*1000).Format("02/01/2006, 15:04:05"))
	msg += fmt.Sprintf("Total bytes sent : %v\n", nodes.TotalSentBytes)
	msg += fmt.Sprintf("Total received bytes : %v\n", nodes.TotalReceivedBytes)
	msg += fmt.Sprintf("Number of peer nodes: %v\n", len(nodes.Peers))
	// check block height
	blochainInfo, err := cl.GetBlockchainInfo()
	if err != nil {
		log.Printf("error current block height")
		return msg
	}
	msg += fmt.Sprintf("Block height: %v\n", blochainInfo.LastBlockHeight)
	msg += fmt.Sprintf("Total power: %.4f PACs\n", util.ChangeToCoin(blochainInfo.TotalPower))
	msg += fmt.Sprintf("Total committee power: %.4f PACs\n", util.ChangeToCoin(blochainInfo.CommitteePower))
	msg += fmt.Sprintf("Total validators: %v\n", blochainInfo.TotalValidators)
	return msg
}
