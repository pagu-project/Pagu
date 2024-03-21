package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	Network       string
	NetworkNodes  []string
	LocalNode     string
	DataBasePath  string
	AuthIDs       []string
	DiscordBotCfg DiscordBotConfig
	GRPCConfig    GRPCConfig
	WalletConfig  WalletConfig
}

type WalletConfig struct {
	Enable   bool
	Address  string
	Path     string
	Password string
	RPCUrl   string
}

type DiscordBotConfig struct {
	Token   string
	GuildID string
}

type GRPCConfig struct {
	Listen string
}

func Load(filePaths ...string) (*Config, error) {
	err := godotenv.Load(filePaths...)
	if err != nil {
		return nil, err
	}

	enableWallet, err := strconv.ParseBool(os.Getenv("ENABLE_WALLET"))
	if err != nil {
		return nil, err
	}

	// Fetch config values from environment variables.
	cfg := &Config{
		Network: os.Getenv("NETWORK"),
		WalletConfig: WalletConfig{
			Address:  os.Getenv("WALLET_ADDRESS"),
			Path:     os.Getenv("WALLET_PATH"),
			Password: os.Getenv("WALLET_PASSWORD"),
			RPCUrl:   os.Getenv("WALLET_PRC"),
			Enable:   enableWallet,
		},
		LocalNode:    os.Getenv("LOCAL_NODE"),
		NetworkNodes: strings.Split(os.Getenv("NETWORK_NODES"), ","),
		DataBasePath: os.Getenv("DATABASE_PATH"),
		AuthIDs:      strings.Split(os.Getenv("AUTHORIZED_DISCORD_IDS"), ","),
		DiscordBotCfg: DiscordBotConfig{
			Token:   os.Getenv("DISCORD_TOKEN"),
			GuildID: os.Getenv("DISCORD_GUILD_ID"),
		},
		GRPCConfig: GRPCConfig{
			Listen: os.Getenv("GRPC_LISTEN"),
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
	if cfg.WalletConfig.Address == "" {
		return fmt.Errorf("WALLET_ADDRESS is not set")
	}

	// Check if the WalletPath exists.
	if !util.PathExists(cfg.WalletConfig.Path) {
		return fmt.Errorf("WALLET_PATH does not exist")
	}

	if len(cfg.NetworkNodes) == 0 {
		return fmt.Errorf("RPCNODES is not set or incorrect")
	}

	return nil
}
