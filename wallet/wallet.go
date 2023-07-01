package wallet

import (
	"log"
	"os"
	"pactus-faucet/config"

	"github.com/k0kubun/pp"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	w "github.com/pactus-project/pactus/wallet"
)

const entropy = 128
const faucetAddressLabel = "faucet"

type Balance struct {
	Available int64
	Staked    int64
}

func Create(cfg *config.Config, mnemonic string) *w.Wallet {
	network := genesis.Testnet
	if mnemonic == "" {
		mnemonic = w.GenerateMnemonic(entropy)
	}
	mywallet, err := w.Create(cfg.WalletPath, mnemonic, cfg.Password, network)

	if err != nil {
		log.Printf("error creating wallet: %v", err)
		return nil
	}
	address, err := mywallet.DeriveNewAddress(faucetAddressLabel)
	if err != nil {
		log.Printf("error deriving wallet faucet address: %v", err)
		return nil
	}
	cfg.FaucetAddress = address
	err = cfg.Save()
	if err != nil {
		log.Printf("error updating configuration faucet address: %v", err)
		return nil
	}
	err = mywallet.Save()
	if err != nil {
		log.Printf("error saving wallet: %v", err)
		return nil
	}
	pp.Printf("Wallet created successfully at: %s\n", mywallet.Path())
	pp.Printf("Seed: \"%v\"\n", mnemonic)
	pp.Printf("Please keep your seed in a safe place;\nif you lose it, you will not be able to restore your wallet.\n")
	return mywallet
}

func Open(cfg *config.Config) *w.Wallet {
	if doesWalletExist(cfg.WalletPath) {
		wt, err := w.Open(cfg.WalletPath, true)
		if err != nil {
			log.Printf("error opening exising wallet: %v", err)
			return nil
		}
		err = wt.Connect(cfg.Server)
		if err != nil {
			log.Printf("error establishing connection: %v", err)
			return nil
		}
		return wt
	}
	//if the wallet does not exist, create one
	return Create(cfg, "")
}

func Transfer(cfg *config.Config, wt *w.Wallet, toAddress string, amount float64) string {

	opts := []w.TxOption{
		w.OptionStamp(""),
		w.OptionFee(util.CoinToChange(0)),
		w.OptionSequence(int32(0)),
		w.OptionMemo(""),
	}
	tx, err := wt.MakeTransferTx(cfg.FaucetAddress, toAddress,
		util.CoinToChange(amount), opts...)
	if err != nil {
		log.Printf("error creating transfer transaction: %v", err)
		return ""
	}
	//sign transaction
	err = wt.SignTransaction(cfg.Password, tx)
	if err != nil {
		log.Printf("error signing transfer transaction: %v", err)
		return ""
	}

	//broadcast transaction
	res, err := wt.BroadcastTransaction(tx)
	if err != nil {
		log.Printf("error broadcasting transfer transaction: %v", err)
		return ""
	}

	err = wt.Save()
	if err != nil {
		log.Printf("error saving wallet transaction history: %v", err)
	}
	return res //return transaction hash
}
func BondTransaction(cfg *config.Config, wt *w.Wallet, toAddress string, amount float64) string {
	info := wt.AddressInfo(toAddress)

	pubKey := ""
	if info != nil {
		pubKey = info.Pub.String()
	}

	opts := []w.TxOption{
		w.OptionStamp(""),
		w.OptionFee(util.CoinToChange(0)),
		w.OptionSequence(int32(0)),
		w.OptionMemo(""),
	}
	tx, err := wt.MakeBondTx(cfg.FaucetAddress, toAddress, pubKey,
		util.CoinToChange(amount), opts...)
	if err != nil {
		log.Printf("error creating bond transaction: %v", err)
		return ""
	}
	//sign transaction
	err = wt.SignTransaction(cfg.Password, tx)
	if err != nil {
		log.Printf("error signing bond transaction: %v", err)
		return ""
	}

	//broadcast transaction
	res, err := wt.BroadcastTransaction(tx)
	if err != nil {
		log.Printf("error broadcasting bond transaction: %v", err)
		return ""
	}

	err = wt.Save()
	if err != nil {
		log.Printf("error saving wallet transaction history: %v", err)
	}
	return res //return transaction hash
}

func GetBalance(wt *w.Wallet, address string) *Balance {
	balance := &Balance{Available: 0, Staked: 0}
	b, err := wt.Balance(address)
	if err != nil {
		log.Printf("error getting balance: %v", err)
		return nil
	}
	balance.Available = b
	stake, err := wt.Stake(address)
	if err != nil {
		log.Printf("error getting staking amount: %v", err)
		return nil
	}
	balance.Staked = stake
	return balance
}

// function to check if file exists
func doesWalletExist(fileName string) bool {
	_, error := os.Stat(fileName)
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}
