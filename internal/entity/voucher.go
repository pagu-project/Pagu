package entity

import (
	"github.com/pagu-project/Pagu/pkg/amount"
	"gorm.io/gorm"
)

type Voucher struct {
	ID          uint          `gorm:"primaryKey;unique"`
	Creator     uint          `gorm:"size:255"`
	Code        string        `gorm:"size:8"`
	Amount      amount.Amount `gorm:"column:amount"`
	Desc        string
	Email       string
	Recipient   string
	ValidMonths uint8
	TxHash      string
	ClaimedBy   uint

	gorm.Model
}

func (Voucher) TableName() string {
	return "voucher"
}

func (v *Voucher) IsClaimed() bool {
	return v.TxHash != ""
}
