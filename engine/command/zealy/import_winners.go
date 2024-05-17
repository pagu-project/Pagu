package zealy

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/pagu-project/Pagu/database"
	"github.com/pagu-project/Pagu/engine/command"
)

func (z *Zealy) importWinnersHandler(cmd command.Command, appID command.AppID, _ string, args ...string) command.CommandResult {
	if appID != command.AppIdCLI {
		return cmd.FailedResult("command not implemented for this app")
	}

	if len(args) == 0 {
		return cmd.FailedResult("please specify a file path to import")
	}

	path := args[0]
	records, err := readCSV(path)
	if err != nil {
		return cmd.FailedResult("csv file is not valid")
	}

	totalInserted := 0
	totalDuplicate := 0

	for _, record := range records[1:] {
		discordID := record[1]
		if _, err := z.db.GetZealyUser(discordID); err == nil {
			totalDuplicate++
			continue
		}

		prizeStr := record[2]
		prize, _ := strconv.Atoi(prizeStr)
		if err := z.db.AddZealyUser(&database.ZealyUser{
			Amount:    int64(prize),
			DiscordID: discordID,
			TxHash:    "",
		}); err != nil {
			return cmd.FailedResult("error in adding zealy user into db. discord ID: %s", discordID)
		}

		totalInserted++
	}

	return cmd.SuccessfulResult("Imported successfully\nTotal inserted: %d\nTotal duplicate: %d", totalInserted, totalDuplicate)
}

func readCSV(path string) ([][]string, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
