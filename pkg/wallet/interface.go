package wallet

import "github.com/pagu-project/Pagu/pkg/amount"

type IWallet interface {
	NewAddress(lb string) (string, error)
	Balance() int64
	Address() string
	TransferTransaction(toAddress, memo string, amt amount.Amount) (string, error)
	BondTransaction(pubKey, toAddress, memo string, amt amount.Amount) (string, error)
}
