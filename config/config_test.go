package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBasicCheck tests the BasicCheck method of the Config struct.
func TestBasicCheck(t *testing.T) {
	// Create a temporary directory for the WalletPath
	tempWalletPath := t.TempDir()

	// Define test cases
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Valid config",
			cfg: Config{
				Wallet: Wallet{
					Address:  "test_wallet_address",
					Path:     tempWalletPath, // Use the temporary directory
					Password: "test_password",
				},
				NetworkNodes: []string{"http://127.0.0.1:8545"},
				DiscordBot: DiscordBot{
					Token:   "MTEabc123",
					GuildID: "123456789",
				},
				Phoenix: PhoenixNetwork{
					NetworkNodes: []string{""},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid RPCNodes",
			cfg: Config{
				Wallet: Wallet{
					Address:  "test_wallet_address",
					Path:     "/valid/path",
					Password: "test_password",
				},
				NetworkNodes: []string{},
				DiscordBot: DiscordBot{
					Token:   "MTEabc123",
					GuildID: "123456789",
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
