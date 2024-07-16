package voucher

import (
	"errors"
	"testing"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/wallet"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (*Voucher, repository.MockDatabase, client.MockManager, wallet.MockIWallet) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockDB := repository.NewMockDatabase(ctrl)
	mockClient := client.NewMockManager(ctrl)
	mockWallet := wallet.NewMockIWallet(ctrl)

	mockVoucher := NewVoucher(mockDB, mockWallet, mockClient)

	return mockVoucher, *mockDB, *mockClient, *mockWallet
}

func TestCreate(t *testing.T) {
	voucher, db, _, _ := setup(t)

	t.Run("normal", func(t *testing.T) {
		db.EXPECT().GetVoucherByCode(gomock.Any()).Return(
			entity.Voucher{}, errors.New(""),
		).AnyTimes()

		db.EXPECT().AddVoucher(gomock.Any()).Return(nil).AnyTimes()

		cmd := &command.Command{
			User: &entity.User{
				ID: 1,
			},
		}

		args := make(map[string]string)
		args["amount"] = "100"
		args["valid-months"] = "1"
		result := voucher.createHandler(cmd, entity.AppIDDiscord, "", args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})

	t.Run("more than 1000 PAC", func(t *testing.T) {
		db.EXPECT().GetVoucherByCode(gomock.Any()).Return(
			entity.Voucher{}, errors.New(""),
		).AnyTimes()

		db.EXPECT().AddVoucher(gomock.Any()).Return(nil).AnyTimes()

		cmd := &command.Command{
			User: &entity.User{
				ID: 1,
			},
		}

		args := make(map[string]string)
		args["amount"] = "1001"
		args["valid-months"] = "1"
		result := voucher.createHandler(cmd, entity.AppIDDiscord, "", args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "stake amount is more than 1000")
	})

	t.Run("wrong month", func(t *testing.T) {
		db.EXPECT().GetVoucherByCode(gomock.Any()).Return(
			entity.Voucher{}, errors.New(""),
		).AnyTimes()

		db.EXPECT().AddVoucher(gomock.Any()).Return(nil).AnyTimes()

		cmd := &command.Command{
			User: &entity.User{
				ID: 1,
			},
		}

		args := make(map[string]string)
		args["amount"] = "100"
		args["valid-months"] = "1.1"
		result := voucher.createHandler(cmd, entity.AppIDDiscord, "", args)
		assert.False(t, result.Successful)
	})

	t.Run("normal with optional arguments", func(t *testing.T) {
		db.EXPECT().GetVoucherByCode(gomock.Any()).Return(
			entity.Voucher{}, errors.New(""),
		).AnyTimes()

		db.EXPECT().AddVoucher(gomock.Any()).Return(nil).AnyTimes()

		cmd := &command.Command{
			User: &entity.User{
				ID: 1,
			},
		}

		args := make(map[string]string)
		args["amount"] = "100"
		args["valid-months"] = "12"
		args["recipient"] = "Kayhan"
		args["description"] = "Testnet node"
		result := voucher.createHandler(cmd, entity.AppIDDiscord, "", args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})
}
