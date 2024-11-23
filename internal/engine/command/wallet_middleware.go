package command

import (
	"fmt"

	"github.com/pagu-project/Pagu/internal/entity"
)

const minWalletBalance = 500

func (h *MiddlewareHandler) WalletBalance(_ *entity.User, _ *Command, _ map[string]string) error {
	if h.wallet.Balance() < minWalletBalance {
		return fmt.Errorf("the Pagu Wallet balance is less than %d PAC", minWalletBalance)
	}

	return nil
}
