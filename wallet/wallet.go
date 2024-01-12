package wallet

import (
	"os"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	pwallet "github.com/pactus-project/pactus/wallet"
)

type Balance struct {
	Available float64
	Staked    float64
}

type Wallet struct {
	address  string
	wallet   *pwallet.Wallet
	password string
}

func Open(cfg *config.Config) IWallet {
	if doesWalletExist(cfg.WalletPath) {

		wt, err := pwallet.Open(cfg.WalletPath, true)
		if err != nil {
			log.Fatal("error opening existing wallet", "err", err)
		}

		err = wt.Connect(cfg.RPCNodes[0])
		if err != nil {
			log.Fatal("error establishing connection", "err", err)
		}

		return &Wallet{
			wallet:   wt,
			address:  cfg.FaucetAddress,
			password: cfg.WalletPassword,
		}
	}

	// if the wallet does not exist, create one
	return nil
}

func (w *Wallet) BondTransaction(pubKey, toAddress, memo string, amount float64) (string, error) {
	opts := []pwallet.TxOption{
		pwallet.OptionFee(util.CoinToChange(0)),
		pwallet.OptionMemo(memo),
	}
	tx, err := w.wallet.MakeBondTx(w.address, toAddress, pubKey,
		util.CoinToChange(amount), opts...)
	if err != nil {
		log.Error("error creating bond transaction", "err", err, "address", toAddress, "amount", amount)
		return "", err
	}
	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing bond transaction", "err", err, "address", toAddress, "amount", amount)
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting bond transaction", "err", err, "address", toAddress, "amount", amount)
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err, "address", toAddress, "amount", amount)
	}
	return res, nil // return transaction hash
}

func (w *Wallet) TransferTransaction(pubKey, toAddress, memo string, amount float64) (string, error) {
	fee, err := w.wallet.CalculateFee(int64(amount), payload.TypeTransfer)
	if err != nil {
		return "", err
	}

	opts := []pwallet.TxOption{
		pwallet.OptionFee(util.CoinToChange(float64(fee))),
		pwallet.OptionMemo(memo),
	}

	tx, err := w.wallet.MakeTransferTx(w.address, toAddress, int64(amount), opts...)
	if err != nil {
		log.Error("error creating transfer transaction", "err", err, "address", toAddress, "amount", amount)
		return "", err
	}

	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing transfer transaction", "err", err, "address", toAddress, "amount", amount)
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting transfer transaction", "err", err, "address", toAddress, "amount", amount)
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err, "address", toAddress, "amount", amount)
	}
	return res, nil // return transaction hash
}

func (w *Wallet) Address() string {
	return w.address
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
