package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DiscordToken          string   `json:"discord_token"`
	WalletPath            string   `json:"wallet_path"`
	WalletPassword        string   `json:"wallet_password"`
	Servers               []string `json:"servers"`
	FaucetAddress         string   `json:"faucet_address"`
	FaucetAmount          float64  `json:"faucet_amount"`
	ReferralerStakeAmount float64  `json:"referraler_stake_amount"` // who get faucet
	ValidatorDataPath     string   `json:"validator_data_path"`
	ReferralDataPath      string   `json:"referral_data_path"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error loading configuration file: %v", err)
		return nil, fmt.Errorf("error loading configuration file: %w", err)
	}

	cfg := &Config{}
	err = json.Unmarshal(file, cfg)

	if err != nil {
		log.Printf("error unmarshalling configuration file: %v", err)
		return nil, fmt.Errorf("error unmarshalling configuration file: %w", err)
	}
	return cfg, nil
}
