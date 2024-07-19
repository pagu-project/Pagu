package voucher

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pactus-project/pactus/util/logger"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/pagu-project/Pagu/pkg/utils"
)

func (v *Voucher) createOneHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	code := utils.RandomString(8, utils.CapitalAlphanumerical)
	for _, err := v.db.GetVoucherByCode(code); err == nil; {
		code = utils.RandomString(8, utils.CapitalAlphanumerical)
	}

	amt, err := amount.FromString(args["amount"])
	if err != nil {
		return cmd.ErrorResult(errors.New("invalid amount param"))
	}

	maxStake, _ := amount.NewAmount(1000)
	if amt > maxStake {
		return cmd.ErrorResult(errors.New("stake amount is more than 1000"))
	}

	expireMonths, err := strconv.Atoi(args["valid-months"])
	if err != nil {
		return cmd.ErrorResult(errors.New("invalid valid-months param"))
	}

	vch := &entity.Voucher{
		Creator:     caller.ID,
		Code:        code,
		ValidMonths: uint8(expireMonths),
		Amount:      amt,
	}

	vch.Recipient = args["recipient"]
	vch.Desc = args["description"]

	err = v.db.AddVoucher(vch)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Voucher created successfully! \n Code: %s", vch.Code)
}

func (v *Voucher) createBulkHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
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

	if len(records) < 2 {
		err = fmt.Errorf("no record founded. please add at least one record to csv file")
		return cmd.ErrorResult(err)
	}

	vouchers := make([]*entity.Voucher, 0)
	for rowIndex := 1; rowIndex < len(records); rowIndex++ {
		code := utils.RandomString(8, utils.CapitalAlphanumerical)
		for _, err := v.db.GetVoucherByCode(code); err == nil; {
			code = utils.RandomString(8, utils.CapitalAlphanumerical)
		}

		email := records[rowIndex][0] // TODO: validate email address using regex
		recipient := records[rowIndex][1]
		amt, err := amount.FromString(strings.TrimSpace(records[rowIndex][2]))
		if err != nil {
			err = fmt.Errorf("invalid amount at row %d", rowIndex)
			return cmd.ErrorResult(err)
		}

		maxStake, _ := amount.NewAmount(1000)
		if amt > maxStake {
			return cmd.ErrorResult(errors.New("stake amount is more than 1000"))
		}

		validMonths, err := strconv.Atoi(strings.TrimSpace(records[rowIndex][3]))
		if err != nil {
			err = fmt.Errorf("invalid validate months at row %d", rowIndex)
			return cmd.ErrorResult(err)
		}

		desc := records[rowIndex][4]
		vouchers = append(vouchers, &entity.Voucher{
			Creator:     caller.ID,
			Code:        code,
			Desc:        desc,
			Recipient:   recipient,
			Email:       email,
			ValidMonths: uint8(validMonths),
			Amount:      amt,
		})
	}

	for _, vch := range vouchers {
		// TODO: add gorm transaction for this two insert
		err = v.db.AddVoucher(vch)
		if err != nil {
			return cmd.ErrorResult(err)
		}

		if args["notify"] == "TRUE" {
			err = v.db.AddNotification(&entity.Notification{
				Type:   entity.NotificationTypeEmail,
				Email:  vch.Email,
				Status: entity.NotificationStatusPending,
			})
			if err != nil {
				return cmd.ErrorResult(err)
			}
		}
	}

	return cmd.SuccessfulResult("Vouchers created successfully!")
}
