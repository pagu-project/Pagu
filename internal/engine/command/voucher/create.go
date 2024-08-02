package voucher

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pactus-project/pactus/util/logger"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/pagu-project/Pagu/pkg/notification"
	"github.com/pagu-project/Pagu/pkg/utils"
	"gorm.io/datatypes"
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
	notify := args["notify"]

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

	vouchers, err := v.createBulkVoucher(records, caller.ID)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	for _, vch := range vouchers {
		// TODO: add gorm transaction for this two insert
		err := v.db.AddVoucher(vch)
		if err != nil {
			return cmd.ErrorResult(err)
		}

		if notify == "TRUE" {
			if v.createNotification(vch.Email, vch.Code) != nil {
				return cmd.ErrorResult(err)
			}
		}
	}

	return cmd.SuccessfulResult("Vouchers created successfully!")
}

func (v *Voucher) createBulkVoucher(records [][]string, callerID uint) ([]*entity.Voucher, error) {
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
			return nil, fmt.Errorf("invalid amount at row %d", rowIndex)
		}

		maxStake, _ := amount.NewAmount(1000)
		if amt > maxStake {
			return nil, fmt.Errorf("stake amount is more than 1000")
		}

		validMonths, err := strconv.Atoi(strings.TrimSpace(records[rowIndex][3]))
		if err != nil {
			return nil, fmt.Errorf("invalid validate months at row %d", rowIndex)
		}

		desc := records[rowIndex][4]
		vouchers = append(vouchers, &entity.Voucher{
			Creator:     callerID,
			Code:        code,
			Desc:        desc,
			Recipient:   recipient,
			Email:       email,
			ValidMonths: uint8(validMonths),
			Amount:      amt,
		})
	}

	return vouchers, nil
}

func (v *Voucher) createNotification(email, code string) error {
	notificationData := entity.VoucherNotificationData{Code: code}
	b, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}
	voucherCodeJSON := datatypes.JSON(b)
	return v.db.AddNotification(&entity.Notification{
		Type:      notification.NotificationTypeMail,
		Recipient: email,
		Data:      voucherCodeJSON,
		Status:    entity.NotificationStatusPending,
	})
}
