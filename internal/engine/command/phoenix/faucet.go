package phoenix

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (pt *Phoenix) faucetHandler(cmd command.Command, _ entity.AppID, _ string, args ...string) command.CommandResult {
	if len(args) == 0 {
		return cmd.ErrorResult(errors.New("invalid wallet address"))
	}

	toAddr := args[0]
	if len(toAddr) != 43 || toAddr[:3] != "tpc" {
		return cmd.ErrorResult(errors.New("invalid wallet address"))
	}

	if !pt.db.CanGetFaucet(cmd.User) {
		return cmd.FailedResult("Uh, you used your share of faucets today!")
	}

	txID, err := pt.wallet.TransferTransaction(toAddr, "Phoenix Testnet Pagu PhoenixFaucet", int64(pt.faucetAmount))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = pt.db.AddFaucet(&entity.PhoenixFaucet{
		UserID:          cmd.User.ID,
		Address:         toAddr,
		Amount:          pt.faucetAmount,
		TransactionHash: txID,
	}); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("You got %d tPAC in %s address on Phoenix Testnet!", pt.faucetAmount, toAddr)
}
