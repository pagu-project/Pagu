package store_test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/stretchr/testify/assert"
)

//go:embed test/store_example.json
var exampleStore []byte

func setup(t *testing.T) (store.IStore, string) {
	cfg, err := config.Load("./test/.env.test")
	assert.NoError(t, err)

	err = os.WriteFile(cfg.StorePath, exampleStore, 0o600)
	assert.NoError(t, err)

	log.InitGlobalLogger()
	logger := log.NewSubLogger("store_test")

	store, err := store.NewStore(cfg, logger)
	assert.NoError(t, err)

	return store, cfg.StorePath
}

func TestStore(t *testing.T) {
	store, path := setup(t)

	t.Run("unknown claimer", func(t *testing.T) {
		claimer := store.ClaimerInfo("unknown-addr")
		assert.Nil(t, claimer)
	})

	t.Run("get claimer", func(t *testing.T) {
		claimer := store.ClaimerInfo("tpc1pqn7uaeduklpg00rqt6uq0m9wy5txnyt0kmxmgf")
		assert.False(t, claimer.IsClaimed())
		assert.Equal(t, int64(100*1e9), claimer.TotalReward)
		assert.Equal(t, "123456789", claimer.DiscordID)
	})

	t.Run("test add claim transaction", func(t *testing.T) {
		txID := "0x123456789"
		discordID := "123456789"
		testNetValAddr := "tpc1pqn7uaeduklpg00rqt6uq0m9wy5txnyt0kmxmgf"

		claimer := store.ClaimerInfo(testNetValAddr)

		isClaimed := claimer.IsClaimed()
		assert.False(t, isClaimed)

		err := store.AddClaimTransaction(testNetValAddr, txID)
		assert.NoError(t, err)

		claimedInfo := store.ClaimerInfo(testNetValAddr)
		assert.Equal(t, discordID, claimedInfo.DiscordID)
		assert.Equal(t, int64(100*1e9), claimedInfo.TotalReward)
		assert.Equal(t, txID, claimedInfo.ClaimedTxID)

		isClaimed = claimedInfo.IsClaimed()
		assert.True(t, isClaimed)
	})

	t.Run("is claimed test", func(t *testing.T) {
		claimer := store.ClaimerInfo("tpc1pesz6kuv7jts6al6la3794fyj5xaj7wm93k7z6y")
		assert.Equal(t, int64(12*1e9), claimer.TotalReward)
		assert.Equal(t, "964550933793103912", claimer.DiscordID)
		assert.True(t, claimer.IsClaimed())
	})

	err := os.Remove(path)
	assert.NoError(t, err)

	err = os.Remove("RoboPac.log")
	assert.NoError(t, err)
}
