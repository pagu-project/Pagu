package command

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/entity"
)

func (h *MiddlewareHandler) WalletBalance(_ *entity.User, _ *Command, _ map[string]string) error {
	if h.wallet.Balance() < 5 {
		return errors.New("the Pagu Wallet balance is less than 5 PAC")
	}

	return nil
}
