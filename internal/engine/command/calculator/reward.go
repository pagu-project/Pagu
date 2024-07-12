package calculator

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/utils"
)

func (bc *Calculator) calcRewardHandler(cmd *command.Command,
	_ entity.AppID, _ string, args map[string]any,
) command.CommandResult {
	stake, ok := args["stake"].(int)
	if !ok {
		return cmd.ErrorResult(errors.New("invalid stake param"))
	}

	if stake < 1 || stake > 1_000 {
		return cmd.ErrorResult(
			fmt.Errorf("%v is invalid amount; minimum stake amount is 1 PAC and maximum is 1,000 PAC", stake))
	}

	numOfDays, ok := args["days"].(int)
	if !ok {
		return cmd.ErrorResult(errors.New("invalid days param"))
	}

	if numOfDays < 1 || numOfDays > 365 {
		return cmd.ErrorResult(fmt.Errorf("%v is invalid time; minimum time value is 1 and maximum is 365", numOfDays))
	}

	blocks := numOfDays * 8640
	bi, err := bc.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	reward := int64(stake*blocks) / int64(amount.Amount(bi.TotalPower).ToPAC())

	return cmd.SuccessfulResult("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator "+
		"in %d days ⏰ with %s total power ⚡ of committee."+
		"\n\n> Note📝: This number is just an estimation. "+
		"It will vary depending on your stake amount and total network power.",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), numOfDays,
		utils.FormatNumber(int64(amount.Amount(bi.TotalPower).ToPAC())))
}
