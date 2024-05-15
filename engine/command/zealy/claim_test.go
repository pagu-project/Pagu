package zealy

import (
	"os"
	"testing"

	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/database"
	"github.com/pagu-project/Pagu/engine/command"
	"github.com/pagu-project/Pagu/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupClaimTest(t *testing.T) *Zealy {
	t.Helper()

	file, err := os.CreateTemp("", "temp-db")
	require.NoError(t, err)

	db, err := database.NewDB(file.Name())
	require.NoError(t, err)

	walletConfig := &config.Wallet{
		Enable:   false,
		Address:  "tpc1zzgvtgd8p6mlwey5e4ajpg5ugn8zltwk2eawfpm",
		RPCUrl:   "localhost:50052",
		Path:     "../../../test/test_wallet",
		Password: "123456789",
	}

	w := wallet.Open(walletConfig)
	zealy := NewZealy(db, w)
	return &zealy
}

func TestClaim(t *testing.T) {
	t.Run("everything normal and good", func(t *testing.T) {
		zealy := setupClaimTest(t)

		discordID := "abcd1234"
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: discordID,
			IsClaimed: false,
			TxHash:    "",
		})

		cmd := zealy.GetCommand()
		expectedRes := zealy.claimHandler(cmd, command.AppIdDiscord, discordID, "tpc1z8da0suum050nqsua0jjwelc62ppq6nuz0msvvw")

		assert.Equal(t, true, expectedRes.Successful)
		assert.Equal(t, "", expectedRes.Error)

		user, err := zealy.db.GetZealyUser(discordID)
		assert.NoError(t, err)
		assert.NotEqual(t, "", user.TxHash)
		assert.Equal(t, true, user.IsClaimed)
	})

	t.Run("should fail, user not found", func(t *testing.T) {
		zealy := setupClaimTest(t)

		discordID := "12345678"
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: discordID,
			IsClaimed: false,
			TxHash:    "",
		})

		cmd := zealy.GetCommand()
		expectedRes := zealy.claimHandler(cmd, command.AppIdDiscord, "some_other_discord_id", "tpc1z8da0suum050nqsua0jjwelc62ppq6nuz0msvvw")

		assert.Equal(t, false, expectedRes.Successful)
	})

	t.Run("should fail, claimed before", func(t *testing.T) {
		zealy := setupClaimTest(t)

		discordID := "12345678"
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: discordID,
			IsClaimed: true,
			TxHash:    "",
		})

		cmd := zealy.GetCommand()
		expectedRes := zealy.claimHandler(cmd, command.AppIdDiscord, discordID, "tpc1z8da0suum050nqsua0jjwelc62ppq6nuz0msvvw")

		assert.Equal(t, false, expectedRes.Successful)
	})

	t.Run("should fail, no wallet address provided", func(t *testing.T) {
		zealy := setupClaimTest(t)

		discordID := "12345678"
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: discordID,
			IsClaimed: false,
			TxHash:    "",
		})

		cmd := zealy.GetCommand()
		expectedRes := zealy.claimHandler(cmd, command.AppIdDiscord, discordID)

		assert.Equal(t, false, expectedRes.Successful)
	})

	t.Run("should fail, no wallet address provided", func(t *testing.T) {
		zealy := setupClaimTest(t)

		discordID := "12345678"
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: discordID,
			IsClaimed: false,
			TxHash:    "",
		})

		cmd := zealy.GetCommand()
		expectedRes := zealy.claimHandler(cmd, command.AppIdDiscord, discordID, "tpc1z8da0suum050nqsua0jjwelc62ppq6nuz0msvvw")

		assert.Equal(t, false, expectedRes.Successful)
	})
}
