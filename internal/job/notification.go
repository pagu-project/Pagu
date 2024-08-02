package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pactus-project/pactus/util/logger"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/notification"
)

type mailSenderJob struct {
	ctx        context.Context
	ticker     *time.Ticker
	cancel     context.CancelFunc
	db         repository.Database
	mailSender notification.ISender
	templates  map[string]string
}

func NewMailSender(db repository.Database, mailSender notification.ISender, templates map[string]string) Job {
	ctx, cancel := context.WithCancel(context.Background())
	return &mailSenderJob{
		ticker:     time.NewTicker(10 * time.Minute),
		ctx:        ctx,
		cancel:     cancel,
		db:         db,
		mailSender: mailSender,
		templates:  templates,
	}
}

func (p *mailSenderJob) Start() {
	p.sendVoucherNotifications()
	go p.runTicker()
}

func (p *mailSenderJob) sendVoucherNotifications() {
	notif, err := p.db.GetPendingMailNotification()
	if err != nil {
		logger.Error("failed to get pending mail from db", "err", err)
		return
	}

	v := entity.VoucherNotificationData{}
	vByte, _ := notif.Data.MarshalJSON()
	_ = json.Unmarshal(vByte, &v)
	tmpl, err := notification.LoadMailTemplate(p.templates["voucher"])
	if err != nil {
		logger.Fatal("failed to load mail template", "err", err)
	}

	err = p.mailSender.SendTemplateMail(
		notification.NotificationProviderZapToMail,
		"no-reply@pactus.org", []string{notif.Recipient}, tmpl, v)
	if err != nil {
		logger.Error("failed to send mail notification", "err", err)
		err = p.db.UpdateNotificationStatus(notif.ID, entity.NotificationStatusFail)
		if err != nil {
			logger.Error("failed to update status of sent mail", "err", err)
		}
	} else {
		err = p.db.UpdateNotificationStatus(notif.ID, entity.NotificationStatusDone)
		if err != nil {
			logger.Error("failed to update status of sent mail", "err", err)
		}
	}
}

func (p *mailSenderJob) runTicker() {
	for {
		select {
		case <-p.ctx.Done():
			return

		case <-p.ticker.C:
			p.sendVoucherNotifications()
		}
	}
}

func (p *mailSenderJob) Stop() {
	p.ticker.Stop()
}
