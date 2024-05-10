package database

import (
	"fmt"
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

func TestUserAndFaucet(t *testing.T) {
	db := setup(t)

	err := db.AddUser(&User{
		ID: "123456789",
	})
	assert.NoError(t, err)

	u, err := db.GetUser("123456789")
	assert.NoError(t, err)
	assert.Equal(t, "123456789", u.ID)

	r := db.CanGetFaucet("123456789")
	assert.True(t, r)

	err = db.AddFaucet(&Faucet{
		Address: "tpc1zlymfcuxlgvvuud2q4zw0scllqn74d2f90hld6w",
		Amount:  5,
		UserID:  "123456789",
	})
	assert.NoError(t, err)

	r = db.CanGetFaucet("123456789")
	assert.False(t, r)

	u, err = db.GetUser("not-exist")
	fmt.Println(u.ID)
	assert.Error(t, err)
}

func TestZealyDB(t *testing.T) {
	db := setup(t)

	err := db.AddZealyUser(&ZealyUser{
		Amount:    100,
		DiscordID: "12345678",
		IsClaimed: false,
		TxHash:    "",
	})
	assert.NoError(t, err)

	uz, err := db.GetZealyUser("12345678")
	assert.NoError(t, err)
	assert.Equal(t, false, uz.IsClaimed)
	assert.Equal(t, "", uz.TxHash)
	assert.Equal(t, int64(100), uz.Amount)

	err = db.UpdateZealyUser("12345678", "0x123456789")
	assert.NoError(t, err)

	uz, err = db.GetZealyUser("12345678")
	assert.NoError(t, err)
	assert.Equal(t, true, uz.IsClaimed)
	assert.Equal(t, "0x123456789", uz.TxHash)
	assert.Equal(t, int64(100), uz.Amount)

	_, err = db.GetZealyUser("87654321")
	assert.Error(t, err)

	azu, err := db.GetAllZealyUser()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(azu))
}
