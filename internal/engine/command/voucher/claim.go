package voucher

import (
	"errors"
	"time"

	amt "github.com/pactus-project/pactus/types/amount"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/log"
)

func (v *Voucher) claimHandler(cmd command.Command, _ entity.AppID, callerID string, args ...string) command.CommandResult {
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

	if len(voucher.TxHash) > 0 {
		return cmd.ErrorResult(errors.New("voucher code claimed before"))
	}

	address := args[1]
	validatorInfo, err := v.clientManager.GetValidatorInfo(address)
	if err != nil {
		log.Error("error get validator info", "err", err)
		return cmd.ErrorResult(errors.New("bond error"))
	}

	pubKey := validatorInfo.GetValidator().GetPublicKey()

	amountInNanoPAC, err := amt.NewAmount(float64(voucher.Amount))
	if err != nil {
		log.Error("error converting amount to nanoPAC", "err", err)
		return cmd.ErrorResult(errors.New("bond error"))
	}

	txHash, err := v.wallet.BondTransaction(pubKey, address, "Voucher claim for bond in validator", amountInNanoPAC.ToNanoPAC())
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if txHash == "" {
		return cmd.ErrorResult(errors.New("can't send bond transaction"))
	}

	if err = v.db.ClaimVoucher(voucher.ID, txHash, cmd.User.ID); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher claimed successfully: https://pacviewer.com/transaction/%s", txHash)
}
