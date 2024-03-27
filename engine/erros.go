package engine

type WalletError struct {
	Reason string
}

func (e WalletError) Error() string {
	return e.Reason
}
