package wallet

type IWallet interface {
	BondTransaction(string, string, string, float64) (string, error)
	TransferTransaction(string, string, string, float64) (string, error)
	Address() string
}
