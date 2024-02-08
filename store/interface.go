package store

type Claimer struct {
	DiscordID   string `json:"did"`
	TotalReward int64  `json:"r"`
	ClaimedTxID string `json:"tx_id"`
}

type TwitterParty struct {
	TwitterName   string `json:"twitter"`
	TweetID       string `json:"tweet"`
	PricePerCents int    `json:"price"`
	ValAddr       string `json:"val_addr"`
	ValPubKey     string `json:"val_pub"`
	AmountInPAC   int    `json:"amount"`
	DiscountCode  string `json:"code"`
}

func (c *Claimer) IsClaimed() bool {
	return c.ClaimedTxID != ""
}

type IStore interface {
	ClaimerInfo(testNetValAddr string) *Claimer
	AddClaimTransaction(testNetValAddr string, txID string) error
	AddTwitterParty(party *TwitterParty) error
	GetTwitterParty(twitterName string) *TwitterParty
	Status() (int64, int64, int64, int64)
}
