package discord

import (
	"fmt"
	"log"
	"pactus-bot/config"
	"pactus-bot/wallet"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	bot := &Bot{cfg: cfg, discordSession: dg, faucetWallet: w, store: ss}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(bot.messageHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

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

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "balance" reply with "available faucet balance"
	if m.Content == "balance" {
		b := b.faucetWallet.GetBalance()
		msg := fmt.Sprintf("Available faucet balance is %.6f PAC", b.Available)
		s.ChannelMessageSend(m.ChannelID, msg)
	}

	// faucet message must contain address/pubkey
	trimedContent := strings.Trim(m.Content, " ")
	subConents := strings.Split(trimedContent, "/")

	if len(subConents) == 2 && wallet.IsValidData(subConents[0], subConents[1]) {
		address := subConents[0]
		pubKey := subConents[1]
		// check if the validator has already been given the faucet
		_, exists := b.store.GetData(m.Author.ID)
		if exists {
			s.ChannelMessageSendReply(m.ChannelID, "You received the faucet!", m.MessageReference)
			return
		}

		//check available balance
		balance := b.faucetWallet.GetBalance()
		if balance.Available < b.cfg.FaucetAmount {
			s.ChannelMessageSend(m.ChannelID, "Insuffcient faucet balance. Try again later")
			return
		}

		//send faucet
		txHash := b.faucetWallet.BondTransaction(pubKey, address, b.cfg.FaucetAmount)
		if txHash != "" {
			err := b.store.SetData(address, m.Author.Username, m.Author.ID, b.cfg.FaucetAmount)
			if err != nil {
				log.Printf("error saving faucet information: %v\n", err)
			}
			msg := fmt.Sprintf("Faucet ( %.6f PAC) is transfered successfully!", b.cfg.FaucetAmount)
			s.ChannelMessageSendReply(m.ChannelID, msg, m.MessageReference)
		}

	}

}
