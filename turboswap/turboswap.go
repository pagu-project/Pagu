package turboswap

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/kehiy/RoboPac/store"
)

type Turboswap struct {
	APIToken string
}

func NewTurboswap(apiToken string) (*Turboswap, error) {
	return &Turboswap{
		APIToken: apiToken,
	}, nil
}

func (ts *Turboswap) GetDiscountCodeStatus(ctx context.Context, pubKey string) error {

	return nil
}

func (ts *Turboswap) AddDiscountCode(ctx context.Context, party *store.TwitterParty) error {
	url := "https://swap-api.sensifia.vc/pactus/discount"
	jsonStr := fmt.Sprintf(`{"apiKey":"%v","code":"%v","validatorPublicKey":"%v","priceInCents":"%v"}`,
		ts.APIToken, party.DiscountCode, party.ValPubKey, party.UnitPrice)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
