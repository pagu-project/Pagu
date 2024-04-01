package wallet

type IWallet interface {
	BondTransaction(string, string, string, int64) (string, error)
	TransferTransaction(string, string, int64) (string, error)
	NewAddress(string) (string, error)
	Address() string
	Balance() int64
}
