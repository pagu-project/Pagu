package config

import (
	"testing"
)

// TestBasicCheck tests the BasicCheck method of the Config struct.
func TestBasicCheck(t *testing.T) {
	// Setup - Create a temporary directory for WalletPath
	tempDir := t.TempDir()

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
				WalletPath:     tempDir, // Use the temporary directory
				WalletPassword: "test_password",
				RPCNodes:       []string{"http://127.0.0.1:8545"},
				StorePath:      tempDir, // Use the temporary directory for StorePath as well
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
				RPCNodes:       []string{},
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

			// Assert the error
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.BasicCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
