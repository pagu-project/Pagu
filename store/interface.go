package store

type Claimer struct {
	DiscordID   string `json:"did"`
	TotalReward int64  `json:"r"`
	ClaimedTxID string `json:"tx_id"`
}

type TwitterParty struct {
	TwitterID    string `json:"twitter_id"`
	TwitterName  string `json:"twitter_name"`
	RetweetID    string `json:"retweet_id"`
	ValAddr      string `json:"val_addr"`
	ValPubKey    string `json:"val_pub"`
	DiscountCode string `json:"discount_code"`
	UnitPrice    int    `json:"unit_price_in_cents"`
	TotalPrice   int    `json:"total_price"`
	AmountInPAC  int    `json:"amount_in_pac"`
	CreatedAt    int64  `json:"created_at"`
}

func (c *Claimer) IsClaimed() bool {
	return c.ClaimedTxID != ""
}

type IStore interface {
	ClaimerInfo(testNetValAddr string) *Claimer
	AddClaimTransaction(testNetValAddr string, txID string) error
	AddTwitterParty(party *TwitterParty) error
	FindTwitterParty(twitterName string) *TwitterParty
	Status() (int64, int64, int64, int64)
}
