package turboswap

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/kehiy/RoboPac/store"
	"github.com/pactus-project/pactus/util/logger"
)

type Turboswap struct {
	APIToken string
}

func NewTurboswap(apiToken string) (*Turboswap, error) {
	return &Turboswap{
		APIToken: apiToken,
	}, nil
}

func (ts *Turboswap) GetStatus(ctx context.Context, party *store.TwitterParty) error {
	url := fmt.Sprintf("https://swap-api.sensifia.vc/pactus/discount/status/%v/%v", party.ValPubKey, ts.APIToken)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}
	client := &http.Client{}

	logger.Info("calling swap-api/status", "twitter", party.TwitterName)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := make([]byte, 0, 1024)
	_, err = resp.Body.Read(buf)
	if err != nil {
		return err
	}

	logger.Info("response from turboswap", "res", string(buf))
	return nil
}

func (ts *Turboswap) SendDiscountCode(ctx context.Context, party *store.TwitterParty) error {
	url := "https://swap-api.sensifia.vc/pactus/discount"
	jsonStr := fmt.Sprintf(`{"api_key":"%v","code":"%v","validator_public_key":"%v","total_coins":"%v","total_price_in_usd":"%v","created_at":"%v"}`,
		ts.APIToken, party.DiscountCode, party.ValPubKey, party.AmountInPAC, party.TotalPrice, party.CreatedAt)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	logger.Info("calling swap-api/discount", "twitter", party.TwitterName)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := make([]byte, 0, 1024)
	_, err = resp.Body.Read(buf)
	if err != nil {
		return err
	}
	logger.Info("response from turboswap", "res", string(buf))
	return nil
}
