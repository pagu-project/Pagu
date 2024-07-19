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

func (v *Validator) importHandler(_ *entity.User, cmd *command.Command, args map[string]string) command.CommandResult {
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

	if len(records) < 2 {
		err = fmt.Errorf("no record founded. please add at least one record to csv file")
		return cmd.ErrorResult(err)
	}

	for rowIndex := 1; rowIndex < len(records); rowIndex++ {
		if len(records[rowIndex]) != 2 {
			err = fmt.Errorf("invalid data at row %d", rowIndex)
			return cmd.ErrorResult(err)
		}

		if rowIndex == 0 {
			continue
		}

		validator := &entity.Validator{
			Name:  records[rowIndex][0],
			Email: records[rowIndex][1],
		}

		if err = v.db.AddValidator(validator); err != nil {
			return cmd.ErrorResult(err)
		}
	}

	return cmd.SuccessfulResult("Validators created successfully!")
}
