package voucher

import (
	"strconv"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (v *Voucher) createHandler(cmd command.Command, _ entity.AppID, callerID string, args ...string) command.CommandResult {
	//! Admin only check

	cID, err := strconv.Atoi(callerID)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	recipient := args[0]
	description := args[1]
	validMonths := args[2]
	amount := args[3]
	discordID := args[4]
	code := args[5]

	expireMonths, err := strconv.Atoi(validMonths)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	intAmount, err := strconv.Atoi(amount)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	err = v.db.AddVoucher(&entity.Voucher{
		Creator:     uint(cID),
		Code:        code,
		Desc:        description,
		DiscordID:   discordID,
		Recipient:   recipient,
		ValidMonths: uint(expireMonths),
		Amount:      uint(intAmount),
	})
	if err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher crated successfully!")
}
