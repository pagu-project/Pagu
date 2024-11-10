package zealy

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (z *Zealy) claimHandler(caller *entity.User, cmd *command.Command, args map[string]string) command.CommandResult {
	user, err := z.db.GetZealyUser(caller.CallerID)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if user.IsClaimed() {
		return cmd.FailedResult("You already claimed your reward: https://pacviewer.com/transaction/%s",
			user.TxHash)
	}

	address := args["address"]
	txHash, err := z.wallet.TransferTransaction(address, "Pagu Zealy reward distribution", user.Amount)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = z.db.UpdateZealyUser(caller.CallerID, txHash); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Zealy reward claimed successfully: https://pacviewer.com/transaction/%s",
		txHash)
}
