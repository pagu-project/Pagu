package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	WalletAddress  string           `json:"wallet_address"`
	WalletPath     string           `json:"wallet_path"`
	WalletPassword string           `json:"wallet_password"`
	RPCNodes       []string         `json:"rpc_nodes"`
	StorePath      string           `json:"store_path"`
	DiscordBotCfg  DiscordBotConfig `json:"discord_bot_config"`
}

type DiscordBotConfig struct {
	DiscordToken   string `json:"discord_token"`
	DiscordGuildID string `json:"discord_guild_id"`
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
