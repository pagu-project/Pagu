package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const configPath = "/bot-data/config.json"

type Config struct {
	DiscordToken      string  `json:"discord_token"`
	WalletPath        string  `json:"wallet_path"`
	WalletPassword    string  `json:"wallet_password"`
	Server            string  `json:"server"`
	FaucetAddress     string  `json:"faucet_address"`
	FaucetAmount      float64 `json:"faucet_amount"`
	ValidatorDataPath string  `json:"validator_data_path"`
}

func Load() (*Config, error) {
	file, err := os.ReadFile(filepath.Join(configPath))
	if err != nil {
		log.Printf("error loading configuration file: %v", err)
		return nil, fmt.Errorf("error loading configuration file: %v", err)
	}

	cfg := &Config{}
	err = json.Unmarshal(file, cfg)

	if err != nil {
		log.Printf("error unmarshalling configuration file: %v", err)
		return nil, fmt.Errorf("error unmarshalling configuration file: %v", err)
	}
	return cfg, nil
}

func (cfg *Config) Save() error {
	data, err := json.MarshalIndent(cfg, "  ", "  ")
	if err != nil {
		log.Printf("error marshalling configuration file: %v", err)
		return fmt.Errorf("error marshalling configuration file: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		log.Printf("failed to write to %s: %v", configPath, err)
		return fmt.Errorf("failed to write to %s: %v", configPath, err)
	}
	return nil
}
