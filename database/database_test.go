package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *DB {
	file, err := os.CreateTemp("", "temp-db")
	require.NoError(t, err)

	db, err := NewDB(file.Name())
	require.NoError(t, err)

	return db
}

func TestMember(t *testing.T) {
	db := setup(t)

	discordID := "123456789"
	user := &DiscordUser{
		DiscordID:      discordID,
		DepositAddress: "pc1zrandomaddr",
	}

	t.Run("test add user", func(t *testing.T) {
		err := db.AddUser(user)
		require.NoError(t, err)
	})

	t.Run("test get not existing user", func(t *testing.T) {
		u, err := db.GetUser("non-existing-member-id")
		require.Error(t, err)
		require.Nil(t, u)
	})

	t.Run("test get user", func(t *testing.T) {
		u, err := db.GetUser(discordID)
		require.NoError(t, err)

		assert.Equal(t, user.DepositAddress, u.DepositAddress)
		assert.Equal(t, user.DiscordID, u.DiscordID)
	})
}

func TestHasUser(t *testing.T) {
	db := setup(t)

	err := db.AddUser(&DiscordUser{
		DiscordID: "123456",
	})
	assert.NoError(t, err)

	assert.True(t, db.HasUser("123456"))
	assert.False(t, db.HasUser("654321"))
}

func TestDB_CreateOffer(t *testing.T) {
	db := setup(t)

	someDiscordId := "123456"
	someAddr := "some-addr"
	_ = db.AddUser(&DiscordUser{
		DiscordID:      someDiscordId,
		DepositAddress: someAddr,
	})

	u, _ := db.GetUser(someDiscordId)

	offer := &Offer{
		TotalAmount: 10,
		TotalPrice:  10,
		ChainType:   "BTCUSDT",
		Address:     "addr1",
		DiscordUser: *u,
	}
	err := db.CreateOffer(offer)

	assert.NoError(t, err)

	var actual Offer
	err = db.Preload("DiscordUser").First(&actual, 1).Error
	assert.NoError(t, err)

	assert.Equal(t, offer.ChainType, actual.ChainType)
	assert.Equal(t, offer.Address, actual.Address)
	assert.Equal(t, u.DepositAddress, offer.DiscordUser.DepositAddress)
}
