package calculator

import (
	"errors"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (bc *Calculator) calcFeeHandler(cmd *command.Command,
	_ entity.AppID, _ string, args map[string]any,
) command.CommandResult {
	amt, ok := args["amount"].(int)
	if !ok {
		return cmd.ErrorResult(errors.New("invalid amount param"))
	}

	fee, err := bc.clientMgr.GetFee(int64(amt))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	feeAmount := amount.Amount(fee)

	return cmd.SuccessfulResult("Sending %d will cost %s with current fee percentage."+
		"\n> Note: Consider unbond and sortition transaction fee is 0 PAC always.", amt, feeAmount.String())
}
