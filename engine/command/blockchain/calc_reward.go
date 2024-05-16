package blockchain

import (
	"fmt"
	"strconv"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/utils"

	"github.com/pagu-project/Pagu/engine/command"
)

func (bc *Blockchain) calcRewardHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	stake, err := strconv.Atoi(args[0])
	if err != nil {
		return cmd.ErrorResult(err)
	}

	time := args[1]

	if stake < 1 || stake > 1_000 {
		return cmd.ErrorResult(fmt.Errorf("%v is invalid amount; minimum stake amount is 1 PAC and maximum is 1,000 PAC", stake))
	}

	var blocks int
	switch time {
	case "day":
		blocks = 8640
	case "month":
		blocks = 259200
	case "year":
		blocks = 3110400
	default:
		blocks = 8640
		time = "day"
	}

	bi, err := bc.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	reward := int64(stake*blocks) / int64(amount.Amount(bi.TotalPower).ToPAC())

	return cmd.SuccessfulResult("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator in one %s ⏰ with %s total power ⚡ of committee."+
		"\n\n> Note📝: This number is just an estimation. It will vary depending on your stake amount and total network power.",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), time, utils.FormatNumber(int64(amount.Amount(bi.TotalPower).ToPAC())))
}
