package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	Network        string
	WalletAddress  string
	WalletPath     string
	WalletPassword string
	NetworkNodes   []string
	LocalNode      string
	StorePath      string
	DiscordBotCfg  DiscordBotConfig
}

type DiscordBotConfig struct {
	DiscordToken   string
	DiscordGuildID string
}

func Load(filePaths ...string) (*Config, error) {
	err := godotenv.Load(filePaths...)
	if err != nil {
		return nil, err
	}

	// Fetch config values from environment variables.
	cfg := &Config{
		Network:        os.Getenv("NETWORK"),
		WalletAddress:  os.Getenv("WALLET_ADDRESS"),
		WalletPath:     os.Getenv("WALLET_PATH"),
		WalletPassword: os.Getenv("WALLET_PASSWORD"),
		LocalNode:      os.Getenv("LOCAL_NODE"),
		NetworkNodes:   strings.Split(os.Getenv("NETWORK_NODES"), ","),
		StorePath:      os.Getenv("STORE_PATH"),
		DiscordBotCfg: DiscordBotConfig{
			DiscordToken:   os.Getenv("DISCORD_TOKEN"),
			DiscordGuildID: os.Getenv("DISCORD_GUILD_ID"),
		},
	}

	// Check if the required configurations are set.
	if err := cfg.BasicCheck(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks for the presence of required environment variables.
func (cfg *Config) BasicCheck() error {
	if cfg.WalletAddress == "" {
		return fmt.Errorf("WALLET_ADDRESS is not set")
	}

	// Check if the WalletPath exists.
	if !util.PathExists(cfg.WalletPath) {
		return fmt.Errorf("WALLET_PATH does not exist")
	}

	if len(cfg.NetworkNodes) == 0 {
		return fmt.Errorf("RPCNODES is not set or incorrect")
	}

	if cfg.StorePath == "" {
		return fmt.Errorf("STORE_PATH is not set or incorrect")
	}

	// if cfg.DiscordBotCfg.DiscordToken == "" {
	// 	return fmt.Errorf("DISCORD_TOKEN is not set or incorrect")
	// }

	// // Check if the DiscordToken starts with 'MTE' which is discord's token prefix.
	// if !strings.HasPrefix(cfg.DiscordBotCfg.DiscordToken, "MTE") {
	// 	return fmt.Errorf("DISCORD_TOKEN does not start with the correct prefix or invalid")
	// }

	// if cfg.DiscordBotCfg.DiscordGuildID == "" {
	// 	return fmt.Errorf("DISCORD_GUILD_ID is not set or incorrect")
	// }

	return nil
}
