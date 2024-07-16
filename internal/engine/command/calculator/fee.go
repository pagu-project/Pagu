package calculator

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
)

func (bc *Calculator) calcFeeHandler(cmd *command.Command,
	_ entity.AppID, _ string, args map[string]string,
) command.CommandResult {
	amt, err := amount.FromString(args["amount"])
	if err != nil {
		return cmd.ErrorResult(errors.New("invalid amount param"))
	}

	fee, err := bc.clientMgr.GetFee(amt.ToNanoPAC())
	if err != nil {
		return cmd.ErrorResult(err)
	}

	feeAmount := amount.Amount(fee)

	return cmd.SuccessfulResult("Sending %s will cost %s with current fee percentage."+
		"\n> Note: Consider unbond and sortition transaction fee is 0 PAC always.", amt.String(), feeAmount.String())
}
