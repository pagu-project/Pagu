package voucher

import (
	"errors"
	"time"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (v *Voucher) claimHandler(cmd command.Command, _ entity.AppID, _ string, args ...string) command.CommandResult {
	code := args[0]
	if len(code) != 8 {
		return cmd.ErrorResult(errors.New("voucher code is not valid"))
	}

	voucher, err := v.db.GetVoucherByCode(code)
	if err != nil {
		return cmd.ErrorResult(errors.New("voucher code is not valid"))
	}

	now := time.Now().Month()
	if voucher.CreatedAt.Month() >= (now + time.Month(voucher.ValidMonths)) {
		return cmd.ErrorResult(errors.New("voucher is expired"))
	}

	if len(voucher.TxHash) > 0 {
		return cmd.ErrorResult(errors.New("voucher code claimed before"))
	}

	address := args[1]
	validatorInfo, err := v.clientManager.GetValidatorInfo(address)
	if err != nil {
		return cmd.ErrorResult(errors.New("bond error"))
	}

	pubKey := validatorInfo.GetValidator().GetPublicKey()
	txHash, err := v.wallet.BondTransaction(pubKey, address, "Voucher claim for bond in validator", int64(voucher.Amount))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if txHash == "" {
		return cmd.ErrorResult(errors.New("can't send bond transaction"))
	}

	if err = v.db.UpdateVoucherTx(voucher.ID, txHash); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher claimed successfully: https://pacviewer.com/transaction/%s", txHash)
}
