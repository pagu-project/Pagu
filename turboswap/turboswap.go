package turboswap

import (
	"bytes"
	"context"
	"encoding/json"
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

func (ts *Turboswap) GetStatus(ctx context.Context, party *store.TwitterParty) (*DiscountStatus, error) {
	url := fmt.Sprintf("https://swap-api.sensifia.vc/pactus/discount/status/%v/%v", party.ValPubKey, ts.APIToken)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}

	logger.Info("calling Turboswap/status", "twitter", party.TwitterName)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// http.StatusOK = 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to call Turboswap/status. Status code: %v", resp.StatusCode)
	}

	buf := make([]byte, 0, 1024)
	_, err = resp.Body.Read(buf)
	if err != nil {
		return nil, err
	}
	logger.Info("Turboswap call successful", "res", string(buf))

	res := &DiscountStatus{}
	err = json.Unmarshal(buf, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ts *Turboswap) SendDiscountCode(ctx context.Context, party *store.TwitterParty) error {
	url := "https://swap-api.sensifia.vc/pactus/discount"
	jsonStr := fmt.Sprintf(`{"api_key":"%v","code":"%v","validator_public_key":"%v","total_coins":"%v","total_price_in_usd":"%v"}`,
		ts.APIToken, party.DiscountCode, party.ValPubKey, party.AmountInPAC, party.TotalPrice)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	logger.Info("calling Turboswap/discount", "twitter", party.TwitterName)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// http.StatusOK = 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to call Turboswap/discount. Status code: %v", resp.StatusCode)
	}

	return nil
}
