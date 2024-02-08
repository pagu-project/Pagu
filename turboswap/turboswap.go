package turboswap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kehiy/RoboPac/store"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type Turboswap struct {
	APIKey string
}

func NewTurboswap(apiKey string) (*Turboswap, error) {

	return &Turboswap{
		APIKey: apiKey,
	}, nil
}

func (ts *Turboswap) GetDiscountCodeStatus(ctx context.Context, pubKey string) error {

	return nil
}

func (ts *Turboswap) AddDiscountCode(ctx context.Context, party *store.TwitterParty) error {

	return nil
}
