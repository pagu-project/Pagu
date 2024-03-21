package wallet

import (
	"os"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	pwallet "github.com/pactus-project/pactus/wallet"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/log"
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

func Open(cfg *config.WalletConfig) IWallet {
	if doesWalletExist(cfg.Path) {

		wt, err := pwallet.Open(cfg.Path, true)
		if err != nil {
			log.Fatal("error opening existing wallet", "err", err)
		}

		err = wt.Connect(cfg.RPCUrl)
		if err != nil {
			log.Fatal("error establishing connection", "err", err)
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
	opts := []pwallet.TxOption{
		pwallet.OptionMemo(memo),
	}
	tx, err := w.wallet.MakeBondTx(w.address, toAddress, pubKey,
		amount, opts...)
	if err != nil {
		log.Error("error creating bond transaction", "err", err, "to",
			toAddress, "amount", util.ChangeToCoin(amount))
		return "", err
	}
	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing bond transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting bond transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
	}
	return res, nil // return transaction hash
}

func (w *Wallet) TransferTransaction(pubKey, toAddress, memo string, amount int64) (string, error) {
	fee, err := w.wallet.CalculateFee(int64(amount), payload.TypeTransfer)
	if err != nil {
		return "", err
	}

	opts := []pwallet.TxOption{
		pwallet.OptionFee(fee),
		pwallet.OptionMemo(memo),
	}

	tx, err := w.wallet.MakeTransferTx(w.address, toAddress, int64(amount), opts...)
	if err != nil {
		log.Error("error creating transfer transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
		return "", err
	}

	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing transfer transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting transfer transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
	}
	return res, nil // return transaction hash
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) Balance() int64 {
	balance, _ := w.wallet.Balance(w.address)
	return balance
}

func (w *Wallet) NewAddress(lb string) (string, error) {
	return w.wallet.NewBLSAccountAddress(lb)
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
