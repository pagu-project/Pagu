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
	user := &User{
		DiscordID:      discordID,
		Address:        "pc1zrandomaddr",
		OpenOffers:     10,
		HasOpenPayment: false,
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

		assert.Equal(t, user.Address, u.Address)
		assert.Equal(t, user.DiscordID, u.DiscordID)
		assert.Equal(t, user.OpenOffers, u.OpenOffers)
		assert.Equal(t, user.HasOpenPayment, u.HasOpenPayment)
	})
}
