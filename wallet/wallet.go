package wallet

import (
	"log"
	"os"

	"pactus-bot/config"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
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

func Open(cfg *config.Config) *Wallet {
	if doesWalletExist(cfg.WalletPath) {
		wt, err := pwallet.Open(cfg.WalletPath, true)
		if err != nil {
			log.Printf("error opening existing wallet: %v", err)
			return nil
		}
		err = wt.Connect(cfg.Servers[0])
		if err != nil {
			log.Printf("error establishing connection: %v", err)
			return nil
		}
		return &Wallet{wallet: wt, address: cfg.FaucetAddress, password: cfg.WalletPassword}
	}
	// if the wallet does not exist, create one
	return nil
}

func (w *Wallet) BondTransaction(pubKey, toAddress string, amount float64) string {
	opts := []pwallet.TxOption{
		pwallet.OptionFee(util.CoinToChange(0)),
		pwallet.OptionMemo("faucet from PactusBot"),
	}
	tx, err := w.wallet.MakeBondTx(w.address, toAddress, pubKey,
		util.CoinToChange(amount), opts...)
	if err != nil {
		log.Printf("error creating bond transaction: %v", err)
		return ""
	}
	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Printf("error signing bond transaction: %v", err)
		return ""
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Printf("error broadcasting bond transaction: %v", err)
		return ""
	}

	err = w.wallet.Save()
	if err != nil {
		log.Printf("error saving wallet transaction history: %v", err)
	}
	return res // return transaction hash
}

func (w *Wallet) GetBalance() *Balance {
	balance := &Balance{Available: 0, Staked: 0}
	b, err := w.wallet.Balance(w.address)
	if err != nil {
		log.Printf("error getting balance: %v", err)
		return balance
	}
	balance.Available = util.ChangeToCoin(b)
	stake, err := w.wallet.Stake(w.address)
	if err != nil {
		log.Printf("error getting staking amount: %v", err)
		return balance
	}
	balance.Staked = util.ChangeToCoin(stake)
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

// function to check if file exists
func doesWalletExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}
