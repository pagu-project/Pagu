package turboswap

import (
	"context"

	"github.com/kehiy/RoboPac/store"
)

type ITurboSwap interface {
	GetStatus(ctx context.Context, pubKey string) error
	SendDiscountCode(ctx context.Context, party *store.TwitterParty) error
}
