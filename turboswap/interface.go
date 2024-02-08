package turboswap

import (
	"context"

	"github.com/kehiy/RoboPac/store"
)

type ITurboSwap interface {
	GetDiscountCodeStatus(ctx context.Context, pubKey string) error
	AddDiscountCode(ctx context.Context, party *store.TwitterParty) error
}
