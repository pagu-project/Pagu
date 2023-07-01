package main

import (
	"log"
	"pactus-faucet/config"
	"pactus-faucet/wallet"

	"github.com/yudai/pp"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Printf("error loading configuration %v\n", err)
	}

	//load or create wallette
	w := wallet.Open(cfg)

	if w != nil {
		log.Println("wallet opened/created successfully")
	}

	validatorAddress := "tpc1prmusy99q9fhkaj6p9rfr9423qcg6epj9zdgglq"

	hash := wallet.Transfer(cfg, w, validatorAddress, 4.0)

	pp.Printf("bonding transaction hash: %v", hash)

	// path := "./store/wallet.json"
	// mnemonic := "mother hat quick unhappy swear nuclear reward glove regular napkin salad weapon"
	// password := "123456"

	// w := wallet.CreateWallet(path, mnemonic, password)

	// address, _ := w.DeriveNewAddress("faucet")
	// pp.Printf("address: %v", address)

	// w.Save()

}
