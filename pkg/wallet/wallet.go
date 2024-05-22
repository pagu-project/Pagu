package wallet

import (
	"os"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	amt "github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	pwallet "github.com/pactus-project/pactus/wallet"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/pkg/log"
)

type Balance struct {
	Available float64
	Staked    float64
}

type Wallet struct {
	address  string
	password string
	wallet   *pwallet.Wallet
}

func Open(cfg *config.Wallet) *Wallet {
	if doesWalletExist(cfg.Path) {

		wt, err := pwallet.Open(cfg.Path, false)
		if err != nil {
			log.Fatal("error opening existing wallet", "err", err)
		}

		return &Wallet{
			wallet:   wt,
			address:  cfg.Address,
			password: cfg.Password,
		}
	}

	// if the wallet does not exist, create one
	return nil
}

func (w *Wallet) BondTransaction(pubKey, toAddress, memo string, amount int64) (string, error) {
	amountInNanoPAC := amt.Amount(amount)

	opts := []pwallet.TxOption{
		pwallet.OptionMemo(memo),
	}
	tx, err := w.wallet.MakeBondTx(w.address, toAddress, pubKey,
		amountInNanoPAC, opts...)
	if err != nil {
		log.Error("error creating bond transaction", "err", err, "to",
			toAddress, "amount", amountInNanoPAC.Format(amt.UnitNanoPAC))
		return "", err
	}
	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing bond transaction", "err", err,
			"to", toAddress, "amount", amountInNanoPAC.String())
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting bond transaction", "err", err,
			"to", toAddress, "amount", amountInNanoPAC.Format(amt.UnitNanoPAC))
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", amountInNanoPAC.Format(amt.UnitNanoPAC))
	}
	return res, nil // return transaction hash
}

func (w *Wallet) TransferTransaction(toAddress, memo string, amount int64) (string, error) {
	// Convert int64 to amt.Amount.
	amountInNanoPAC, err := amt.NewAmount(float64(amount))
	if err != nil {
		return "", err
	}

	// claculate fee using amount struct.
	fee, err := w.wallet.CalculateFee(amountInNanoPAC, payload.TypeTransfer)
	if err != nil {
		return "", err
	}

	opts := []pwallet.TxOption{
		pwallet.OptionFee(fee),
		pwallet.OptionMemo(memo),
	}

	// Use amt.Amount for transaction amount.
	tx, err := w.wallet.MakeTransferTx(w.address, toAddress, amountInNanoPAC, opts...)
	if err != nil {
		log.Error("error creating transfer transaction", "err", err,
			"to", toAddress, "amount", amountInNanoPAC.Format(amt.UnitNanoPAC))
		return "", err
	}

	// sign transaction.
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing transfer transaction", "err", err,
			"to", toAddress, "amount", amountInNanoPAC.Format(amt.UnitNanoPAC))
		return "", err
	}

	// broadcast transaction.
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting transfer transaction", "err", err,
			"to", toAddress, "amount", amountInNanoPAC.Format(amt.UnitNanoPAC))
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", amountInNanoPAC.Format(amt.UnitNanoPAC))
	}
	return res, nil // return transaction hash.
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) Balance() int64 {
	balance, _ := w.wallet.Balance(w.address)
	return int64(balance)
}

func (w *Wallet) NewAddress(lb string) (string, error) {
	addressInfo, err := w.wallet.NewBLSAccountAddress(lb)
	if err != nil {
		return "", err
	}
	return addressInfo.Address, nil
}

func IsValidData(address, pubKey string) bool {
	addr, err := crypto.AddressFromString(address)
	if err != nil {
		return false
	}
	pub, err := bls.PublicKeyFromString(pubKey)
	if err != nil {
		return false
	}
	err = pub.VerifyAddress(addr)
	return err == nil
}

// function to check if file exists.
func doesWalletExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}
