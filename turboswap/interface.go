package turboswap

import (
	"context"

	"github.com/kehiy/RoboPac/store"
)

type DiscountStatus struct {
	Status          string `json:"status"`
	TransactionHash string `json:"transactionHash"`
}
type ITurboSwap interface {
	GetStatus(ctx context.Context, party *store.TwitterParty) (*DiscountStatus, error)
	SendDiscountCode(ctx context.Context, party *store.TwitterParty) error
}
