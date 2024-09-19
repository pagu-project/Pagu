package entity

type Price struct {
	XeggexPacToUSDT XeggexPriceResponse
	AzbitPacToUSDT  AzbitPriceResponse
}

type XeggexPriceResponse struct {
	LastPrice      string  `json:"lastPrice"`
	YesterdayPrice string  `json:"yesterdayPrice"`
	HighPrice      string  `json:"highPrice"`
	LowPrice       string  `json:"lowPrice"`
	Volume         string  `json:"volume"`
	Decimal        int     `json:"priceDecimals"`
	BestAsk        string  `json:"bestAsk"`
	BestBid        string  `json:"bestBid"`
	SpreadPercent  string  `json:"spreadPercent"`
	ChangePercent  string  `json:"changePercent"`
	MarketCap      float64 `json:"marketcapNumber"`
}

type AzbitPriceResponse struct {
	Timestamp                int     `json:"timestamp"`
	CurrencyPairCode         string  `json:"currencyPairCode"`
	Price                    float64 `json:"price"`
	Price24HAgo              float64 `json:"price24hAgo"`
	PriceChangePercentage24H float64 `json:"priceChangePercentage24h"`
	Volume24H                float64 `json:"volume24h"`
	BidPrice                 float64 `json:"bidPrice"`
	AskPrice                 float64 `json:"askPrice"`
	Low24H                   float64 `json:"low24h"`
	High24H                  float64 `json:"high24h"`
}
