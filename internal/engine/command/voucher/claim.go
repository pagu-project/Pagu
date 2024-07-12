package voucher

import (
	"errors"
	"time"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/log"
)

func (v *Voucher) claimHandler(cmd *command.Command,
	_ entity.AppID, _ string, args ...string,
) command.CommandResult {
	code := args[0]
	if len(code) != 8 {
		return cmd.ErrorResult(errors.New("voucher code is not valid, length must be 8"))
	}

	voucher, err := v.db.GetVoucherByCode(code)
	if err != nil {
		return cmd.ErrorResult(errors.New("voucher code is not valid, no voucher found"))
	}

	now := time.Now().Month()
	if voucher.CreatedAt.Month() >= (now + time.Month(voucher.ValidMonths)) {
		return cmd.ErrorResult(errors.New("voucher is expired"))
	}

	if voucher.IsClaimed() {
		return cmd.ErrorResult(errors.New("voucher code claimed before"))
	}

	address := args[1]
	validatorInfo, err := v.clientManager.GetValidatorInfo(address)
	if err != nil {
		log.Error("error get validator info", "err", err)
		return cmd.ErrorResult(err)
	}

	pubKey := validatorInfo.GetValidator().GetPublicKey()
	txHash, err := v.wallet.BondTransaction(pubKey, address, "Voucher claim from Pagu", voucher.Amount)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if txHash == "" {
		return cmd.ErrorResult(errors.New("can't send bond transaction"))
	}

	if err = v.db.ClaimVoucher(voucher.ID, txHash, cmd.User.ID); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher claimed successfully!\n\n https://pacviewer.com/transaction/%s", txHash)
}
