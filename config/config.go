package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configPath = "./data/config.json"

type Config struct {
	DiscordToken  string `json:"discord_token"`
	BotPrefix     string `json:"bot_prefix"`
	WalletPath    string `json:"wallet_path"`
	FaucetAddress string `json:"faucet_address"`
	Password      string `json:"password"`
	Server        string `json:"server"`
}

func Load() (*Config, error) {
	file, err := os.ReadFile(configPath)
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
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		log.Printf("failed to write to %s: %v", configPath, err)
		return fmt.Errorf("failed to write to %s: %v", configPath, err)
	}
	return nil
}
