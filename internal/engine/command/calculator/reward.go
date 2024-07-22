package calculator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/utils"
)

func (bc *Calculator) calcRewardHandler(
	_ *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	stake, err := amount.FromString(args["stake"])
	if err != nil {
		return cmd.ErrorResult(errors.New("invalid stake param"))
	}

	minStake, _ := amount.NewAmount(1)
	maxStake, _ := amount.NewAmount(1000)
	if stake < minStake || stake > maxStake {
		return cmd.ErrorResult(
			fmt.Errorf("%v is invalid amount; minimum stake amount is 1 PAC and maximum is 1,000 PAC", stake))
	}

	numOfDays, err := strconv.Atoi(args["days"])
	if err != nil {
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

	reward := (stake.ToNanoPAC() * int64(blocks)) / bi.TotalPower

	return cmd.SuccessfulResult("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator "+
		"in %d days ⏰ with %s total power ⚡ of committee."+
		"\n\n> Note📝: This number is just an estimation. "+
		"It will vary depending on your stake amount and total network power.",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), numOfDays,
		utils.FormatNumber(int64(amount.Amount(bi.TotalPower).ToPAC())))
}
