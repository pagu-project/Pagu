package voucher

import (
	"errors"
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/stretchr/testify/assert"
)

func TestClaimNormal(t *testing.T) {
	voucher, db, client, wallet := setup(t)

	t.Run("normal", func(t *testing.T) {
		amt, _ := amount.NewAmount(100)
		db.EXPECT().GetVoucherByCode("12345678").Return(
			entity.Voucher{
				ValidMonths: 1,
				Amount:      amt,
				ID:          1,
			}, nil,
		).AnyTimes()

		client.EXPECT().GetValidatorInfo("pc1z").Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					PublicKey: "pc1z",
				},
			}, nil,
		).AnyTimes()

		wallet.EXPECT().BondTransaction("pc1z", "pc1z", "Voucher claim from Pagu", amt).Return(
			"0x1", nil,
		).AnyTimes()

		db.EXPECT().ClaimVoucher(uint(1), "0x1", uint(1)).Return(
			nil,
		).AnyTimes()

		cmd := &command.Command{
			User: &entity.User{
				ID: 1,
			},
		}

		args := make(map[string]string)
		args["code"] = "12345678"
		args["address"] = "pc1z"
		result := voucher.claimHandler(cmd, entity.AppIDDiscord, "", args)
		assert.True(t, result.Successful)
		assert.Equal(t, result.Message, "Voucher claimed successfully!\n\n https://pacviewer.com/transaction/0x1")
	})

	t.Run("wrong code", func(t *testing.T) {
		cmd := &command.Command{
			User: &entity.User{
				ID: 1,
			},
		}

		args := make(map[string]string)
		args["code"] = "0"
		args["address"] = "pc1z"
		result := voucher.claimHandler(cmd, entity.AppIDDiscord, "", args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "An error occurred: voucher code is not valid, length must be 8")
	})
}

func TestClaimNotFound(t *testing.T) {
	voucher, db, _, _ := setup(t)

	db.EXPECT().GetVoucherByCode("12345678").Return(
		entity.Voucher{}, errors.New(""),
	).AnyTimes()

	cmd := &command.Command{
		User: &entity.User{
			ID: 1,
		},
	}

	args := make(map[string]string)
	args["code"] = "12345678"
	args["address"] = "pc1z"
	result := voucher.claimHandler(cmd, entity.AppIDDiscord, "", args)
	assert.False(t, result.Successful)
	assert.Equal(t, result.Message, "An error occurred: voucher code is not valid, no voucher found")
}

func TestClaimAlreadyClaimed(t *testing.T) {
	voucher, db, _, _ := setup(t)

	db.EXPECT().GetVoucherByCode("12345678").Return(
		entity.Voucher{
			TxHash: "123456789",
		}, nil,
	).AnyTimes()

	cmd := &command.Command{
		User: &entity.User{
			ID: 1,
		},
	}

	args := make(map[string]string)
	args["code"] = "12345678"
	args["address"] = "pc1z"
	result := voucher.claimHandler(cmd, entity.AppIDDiscord, "", args)
	assert.False(t, result.Successful)
	assert.Equal(t, result.Message, "An error occurred: voucher code claimed before")
}
