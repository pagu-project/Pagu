package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	WalletAddress  string
	WalletPath     string
	WalletPassword string
	RPCNodes       []string
	StorePath      string
	DiscordBotCfg  DiscordBotConfig
}

type DiscordBotConfig struct {
	DiscordToken   string
	DiscordGuildID string
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Fetch config values from environment variables.
	cfg := &Config{
		WalletAddress:  os.Getenv("WALLET_ADDRESS"),                // create a .env file and make a variable named WALLET_ADDRESS,put your wallet address there.
		WalletPath:     os.Getenv("WALLET_PATH"),                   // in .env file , make WALLET_PATH varaible and put your wallet path there.
		WalletPassword: os.Getenv("WALLET_PASSWORD"),               // in .env file, create WALLET PASSWORD variable and put your wallet password.
		RPCNodes:       strings.Split(os.Getenv("RPC_NODES"), ","), // in .env file, make RPC_NODES variable and put your RPC_NODES key there.
		StorePath:      os.Getenv("STORE_PATH"),                    // in .env file, make STORE_PATH variable and put your store_path there.
		DiscordBotCfg: DiscordBotConfig{
			DiscordToken:   os.Getenv("DISCORD_TOKEN"),    // in .env file, make Discord_TOKEN variable and put your Discord_TOKEN.
			DiscordGuildID: os.Getenv("DISCORD_GUILD_ID"), // in .env file, make Discord_Guild_ID variable and put your Discord_GUILD_ID.
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

	if len(cfg.RPCNodes) == 0 {
		return fmt.Errorf("RPCNODES is not set or incorrect")
	}

	if cfg.StorePath == "" {
		return fmt.Errorf("STORE_PATH is not set or incorrect")
	}

	if cfg.DiscordBotCfg.DiscordToken == "" {
		return fmt.Errorf("DISCORD_TOKEN is not set or incorrect")
	}

	// Check if the DiscordToken starts with 'MTE' which is discord's token prefix.
	if !strings.HasPrefix(cfg.DiscordBotCfg.DiscordToken, "MTE") {
		return fmt.Errorf("DISCORD_TOKEN does not start with the correct prefix or invalid")
	}

	if cfg.DiscordBotCfg.DiscordGuildID == "" {
		return fmt.Errorf("DISCORD_GUILD_ID is not set or incorrect")
	}

	return nil
}
