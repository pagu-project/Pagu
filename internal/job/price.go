package job

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/cache"
	"github.com/pagu-project/Pagu/pkg/log"
)

const (
	_defaultXeggexPriceEndpoint   = "https://api.xeggex.com/api/v2/market/getbysymbol/Pactus%2Fusdt"
	_defaultExbitronPriceEndpoint = "https://api.exbitron.digital/api/v1/cg/tickers"
	_defaultP2BPriceEndpoint      = "https://api.p2pb2b.com/api/v2/all/ticker?market=PAC_USDT"
)

var (
	P2BAPIKey    string
	P2BSecretKey string
)

type price struct {
	ctx    context.Context
	cache  cache.Cache[string, entity.Price]
	ticker *time.Ticker
	cancel context.CancelFunc
}

func NewPrice(
	cch cache.Cache[string, entity.Price],
	p2bApiKey, p2bSecretKey string,
) Job {
	ctx, cancel := context.WithCancel(context.Background())
	P2BAPIKey = p2bApiKey
	P2BSecretKey = p2bSecretKey

	return &price{
		cache:  cch,
		ticker: time.NewTicker(128 * time.Second),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *price) Start() {
	p.start()
	go p.runTicker()
}

func (p *price) start() {
	var (
		wg       sync.WaitGroup
		price    entity.Price
		xeggex   entity.XeggexPriceResponse
		exbitron entity.ExbitronPriceResponse
		p2b      entity.P2BPriceResponse
	)

	ctx := context.Background()

	// xeggex
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultXeggexPriceEndpoint, &xeggex); err != nil {
			log.Error(err.Error())
			return
		}
	}()

	// exbitron
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.getPrice(ctx, _defaultExbitronPriceEndpoint, &exbitron); err != nil {
			log.Error(err.Error())
		}
	}()

	// p2b
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getP2BPrice(ctx, P2BAPIKey, P2BSecretKey, &p2b); err != nil {
			log.Error(err.Error())
			return
		}
	}()

	wg.Wait()

	price.XeggexPacToUSDT = xeggex
	price.ExbitronPacToUSDT = exbitron
	price.P2BPacToUSDT = p2b

	ok := p.cache.Exists(config.PriceCacheKey)
	if ok {
		p.cache.Update(config.PriceCacheKey, price, 0)
	} else {
		p.cache.Add(config.PriceCacheKey, price, 0)
	}
}

func (p *price) runTicker() {
	for {
		select {
		case <-p.ctx.Done():
			return

		case <-p.ticker.C:
			p.start()
		}
	}
}

func (p *price) getPrice(ctx context.Context, endpoint string, priceResponse any) error {
	cli := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is %v", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(priceResponse)
}

func getP2BPrice(ctx context.Context, apiKey, secretKey string, p2bResponse *entity.P2BPriceResponse) error {
	payload := struct {
		Request string `json:"request"`
		Nonce   string `json:"nonce"`
	}{
		Request: "/api/v2/all/ticker",
		Nonce:   fmt.Sprintf("%d", time.Now().UnixMilli()),
	}

	payloadByte, _ := json.Marshal(payload)
	cli := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, _defaultP2BPriceEndpoint, bytes.NewBuffer(payloadByte))
	if err != nil {
		return err
	}

	base64Data := base64.StdEncoding.EncodeToString(payloadByte)
	h := hmac.New(sha512.New, []byte(secretKey))
	h.Write([]byte(base64Data))
	hmacData := h.Sum(nil)
	hmacHex := hex.EncodeToString(hmacData)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-TXC-APIKEY", apiKey)
	req.Header.Add("X-TXC-PAYLOAD", base64Data)
	req.Header.Add("X-TXC-SIGNATURE", hmacHex)

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("p2b request failed with status: %v", resp.StatusCode)
	}

	res := struct {
		Result    entity.P2BPriceResponse `json:"result"`
		Success   bool                    `json:"success"`
		ErrorCode string                  `json:"errorCode"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf("p2b response failed with code: %s", res.ErrorCode)
	}

	p2bResponse = &res.Result
	return nil
}

func (p *price) Stop() {
	p.ticker.Stop()
}
