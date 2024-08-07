package voucher

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jszwec/csvutil"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/pagu-project/Pagu/pkg/notification"
	"github.com/pagu-project/Pagu/pkg/utils"
	"gorm.io/datatypes"
)

type BulkRecorder struct {
	Recipient        string  `csv:"Recipient"`
	Email            string  `csv:"Email"`
	Amount           float64 `csv:"Amount"`
	ValidatedInMonth int     `csv:"Validated"`
	Description      string  `csv:"Description"`
}

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

	csvReader := csv.NewReader(resp.Body)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		logger.Error(err.Error())
		return cmd.ErrorResult(errors.New("failed to read csv content"))
	}

	var records []BulkRecorder
	for {
		r := BulkRecorder{}
		if err = dec.Decode(&r); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			logger.Error(err.Error())
			return cmd.ErrorResult(errors.New("failed to parse csv content"))
		}

		records = append(records, r)
	}

	if len(records) == 0 {
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
			if v.createNotification(vch.Email, vch.Code, vch.Recipient, vch.Amount.ToPAC()) != nil {
				return cmd.ErrorResult(err)
			}
		}
	}

	return cmd.SuccessfulResult("Vouchers created successfully!")
}

func (v *Voucher) createBulkVoucher(records []BulkRecorder, callerID uint) ([]*entity.Voucher, error) {
	vouchers := make([]*entity.Voucher, 0)
	for index, record := range records {
		code := utils.RandomString(8, utils.CapitalAlphanumerical)
		for _, err := v.db.GetVoucherByCode(code); err == nil; {
			code = utils.RandomString(8, utils.CapitalAlphanumerical)
		}

		email := record.Email // TODO: validate email address using regex
		recipient := record.Recipient
		desc := record.Description

		amt, err := amount.NewAmount(record.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid amount at row %d", index+1)
		}

		maxStake, _ := amount.NewAmount(1000)
		if amt > maxStake {
			return nil, fmt.Errorf("stake amount is more than 1000")
		}

		validMonths := record.ValidatedInMonth
		if validMonths < 1 {
			return nil, fmt.Errorf("num of validated month of code must be greater than 0 at row %d", index+1)
		}

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

func (v *Voucher) createNotification(email, code, recipient string, amt float64) error {
	notificationData := entity.VoucherNotificationData{
		Code:      code,
		Recipient: recipient,
		Amount:    amt,
	}
	b, err := json.Marshal(notificationData)
	if err != nil {
		return err
	}
	voucherCodeJSON := datatypes.JSON(b)
	return v.db.AddNotification(&entity.Notification{
		Type:      notification.NotificationTypeMail,
		Status:    entity.NotificationStatusPending,
		Recipient: email,
		Data:      voucherCodeJSON,
	})
}
