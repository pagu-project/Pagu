package store

import "time"

type ClaimTransaction struct {
	TxID   string `json:"transaction_id"`
	Data   string `json:"data"`
	Amount int64  `json:"amount"`
	Time   int64  `json:"time"`
}

type Claimer struct {
	DiscordID        string            `json:"discord_id"`
	TotalReward      int64             `json:"total_reward"`
	ClaimTransaction *ClaimTransaction `json:"claim_transaction"`
}

func (c *Claimer) IsClaimed() bool {
	return c.ClaimTransaction != nil
}

type IStore interface {
	ClaimerInfo(discordID string) *Claimer
	AddClaimTransaction(TxID string, Amount int64, Time time.Time, Data string, discordID string) error
}
