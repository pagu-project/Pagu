package store

type ClaimTransaction struct {
	TxID   string  `json:"transaction_id"`
	Amount float64 `json:"amount"`
	Time   int64   `json:"time"`
}

type Claimer struct {
	DiscordID string `json:"did"`
	// ValAddr          string            `json:"main_net_validator_address"`
	TotalReward      float64           `json:"r"`
	ClaimTransaction *ClaimTransaction `json:"claim_transaction"`
}

func (c *Claimer) IsClaimed() bool {
	return c.ClaimTransaction != nil
}

type IStore interface {
	ClaimerInfo(testNetValAddr string) *Claimer
	AddClaimTransaction(amount float64, time int64, txID, discordID, testNetValAddr string) error
}
