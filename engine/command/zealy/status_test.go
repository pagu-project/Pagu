package zealy

import (
	"os"
	"testing"

	"github.com/pagu-project/Pagu/database"
	"github.com/pagu-project/Pagu/engine/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStatusTest(t *testing.T) *Zealy {
	t.Helper()

	file, err := os.CreateTemp("", "temp-db")
	require.NoError(t, err)

	db, err := database.NewDB(file.Name())
	require.NoError(t, err)
	zealy := NewZealy(db, nil)
	return &zealy
}

func TestStatus(t *testing.T) {
	t.Run("everything normal and good", func(t *testing.T) {
		zealy := setupStatusTest(t)

		// claimed
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: "userID1",
			IsClaimed: true,
			TxHash:    "txHash1",
		})

		// claimed
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: "userID2",
			IsClaimed: true,
			TxHash:    "txHash2",
		})

		// claimed
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: "userID3",
			IsClaimed: true,
			TxHash:    "txHash3",
		})

		// not claimed
		_ = zealy.db.AddZealyUser(&database.ZealyUser{
			Amount:    100,
			DiscordID: "userID4",
			IsClaimed: false,
			TxHash:    "txHash4",
		})

		cmd := zealy.GetCommand()
		expectedRes := zealy.statusHandler(cmd, command.AppIdDiscord, "")

		assert.Equal(t, true, expectedRes.Successful)
		assert.Equal(t, "Total Users: 4\nTotal Claims: 3\nTotal not remained claims: 1\nTotal Coins: 400 PAC\n"+
			"Total claimed coins: 300 PAC\nTotal not claimed coins: 100 PAC\n", expectedRes.Message)
	})
}
