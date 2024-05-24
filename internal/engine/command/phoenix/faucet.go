package phoenix

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/repository/faucet"
	"github.com/pagu-project/Pagu/internal/repository/user"
)

func (pt *Phoenix) faucetHandler(cmd command.Command, _ command.AppID, callerID string, args ...string) command.CommandResult {
	if !pt.db.HasUser(callerID) {
		if err := pt.db.AddUser(
			&user.User{
				ID: callerID,
			},
		); err != nil {
			return cmd.ErrorResult(err)
		}
	}

	if !pt.db.CanGetFaucet(callerID) {
		return cmd.FailedResult("Uh, you used your share of faucets today!")
	}

	if pt.wallet.Balance() < 5 {
		return cmd.FailedResult("RoboPac Phoenix wallet is empty, please contact the team!")
	}

	toAddr := args[0]
	txID, err := pt.wallet.TransferTransaction(toAddr, "Phoenix Testnet Pagu Faucet", 5) //! define me on config?
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = pt.db.AddFaucet(&faucet.Faucet{
		Address:         toAddr,
		Amount:          5,
		TransactionHash: txID,
		UserID:          callerID,
	}); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("You got %d tPAC in %s address on Phoenix Testnet!", 5, toAddr)
}
