package voucher

import (
	"errors"
	"strconv"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/pagu-project/Pagu/pkg/utils"
)

func (v *Voucher) createHandler(cmd *command.Command, _ entity.AppID, _ string, args ...string) command.CommandResult {
	code := utils.RandomString(8, utils.CapitalAlphanumerical)
	for _, err := v.db.GetVoucherByCode(code); err == nil; {
		code = utils.RandomString(8, utils.CapitalAlphanumerical)
	}

	amountStr := args[0]
	amt, err := amount.FromString(amountStr)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	maxStake, _ := amount.NewAmount(1000)
	if amt > maxStake {
		return cmd.ErrorResult(errors.New("stake amount is more than 1000"))
	}

	validMonths := args[1]
	expireMonths, err := strconv.Atoi(validMonths)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	vch := &entity.Voucher{
		Creator:     cmd.User.ID,
		Code:        code,
		ValidMonths: uint8(expireMonths),
		Amount:      amt,
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

	return cmd.SuccessfulResult("Voucher created successfully! \n Code: %s", vch.Code)
}
