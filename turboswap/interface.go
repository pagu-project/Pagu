package turboswap

import (
	"context"

	"github.com/kehiy/RoboPac/store"
)

type ITurboSwap interface {
	GetStatus(ctx context.Context, party *store.TwitterParty) error
	SendDiscountCode(ctx context.Context, party *store.TwitterParty) error
}
