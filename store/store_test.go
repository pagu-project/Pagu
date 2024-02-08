package store_test

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func setup(t *testing.T) store.IStore {
	tempDir, err := os.MkdirTemp("", "RoboPAC")
	require.NoError(t, err)

	_, err = copy("./test/.env.test", path.Join(tempDir, "/.env"))
	require.NoError(t, err)

	_, err = copy("./test/store_example.json", path.Join(tempDir, "/store_example.json"))
	require.NoError(t, err)

	_, err = copy("./test/twitter_example.json", path.Join(tempDir, "/twitter_example.json"))
	require.NoError(t, err)

	_, err = copy("./test/test_wallet.json", path.Join(tempDir, "/test_wallet.json"))
	require.NoError(t, err)

	os.Chdir(tempDir)

	cfg, err := config.Load(".env")
	require.NoError(t, err)

	log.InitGlobalLogger()
	logger := log.NewSubLogger("store_test")

	store, err := store.NewStore(cfg, logger)
	require.NoError(t, err)

	return store
}

func TestStore(t *testing.T) {
	store := setup(t)

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
}
