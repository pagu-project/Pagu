package store

type Claimer struct {
	DiscordID   string `json:"did"`
	TotalReward int64  `json:"r"`
	ClaimedTxID string `json:"tx_id"`
}

func (c *Claimer) IsClaimed() bool {
	return c.ClaimedTxID != ""
}

type IStore interface {
	ClaimerInfo(testNetValAddr string) *Claimer
	AddClaimTransaction(testNetValAddr string, txID string) error
	Status() (int64, int64, int64, int64)
}
