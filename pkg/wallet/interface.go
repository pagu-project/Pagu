package wallet

type IWallet interface {
	NewAddress(lb string) (string, error)
	Balance() int64
	Address() string
	TransferTransaction(toAddress, memo string, amount int64) (string, error)
	BondTransaction(pubKey, toAddress, memo string, amount int64) (string, error)
}
