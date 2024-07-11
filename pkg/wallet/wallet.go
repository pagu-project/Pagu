package wallet

import (
	"errors"
	"os"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx/payload"
	pwallet "github.com/pactus-project/pactus/wallet"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/pkg/amount"
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

func Open(cfg *config.Wallet) (IWallet, error) {
	if doesWalletExist(cfg.Path) {
		wt, err := pwallet.Open(cfg.Path, false)
		if err != nil {
			return &Wallet{}, err
		}

		return &Wallet{
			wallet:   wt,
			address:  cfg.Address,
			password: cfg.Password,
		}, nil
	}

	// if the wallet does not exist, create one
	return &Wallet{}, errors.New("can't open the wallet")
}

func (w *Wallet) BondTransaction(pubKey, toAddress, memo string, amt amount.Amount) (string, error) {
	opts := []pwallet.TxOption{
		pwallet.OptionMemo(memo),
	}
	tx, err := w.wallet.MakeBondTx(w.address, toAddress, pubKey, amt.ToPactusAmount(), opts...)
	if err != nil {
		log.Error("error creating bond transaction", "err", err, "to",
			toAddress, "amount", amt)
		return "", err
	}
	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing bond transaction", "err", err,
			"to", toAddress, "amount", amt)
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting bond transaction", "err", err,
			"to", toAddress, "amount", amt)
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", amt)
	}
	return res, nil // return transaction hash
}

func (w *Wallet) TransferTransaction(toAddress, memo string, amt amount.Amount) (string, error) {
	// calculate fee using amount struct.
	fee, err := w.wallet.CalculateFee(amt.ToPactusAmount(), payload.TypeTransfer)
	if err != nil {
		log.Error("error calculating fee", "err", err, "client")
		return "", err
	}

	opts := []pwallet.TxOption{
		pwallet.OptionFee(fee),
		pwallet.OptionMemo(memo),
	}

	// Use amt.Amount for transaction amount.
	tx, err := w.wallet.MakeTransferTx(w.address, toAddress, amt.ToPactusAmount(), opts...)
	if err != nil {
		log.Error("error creating transfer transaction", "err", err,
			"from", w.address, "to", toAddress, "amount", amt)
		return "", err
	}

	// sign transaction.
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing transfer transaction", "err", err,
			"to", toAddress, "amount", amt)
		return "", err
	}

	// broadcast transaction.
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting transfer transaction", "err", err,
			"to", toAddress, "amount", amt)
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", amt)
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
