package config

import (
	"fmt"
	"os"

	"github.com/pactus-project/pactus/util"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Network      string         `yaml:"network"`
	NetworkNodes []string       `yaml:"network_nodes"`
	LocalNode    string         `yaml:"local_node"`
	Database     Database       `yaml:"database"`
	AuthIDs      []string       `yaml:"auth_ids"`
	GRPC         GRPC           `yaml:"grpc"`
	Wallet       Wallet         `yaml:"wallet"`
	Logger       Logger         `yaml:"logger"`
	HTTP         HTTP           `yaml:"http"`
	Phoenix      PhoenixNetwork `yaml:"phoenix"`
	DiscordBot   DiscordBot     `yaml:"discord"`
	Telegram     Telegram       `yaml:"telegram"`
	TargetMask   int            `yaml:"target_mask"`
}

type Database struct {
	URL string `yaml:"url"`
}

type Wallet struct {
	Enable   bool   `yaml:"enable"`
	Address  string `yaml:"address"`
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
	RPCUrl   string `yaml:"rpc"`
}

type DiscordBot struct {
	Token   string `yaml:"token"`
	GuildID string `yaml:"guild_id"`
}

type GRPC struct {
	Listen string `yaml:"listen"`
}

type HTTP struct {
	Listen string `yaml:"listen"`
}

type PhoenixNetwork struct {
	FaucetAmount uint `yaml:"faucet_amount"`
}

type Logger struct {
	Filename   string   `yaml:"filename"`
	LogLevel   string   `yaml:"level"`
	Targets    []string `yaml:"targets"`
	MaxSize    int      `yaml:"max_size"`
	MaxBackups int      `yaml:"max_backups"`
	Compress   bool     `yaml:"compress"`
}

type Telegram struct {
	BotToken  string `yaml:"bot_token"`
	ChatID    int64  `yaml:"chat_id"`
	GroupLink string `yaml:"group_link"`
}

func Load(path string) (*Config, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(payload, cfg); err != nil {
		return nil, err
	}

	// Check if the required configurations are set
	if err := cfg.BasicCheck(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// BasicCheck validate presence of required config variables.
func (cfg *Config) BasicCheck() error {
	if cfg.Wallet.Enable {
		if cfg.Wallet.Address == "" {
			return fmt.Errorf("config: basic check error: WALLET_ADDRESS dose not set")
		}

		// Check if the WalletPath exists.
		if !util.PathExists(cfg.Wallet.Path) {
			return fmt.Errorf("config: basic check error: WALLET_PATH does not exist: %s", cfg.Wallet.Path)
		}
	}

	if len(cfg.NetworkNodes) == 0 {
		return fmt.Errorf("config: basic check error: NETWORK_NODES is not set or incorrect")
	}

	return nil
}
