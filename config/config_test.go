package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBasicCheck tests the BasicCheck method of the Config struct.
func TestBasicCheck(t *testing.T) {
	// Create a temporary directory for the WalletPath
	tempWalletPath := t.TempDir()
	tempStorePath := t.TempDir()

	// Define test cases
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Valid config",
			cfg: Config{
				WalletAddress:  "test_wallet_address",
				WalletPath:     tempWalletPath, // Use the temporary directory
				WalletPassword: "test_password",
				NetworkNodes:   []string{"http://127.0.0.1:8545"},
				StorePath:      tempStorePath, // Use the temporary directory
				DiscordBotCfg: DiscordBotConfig{
					DiscordToken:   "MTEabc123",
					DiscordGuildID: "123456789",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid RPCNodes",
			cfg: Config{
				WalletAddress:  "test_wallet_address",
				WalletPath:     "/valid/path",
				WalletPassword: "test_password",
				NetworkNodes:   []string{},
				StorePath:      "/valid/storepath",
				DiscordBotCfg: DiscordBotConfig{
					DiscordToken:   "MTEabc123",
					DiscordGuildID: "123456789",
				},
			},
			wantErr: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform the check
			err := tt.cfg.BasicCheck()

			// Assert the error based on wantErr
			if tt.wantErr {
				assert.Error(t, err, "Config.BasicCheck() should return an error")
			} else {
				assert.NoError(t, err, "Config.BasicCheck() should not return an error")
			}
		})
	}
}
