package voucher

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestStatusNormal(t *testing.T) {
	voucher, db, _, _ := setup(t)

	t.Run("one code status normal", func(t *testing.T) {
		now := time.Now()
		validMonths := uint8(2)
		voucherAmount, _ := amount.NewAmount(100)

		db.EXPECT().GetVoucherByCode("12345678").Return(
			entity.Voucher{
				ID:          1,
				Code:        "12345678",
				Desc:        "some_desc",
				Recipient:   "some_recipient",
				ValidMonths: validMonths,
				Amount:      voucherAmount,
				TxHash:      "some_transaction_hash",
				ClaimedBy:   0,
				Model: gorm.Model{
					CreatedAt: now,
				},
			}, nil,
		).AnyTimes()

		expTime := now.AddDate(0, int(validMonths), 0).Format("02/01/2006, 15:04:05")

		cmd := &command.Command{}
		caller := &entity.User{ID: 1}

		args := make(map[string]string)
		args["code"] = "12345678"
		result := voucher.statusHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Equal(t, result.Message, fmt.Sprintf("Code: 12345678\nAmount: 100 PAC\n"+
			"Expire At: %s\nRecipient: some_recipient\nDescription: some_desc\nClaimed: YES\n"+
			"Tx Link: https://pacviewer.com/transaction/some_transaction_hash"+
			"\n", expTime))
	})

	t.Run("wrong code", func(t *testing.T) {
		db.EXPECT().GetVoucherByCode("000").Return(
			entity.Voucher{}, errors.New(""),
		).AnyTimes()

		cmd := &command.Command{}
		caller := &entity.User{ID: 1}

		args := make(map[string]string)
		args["code"] = "000"
		result := voucher.statusHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "An error occurred: voucher code is not valid, no voucher found")
	})

	t.Run("list vouchers status normal", func(t *testing.T) {
		now := time.Now()
		validMonths := uint8(2)
		voucherAmount, _ := amount.NewAmount(100)

		db.EXPECT().ListVoucher().Return(
			[]*entity.Voucher{
				{
					ID:          1,
					Code:        "code1",
					ValidMonths: validMonths,
					Amount:      voucherAmount,
					TxHash:      "some_transaction_hash",
					Model: gorm.Model{
						CreatedAt: now,
					},
				},
				{
					ID:          2,
					Code:        "code2",
					ValidMonths: validMonths,
					Amount:      voucherAmount,
					Model: gorm.Model{
						CreatedAt: now,
					},
				},
				{
					ID:          3,
					Code:        "code3",
					ValidMonths: validMonths,
					Amount:      voucherAmount,
					Model: gorm.Model{
						CreatedAt: now.AddDate(0, -3, 0),
					},
				},
			}, nil,
		).AnyTimes()

		cmd := &command.Command{}
		caller := &entity.User{ID: 1}

		args := make(map[string]string)
		result := voucher.statusHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Equal(t, result.Message, "Total Codes: 3\nTotal Amount: 300 PAC\n\n\n"+
			"Claimed: 1\nTotal Claimed Amount: 100 PAC\nTotal Expired: 1"+
			"\n")
	})
}
