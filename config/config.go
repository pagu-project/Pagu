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
	Network      string
	NetworkNodes []string
	LocalNode    string
	DataBasePath string
	AuthIDs      []string
	DiscordBot   DiscordBot
	GRPC         GRPC
	PTWallet     PhoenixTestNetWallet
	Logger       Logger
	HTTP         HTTP
	Phoenix      PhoenixNetwork
	Telegram     Telegram
}

type PhoenixTestNetWallet struct {
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
	BotToken string
	ChatID   int64
	TgLink   string
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

	targets := strings.Split(os.Getenv("LOG_TARGETS"), ",")

	chatID, err := strconv.ParseInt(os.Getenv("ChatID"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ChatID: %w", err)
	}

	// Fetch config values from environment variables.
	cfg := &Config{
		Network: os.Getenv("NETWORK"),
		PTWallet: PhoenixTestNetWallet{
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
			Targets:    targets,
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
			BotToken: os.Getenv("BotToken"),
			ChatID:   chatID,
			TgLink:   os.Getenv("tg_link"),
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
	if cfg.PTWallet.Enable {
		if cfg.PTWallet.Address == "" {
			return fmt.Errorf("WALLET_ADDRESS is not set")
		}

		// Check if the WalletPath exists.
		if !util.PathExists(cfg.PTWallet.Address) {
			return fmt.Errorf("WALLET_PATH does not exist")
		}
	}

	if len(cfg.NetworkNodes) == 0 || len(cfg.Phoenix.NetworkNodes) == 0 {
		return fmt.Errorf("NETWORK_NODES is not set or incorrect")
	}

	return nil
}
