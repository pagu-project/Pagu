package repository

import "github.com/pagu-project/Pagu/internal/entity"

type Database interface {
	AddFaucet(f *entity.PhoenixFaucet) error
	CanGetFaucet(user *entity.User) bool
	AddUser(u *entity.User) error
	HasUser(id string) bool
	GetUserInApp(appID entity.AppID, callerID string) (*entity.User, error)
	AddVoucher(v *entity.Voucher) error
	GetVoucherByCode(code string) (entity.Voucher, error)
	ClaimVoucher(id uint, txHash string, claimer uint) error
	GetZealyUser(id string) (*entity.ZealyUser, error)
	AddZealyUser(u *entity.ZealyUser) error
	UpdateZealyUser(id string, txHash string) error
	GetAllZealyUser() ([]*entity.ZealyUser, error)
}
