package wallet

import (
	"os"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx/payload"
	pwallet "github.com/pactus-project/pactus/wallet"
)

type Balance struct {
	Available float64
	Staked    float64
}

type Wallet struct {
	address  string
	password string
	wallet   *pwallet.Wallet
	logger   *log.SubLogger
}

func Open(cfg *config.Config, logger *log.SubLogger) IWallet {
	if doesWalletExist(cfg.WalletPath) {

		wt, err := pwallet.Open(cfg.WalletPath, true)
		if err != nil {
			logger.Fatal("error opening existing wallet", "err", err)
		}

		err = wt.Connect(cfg.LocalNode)
		if err != nil {
			logger.Fatal("error establishing connection", "err", err)
		}

		return &Wallet{
			wallet:   wt,
			address:  cfg.WalletAddress,
			password: cfg.WalletPassword,
			logger:   logger,
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
		w.logger.Error("error creating bond transaction", "err", err, "to",
			toAddress, "amount", utils.AtomicToCoin(amount))
		return "", err
	}
	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		w.logger.Error("error signing bond transaction", "err", err,
			"to", toAddress, "amount", utils.AtomicToCoin(amount))
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		w.logger.Error("error broadcasting bond transaction", "err", err,
			"to", toAddress, "amount", utils.AtomicToCoin(amount))
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		w.logger.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", utils.AtomicToCoin(amount))
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
		w.logger.Error("error creating transfer transaction", "err", err,
			"to", toAddress, "amount", utils.AtomicToCoin(amount))
		return "", err
	}

	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		w.logger.Error("error signing transfer transaction", "err", err,
			"to", toAddress, "amount", utils.AtomicToCoin(amount))
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		w.logger.Error("error broadcasting transfer transaction", "err", err,
			"to", toAddress, "amount", utils.AtomicToCoin(amount))
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		w.logger.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", utils.AtomicToCoin(amount))
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
