package command

import (
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

type MiddlewareFunc func(caller *entity.User, cmd *Command, args map[string]string) error

type MiddlewareHandler struct {
	db     repository.Database
	wallet wallet.IWallet
}

func NewMiddlewareHandler(d repository.Database, w wallet.IWallet) *MiddlewareHandler {
	return &MiddlewareHandler{
		db:     d,
		wallet: w,
	}
}
