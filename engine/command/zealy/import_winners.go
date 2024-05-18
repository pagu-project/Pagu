package zealy

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/pagu-project/Pagu/database"
	"github.com/pagu-project/Pagu/engine/command"
)

/*
importWinnersHandler gives a csv file with below format and push the data into ZealyUser table.
Position| Discord User 	| 	Prize	.
1st		|	user_id_1	|	amount	.
2nd		|	user_id_2	|	amount	.
*/
func (z *Zealy) importWinnersHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {

	if len(args) == 0 {
		return cmd.FailedResult("please specify a file path to import")
	}

	path := args[0]
	records, err := readCSV(path)
	if err != nil {
		return cmd.ErrorResult(fmt.Errorf("csv file is not valid"))
	}

	totalInserted := 0

	for _, record := range records[1:] {
		discordID := record[1]
		if _, err := z.db.GetZealyUser(discordID); err == nil {
			return cmd.ErrorResult(fmt.Errorf("duplicate zealy user with discord ID: %s", discordID))
		}

		prizeStr := record[2]
		prize, _ := strconv.Atoi(prizeStr)
		if err := z.db.AddZealyUser(&database.ZealyUser{
			Amount:    int64(prize),
			DiscordID: discordID,
			TxHash:    "",
		}); err != nil {
			return cmd.ErrorResult(fmt.Errorf("adding zealy user into db with discord ID: %s", discordID))
		}

		totalInserted++
	}

	return cmd.SuccessfulResult("Imported successfully\nTotal inserted: %d", totalInserted)
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
