package command

import (
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

type MiddlewareFunc func(cmd Command, appID entity.AppID, callerID string, args ...string) error

type MiddlewareHandler struct {
	db     *repository.DB
	wallet *wallet.Wallet
}

func NewMiddlewareHandler(d *repository.DB, w *wallet.Wallet) *MiddlewareHandler {
	return &MiddlewareHandler{
		db:     d,
		wallet: w,
	}
}
