package voucher

import (
	"strconv"

	"github.com/pagu-project/Pagu/pkg/utils"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (v *Voucher) createHandler(cmd command.Command, _ entity.AppID, _ string, args ...string) command.CommandResult {
	//! Admin only check

	code := utils.RandomString(8, utils.CapitalLetterNumbers)
	for _, err := v.db.GetVoucherByCode(code); err == nil; {
		code = utils.RandomString(8, utils.CapitalLetterNumbers)
	}

	amount := args[0]
	intAmount, err := strconv.Atoi(amount)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	validMonths := args[1]
	expireMonths, err := strconv.Atoi(validMonths)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	vch := &entity.Voucher{
		Creator: cmd.User.ID,
		Code:    code,

		ValidMonths: uint8(expireMonths),
		Amount:      uint(intAmount),
	}

	if len(args) > 2 {
		vch.Recipient = args[2]
	}
	if len(args) > 3 {
		vch.Desc = args[3]
	}

	err = v.db.AddVoucher(vch)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher crated successfully!")
}
