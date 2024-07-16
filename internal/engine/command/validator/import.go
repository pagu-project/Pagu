package validator

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"

	"github.com/pactus-project/pactus/util/logger"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (v *Validator) importHandler(cmd *command.Command, _ entity.AppID, _ string, args map[string]string,
) command.CommandResult {
	fileURL := args["file"]

	httpClient := new(http.Client)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fileURL, http.NoBody)
	if err != nil {
		logger.Error(err.Error())
		return cmd.ErrorResult(errors.New("failed to fetch attachment content"))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return cmd.ErrorResult(errors.New("failed to fetch attachment content"))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	r := csv.NewReader(resp.Body)
	records, err := r.ReadAll()
	if err != nil {
		logger.Error(err.Error())
		return cmd.ErrorResult(errors.New("failed to read attachment content"))
	}

	for rowIndex, row := range records {
		if len(row) != 2 {
			err = fmt.Errorf("invalid data at row %d", rowIndex)
			return cmd.ErrorResult(err)
		}

		validator := &entity.Validator{
			Name:  row[0],
			Email: row[1],
		}

		if err = v.db.AddValidator(validator); err != nil {
			return cmd.ErrorResult(err)
		}
	}

	return cmd.SuccessfulResult("Validators created successfully!")
}
