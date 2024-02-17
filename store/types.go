package store

type Claimer struct {
	DiscordID   string `json:"did"`
	TotalReward int64  `json:"r"`
	ClaimedTxID string `json:"tx_id"`
}

type TwitterParty struct {
	TwitterID            string `json:"twitter_id"`
	TwitterName          string `json:"twitter_name"`
	RetweetID            string `json:"retweet_id"`
	ValAddr              string `json:"val_addr"`
	ValPubKey            string `json:"val_pub"`
	DiscordID            string `json:"discord_id"`
	DiscountCode         string `json:"discount_code"`
	TotalPrice           int    `json:"total_price"`
	AmountInPAC          int64  `json:"amount_in_pac"`
	CreatedAt            int64  `json:"created_at"`
	NowPaymentsInvoiceID string `json:"nowpayments_id"`
	NowPaymentsFinished  bool   `json:"nowpayments_finished"`
	TransactionID        string `json:"tx_id"`
}

type WhitelistInfo struct {
	TwitterID     string `json:"twitter_id"`
	TwitterName   string `json:"twitter_name"`
	WhitelistedBy string `json:"whitelisted_by"`
}

type BoosterStatus struct {
	Pac            int
	Usdt           int
	AllPkgs        int
	ClaimedPkgs    int
	UnClaimedPkgs  int
	PaymentDone    int
	PaymentWaiting int
	Whitelists     int
}

type ClaimStatus struct {
	Claimed          int
	ClaimedAmount    int64
	NotClaimed       int
	NotClaimedAmount int64
}

func (c *Claimer) IsClaimed() bool {
	return c.ClaimedTxID != ""
}
