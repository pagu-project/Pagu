package zealy

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (z *Zealy) claimHandler(cmd command.Command, _ entity.AppID, callerID string, args ...string) command.CommandResult {
	user, err := z.db.GetZealyUser(callerID)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if user.IsClaimed() {
		return cmd.FailedResult("You already claimed your reward: https://pacviewer.com/transaction/%s",
			user.TxHash)
	}

	address := args[0]
	txHash, err := z.wallet.TransferTransaction(address, "PaGu Zealy reward distribution", int64(user.Amount))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = z.db.UpdateZealyUser(callerID, txHash); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Zealy reward claimed successfully: https://pacviewer.com/transaction/%s",
		txHash)
}
