package zealy

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
)

func (z *Zealy) claimHandler(cmd *command.Command,
	_ entity.AppID, callerID string, args map[string]any,
) command.CommandResult {
	user, err := z.db.GetZealyUser(callerID)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if user.IsClaimed() {
		return cmd.FailedResult("You already claimed your reward: https://pacviewer.com/transaction/%s",
			user.TxHash)
	}

	address, ok := args["address"].(string)
	if !ok {
		return cmd.ErrorResult(errors.New("invalid address param"))
	}

	amt, _ := amount.NewAmount(float64(user.Amount))
	txHash, err := z.wallet.TransferTransaction(address, "Pagu Zealy reward distribution", amt)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = z.db.UpdateZealyUser(callerID, txHash); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Zealy reward claimed successfully: https://pacviewer.com/transaction/%s",
		txHash)
}
