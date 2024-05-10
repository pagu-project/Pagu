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
	DiscordBot    DiscordBot
	GRPC          GRPC
	Wallet        Wallet
	TestNetWallet Wallet
	Logger        Logger
	HTTP          HTTP
	Phoenix       PhoenixNetwork
	Telegram      Telegram
}

type Wallet struct {
	Enable   bool
	Address  string
	Path     string
	Password string
	RPCUrl   string
}

type DiscordBot struct {
	Token   string
	GuildID string
}

type GRPC struct {
	Listen string
}

type HTTP struct {
	Listen string
}

type PhoenixNetwork struct {
	NetworkNodes []string
	FaucetAmount uint
}

type Logger struct {
	Filename   string
	LogLevel   string
	Targets    []string
	MaxSize    int
	MaxBackups int
	Compress   bool
}

type Telegram struct {
	BotToken  string
	ChatID    int64
	GroupLink string
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

	enableTestNetWallet, err := strconv.ParseBool(os.Getenv("ENABLE_TESTNET_WALLET"))
	if err != nil {
		return nil, err
	}

	maxSizeStr := os.Getenv("LOG_MAX_SIZE")
	maxSize, err := strconv.Atoi(maxSizeStr)
	if err != nil {
		return nil, err
	}

	maxBackupsStr := os.Getenv("LOG_MAX_BACKUPS")
	maxBackups, err := strconv.Atoi(maxBackupsStr)
	if err != nil {
		return nil, err
	}

	compressStr := os.Getenv("LOG_COMPRESS")
	compress, err := strconv.ParseBool(compressStr)
	if err != nil {
		return nil, err
	}

	faucetAmountStr := os.Getenv("PHOENIX_FAUCET_AMOUNT")
	faucetAmount, err := strconv.ParseUint(faucetAmountStr, 10, 8)
	if err != nil {
		return nil, err
	}

	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		return nil, err
	}

	// Fetch config values from environment variables.
	cfg := &Config{
		Network: os.Getenv("NETWORK"),
		Wallet: Wallet{
			Address:  os.Getenv("WALLET_ADDRESS"),
			Path:     os.Getenv("WALLET_PATH"),
			Password: os.Getenv("WALLET_PASSWORD"),
			RPCUrl:   os.Getenv("WALLET_PRC"),
			Enable:   enableWallet,
		},
		TestNetWallet: Wallet{
			Address:  os.Getenv("TESTNET_WALLET_ADDRESS"),
			Path:     os.Getenv("TESTNET_WALLET_PATH"),
			Password: os.Getenv("TESTNET_WALLET_PASSWORD"),
			RPCUrl:   os.Getenv("TESTNET_WALLET_PRC"),
			Enable:   enableTestNetWallet,
		},
		LocalNode:    os.Getenv("LOCAL_NODE"),
		NetworkNodes: strings.Split(os.Getenv("NETWORK_NODES"), ","),
		DataBasePath: os.Getenv("DATABASE_PATH"),
		AuthIDs:      strings.Split(os.Getenv("AUTHORIZED_DISCORD_IDS"), ","),
		DiscordBot: DiscordBot{
			Token:   os.Getenv("DISCORD_TOKEN"),
			GuildID: os.Getenv("DISCORD_GUILD_ID"),
		},
		GRPC: GRPC{
			Listen: os.Getenv("GRPC_LISTEN"),
		},
		Logger: Logger{
			LogLevel:   os.Getenv("LOG_LEVEL"),
			Filename:   os.Getenv("LOG_FILENAME"),
			Targets:    strings.Split(os.Getenv("LOG_TARGETS"), ","),
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			Compress:   compress,
		},
		HTTP: HTTP{
			Listen: os.Getenv("HTTP_LISTEN"),
		},
		Phoenix: PhoenixNetwork{
			NetworkNodes: strings.Split(os.Getenv("PHOENIX_NETWORK_NODES"), ","),
			FaucetAmount: uint(faucetAmount),
		},
		Telegram: Telegram{
			BotToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
			ChatID:    chatID,
			GroupLink: os.Getenv("TELEGRAM_GROUP_LINK"),
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
	if cfg.Wallet.Enable {
		if cfg.Wallet.Address == "" {
			return fmt.Errorf("config: basic check error: WALLET_ADDRESS dose not set")
		}

		// Check if the WalletPath exists.
		if !util.PathExists(cfg.Wallet.Path) {
			return fmt.Errorf("config: basic check error: WALLET_PATH does not exist")
		}
	}

	if cfg.TestNetWallet.Enable {
		if cfg.TestNetWallet.Address == "" {
			return fmt.Errorf("config: basic check error: TESTNET_WALLET_ADDRESS dose not set")
		}

		// Check if the WalletPath exists.
		if !util.PathExists(cfg.TestNetWallet.Path) {
			return fmt.Errorf("config: basic check error: TESTNET_WALLET_PATH does not exist")
		}
	}

	if len(cfg.NetworkNodes) == 0 || len(cfg.Phoenix.NetworkNodes) == 0 {
		return fmt.Errorf("config: basic check error: NETWORK_NODES is not set or incorrect")
	}

	return nil
}
