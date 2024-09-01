package entity

type Price struct {
	XeggexPacToUSDT   XeggexPriceResponse
	ExbitronPacToUSDT ExbitronPriceResponse
	P2BPacToUSDT      P2BPriceResponse
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

type P2BPriceResponse struct {
	LastPrice string `json:"last"`
}

type ExbitronPriceResponse []struct {
	TickerID       string `json:"ticker_id"`
	BaseCurrency   string `json:"base_currency"`
	TargetCurrency string `json:"target_currency"`
	LastPrice      string `json:"last_price"`
	BaseVolume     string `json:"base_volume"`
	TargetVolume   string `json:"target_volume"`
	Bid            string `json:"bid"`
	Ask            string `json:"ask"`
	High           string `json:"high"`
	Low            string `json:"low"`
}

type ExbitronTicker struct {
	TickerID       string `json:"ticker_id"`
	BaseCurrency   string `json:"base_currency"`
	TargetCurrency string `json:"target_currency"`
	LastPrice      string `json:"last_price"`
	BaseVolume     string `json:"base_volume"`
	TargetVolume   string `json:"target_volume"`
	Bid            string `json:"bid"`
	Ask            string `json:"ask"`
	High           string `json:"high"`
	Low            string `json:"low"`
}

func (e ExbitronPriceResponse) GetPacToUSDT() ExbitronTicker {
	const tickerID = "PAC-USDT"

	for index := range e {
		ticker := e[index]
		if ticker.TickerID == tickerID {
			return ExbitronTicker{
				TickerID:       tickerID,
				BaseCurrency:   ticker.BaseCurrency,
				TargetCurrency: ticker.TargetCurrency,
				LastPrice:      ticker.LastPrice,
				BaseVolume:     ticker.BaseVolume,
				TargetVolume:   ticker.TargetVolume,
				Bid:            ticker.Bid,
				Ask:            ticker.Ask,
				High:           ticker.High,
				Low:            ticker.Low,
			}
		}
	}

	return ExbitronTicker{}
}
