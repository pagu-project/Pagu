package job

import (
	"context"
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
}

func NewMailSender(db repository.Database, mailSender notification.ISender) Job {
	ctx, cancel := context.WithCancel(context.Background())
	return &mailSenderJob{
		ticker:     time.NewTicker(10 * time.Minute),
		ctx:        ctx,
		cancel:     cancel,
		db:         db,
		mailSender: mailSender,
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

	tmpl, _ := notification.LoadMailTemplate("./templates/voucher.html")
	err = p.mailSender.SendTemplateMail(
		notification.NotificationProviderZapToMail,
		"no-reply@pactus.org", []string{notif.Recipient}, tmpl, notif.Data)
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
