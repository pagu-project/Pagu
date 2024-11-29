package zealy

import (
	"fmt"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/log"
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
		log.Error("error in transfer zealy reward", "err", err)
		transferErr := fmt.Errorf("Failed to transfer zealy reward. Please make sure the address is valid") //nolint
		return cmd.ErrorResult(transferErr)
	}

	if err = z.db.UpdateZealyUser(caller.CallerID, txHash); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Zealy reward claimed successfully: https://pacviewer.com/transaction/%s",
		txHash)
}
