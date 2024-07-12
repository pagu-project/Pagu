package voucher

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/pagu-project/Pagu/pkg/utils"
)

func (v *Voucher) createHandler(cmd *command.Command, _ entity.AppID,
	_ string, args map[string]any,
) command.CommandResult {
	code := utils.RandomString(8, utils.CapitalAlphanumerical)
	for _, err := v.db.GetVoucherByCode(code); err == nil; {
		code = utils.RandomString(8, utils.CapitalAlphanumerical)
	}

	amountStr, ok := args["amount"].(float64)
	if !ok {
		return cmd.ErrorResult(errors.New("invalid amount param"))
	}

	amt, err := amount.NewAmount(amountStr)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	maxStake, _ := amount.NewAmount(1000)
	if amt > maxStake {
		return cmd.ErrorResult(errors.New("stake amount is more than 1000"))
	}

	expireMonths, ok := args["valid-months"].(int)
	if !ok {
		return cmd.ErrorResult(errors.New("invalid valid-months param"))
	}

	vch := &entity.Voucher{
		Creator:     cmd.User.ID,
		Code:        code,
		ValidMonths: uint8(expireMonths),
		Amount:      amt,
	}

	if args["recipient"] != nil {
		vch.Recipient, ok = args["recipient"].(string)
		if !ok {
			return cmd.ErrorResult(errors.New("invalid recipient param"))
		}
	}

	if args["description"] != nil {
		vch.Desc, ok = args["description"].(string)
		if !ok {
			return cmd.ErrorResult(errors.New("invalid description param"))
		}
	}

	err = v.db.AddVoucher(vch)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher created successfully! \n Code: %s", vch.Code)
}
