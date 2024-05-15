package zealy

import (
	"fmt"
	"os"
	"testing"

	"github.com/pagu-project/Pagu/database"
	"github.com/pagu-project/Pagu/engine/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportWinnersHandler(t *testing.T) {
	t.Run("import list of winners successfully", func(t *testing.T) {
		zealy := setup(t)
		tempPath := "temp-csv"
		csvData := "Position,Discord ID,Prize\n1,id1,1\n2,id2,2\n3,id3,3"
		tempFile := createTempFile(t, tempPath, csvData)

		cmd := zealy.GetCommand()
		expectedRes := zealy.importWinnersHandler(cmd, command.AppIdCLI, "", tempFile.Name())

		assert.Equal(t, true, expectedRes.Successful)
		assert.Equal(t, "Imported successfully\nTotal inserted: 3\nTotal duplicate: 0", expectedRes.Message)

		users, err := zealy.db.GetAllZealyUser()
		assert.Equal(t, nil, err)
		for i, u := range users {
			assert.Equal(t, false, u.IsClaimed)
			assert.Equal(t, fmt.Sprintf("id%d", i+1), u.DiscordID)
			assert.Equal(t, int64(i+1), u.Amount)
		}
	})

	t.Run("import list with duplicate items", func(t *testing.T) {
		zealy := setup(t)
		tempPath := "temp-csv"
		csvData := "Position,Discord ID,Prize\n1,id1,1\n2,id2,2\n3,id3,3\n4,id3,4"
		tempFile := createTempFile(t, tempPath, csvData)

		cmd := zealy.GetCommand()
		expectedRes := zealy.importWinnersHandler(cmd, command.AppIdCLI, "", tempFile.Name())

		assert.Equal(t, true, expectedRes.Successful)
		assert.Equal(t, "Imported successfully\nTotal inserted: 3\nTotal duplicate: 1", expectedRes.Message)

		users, err := zealy.db.GetAllZealyUser()
		assert.Equal(t, nil, err)
		for i, u := range users {
			assert.Equal(t, false, u.IsClaimed)
			assert.Equal(t, fmt.Sprintf("id%d", i+1), u.DiscordID)
			assert.Equal(t, int64(i+1), u.Amount)
		}
	})

	t.Run("not implemented in all apps", func(t *testing.T) {
		zealy := setup(t)
		tempPath := "temp-csv"
		csvData := "Position,Discord ID,Prize\n1,id1,100\n2,id2,100\n3,id3,100"
		tempFile := createTempFile(t, tempPath, csvData)

		cmd := zealy.GetCommand()
		expectedRes := zealy.importWinnersHandler(cmd, command.AppIdDiscord, "", tempFile.Name())

		assert.Equal(t, false, expectedRes.Successful)
		assert.Equal(t, "command not implemented for this app", expectedRes.Message)
	})

	t.Run("no csv file passed", func(t *testing.T) {
		zealy := setup(t)

		cmd := zealy.GetCommand()
		expectedRes := zealy.importWinnersHandler(cmd, command.AppIdCLI, "")

		assert.Equal(t, false, expectedRes.Successful)
		assert.Equal(t, "please specify a file path to import", expectedRes.Message)
	})
}

func setup(t *testing.T) *Zealy {
	t.Helper()

	dbFile, err := os.CreateTemp("", "temp-db")
	require.NoError(t, err)
	db, err := database.NewDB(dbFile.Name())
	require.NoError(t, err)

	zealy := NewZealy(db, nil)
	return &zealy
}

func createTempFile(t *testing.T, path, data string) *os.File {
	t.Helper()

	tempFile, err := os.CreateTemp("", path)
	if err != nil {
		t.Fatalf("createTempFile() error = %v", err)
	}

	defer tempFile.Close()

	if _, err := tempFile.Write([]byte(data)); err != nil {
		t.Fatalf("createTempFile() error writing data to temp file: %v", err)
	}

	return tempFile
}
