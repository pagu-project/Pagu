package blockchain

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/engine/command"
)

func (bc *Blockchain) calcFeeHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	amt, err := amount.FromString(args[0])
	if err != nil {
		return cmd.ErrorResult(err)
	}

	fee, err := bc.clientMgr.GetFee(int64(amt))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	calcedFee := amount.Amount(fee)

	return cmd.SuccessfulResult("Sending %s will cost %s with current fee percentage."+
		"\n> Note: Consider unbond and sortition transaction fee is 0 PAC always.", amt, calcedFee.String())
}
