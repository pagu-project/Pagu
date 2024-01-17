package store

import "time"

type ClaimedTransaction struct {
	TxID   string
	Amount int64
	Time   time.Time
	Data   string
}

type Claimer struct {
	DiscordID          string // user ID
	TotalReward        int64
	ClaimedTransaction *ClaimedTransaction
}

func (c *Claimer) IsClaimed() bool {
	return c.ClaimedTransaction != nil
}

type IStore interface {
	ClaimerInfo(discordID string) *Claimer
	AddClaimTransaction(TxID string, Amount int64, Time time.Time, Data string) error
}
