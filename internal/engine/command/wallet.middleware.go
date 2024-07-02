package command

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/entity"
)

func (h *MiddlewareHandler) WalletBalance(_ *Command, _ entity.AppID, _ string, _ ...string) error {
	if h.wallet.Balance() < 5 {
		return errors.New("empty pagu wallet balance")
	}

	return nil
}
