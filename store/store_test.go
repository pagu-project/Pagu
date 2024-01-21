package store_test

import (
	_ "embed"
	"os"
	"testing"
	"time"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/stretchr/testify/assert"
)

//go:embed store_example.json
var exampleStore []byte

func setup(t *testing.T) (store.IStore, string) {
	cfg, err := config.Load("../.env")
	assert.NoError(t, err)

	err = os.WriteFile(cfg.StorePath, exampleStore, 0o600)
	assert.NoError(t, err)

	log.InitGlobalLogger()
	logger := log.NewSubLogger("store_test")

	store, err := store.LoadStore(cfg, logger)
	assert.NoError(t, err)

	return store, cfg.StorePath
}

func TestStore(t *testing.T) {
	store, path := setup(t)

	t.Run("get claimer", func(t *testing.T) {
		claimer := store.ClaimerInfo("123456789")
		assert.Equal(t, int64(100), claimer.TotalReward)
		assert.Equal(t, "123456789", claimer.DiscordID)
	})

	t.Run("test add claim transaction", func(t *testing.T) {
		txID := "0x123456789"
		time := time.Now()
		data := "data"
		discordID := "123456789"

		claimer := store.ClaimerInfo(discordID)

		isClaimed := claimer.IsClaimed()
		assert.False(t, isClaimed)

		err := store.AddClaimTransaction(txID, claimer.TotalReward, time, data, discordID)
		assert.NoError(t, err)

		claimedInfo := store.ClaimerInfo(discordID)
		assert.Equal(t, data, claimedInfo.ClaimTransaction.Data)
		assert.Equal(t, discordID, claimedInfo.DiscordID)
		assert.Equal(t, int64(100), claimedInfo.ClaimTransaction.Amount)
		assert.Equal(t, txID, claimedInfo.ClaimTransaction.TxID)
		assert.Equal(t, time.Unix(), claimedInfo.ClaimTransaction.Time)
		assert.Equal(t, claimer.TotalReward, claimedInfo.ClaimTransaction.Amount)
		assert.Equal(t, int64(0), claimedInfo.TotalReward)

		isClaimed = claimedInfo.IsClaimed()
		assert.True(t, isClaimed)
	})

	err := os.Remove(path)
	assert.NoError(t, err)

	err = os.Remove("RoboPac.log")
	assert.NoError(t, err)
}
