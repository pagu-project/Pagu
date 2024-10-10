package voucher

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

const (
	CommandName           = "voucher"
	ClaimCommandName      = "claim"
	CreateOneCommandName  = "create-one"
	CreateBulkCommandName = "create-bulk"
	StatusCommandName     = "status"
	HelpCommandName       = "help"
)

type Voucher struct {
	db            repository.Database
	wallet        wallet.IWallet
	clientManager client.Manager
}

func NewVoucher(db repository.Database, wlt wallet.IWallet, cli client.Manager) *Voucher {
	return &Voucher{
		db:            db,
		wallet:        wlt,
		clientManager: cli,
	}
}

func (v *Voucher) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(v.db, v.wallet)

	subCmdClaim := &command.Command{
		Name: ClaimCommandName,
		Help: "Claim your voucher coins and bond to validator",
		Args: []command.Args{
			{
				Name:     "code",
				Desc:     "voucher code",
				InputBox: command.InputBoxText,
				Optional: false,
			},
			{
				Name:     "address",
				Desc:     "your pactus validator address",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.WalletBalance},
		Handler:     v.claimHandler,
		TargetFlag:  command.TargetMaskMain,
	}

	subCmdCreateOne := &command.Command{
		Name: CreateOneCommandName,
		Help: "Create a new voucher code",
		Args: []command.Args{
			{
				Name:     "amount",
				Desc:     "Amount of PAC to bond",
				InputBox: command.InputBoxAmount,
				Optional: false,
			},
			{
				Name:     "valid-months",
				Desc:     "Indicates how many months the voucher is valid after it is issued",
				InputBox: command.InputBoxNumber,
				Optional: false,
			},
			{
				Name:     "recipient",
				Desc:     "Indicates the name of the recipient of the voucher",
				InputBox: command.InputBoxText,
				Optional: true,
			},
			{
				Name:     "description",
				Desc:     "Describes the reason for issuing the voucher",
				InputBox: command.InputBoxText,
				Optional: true,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.OnlyModerator},
		Handler:     v.createOneHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	subCmdCreateBulk := &command.Command{
		Name: CreateBulkCommandName,
		Help: "Create more than one voucher code by importing file",
		Args: []command.Args{
			{
				Name:     "file",
				Desc:     "include list of vouchers receivers",
				InputBox: command.InputBoxFile,
				Optional: false,
			},
			{
				Name:     "notify",
				Desc:     "Notify receivers by sending mail",
				InputBox: command.InputBoxToggle,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.OnlyModerator},
		Handler:     v.createBulkHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	subCmdStatus := &command.Command{
		Name: StatusCommandName,
		Help: "Get status of vouchers/one voucher",
		Args: []command.Args{
			{
				Name:     "code",
				Desc:     "Voucher code (8 characters)",
				InputBox: command.InputBoxText,
				Optional: true,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.OnlyModerator},
		Handler:     v.statusHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdVoucher := &command.Command{
		Name:        CommandName,
		Help:        "Voucher Commands",
		Args:        nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMain | command.TargetMaskModerator,
	}

	cmdVoucher.AddSubCommand(subCmdClaim)
	cmdVoucher.AddSubCommand(subCmdCreateOne)
	cmdVoucher.AddSubCommand(subCmdCreateBulk)
	cmdVoucher.AddSubCommand(subCmdStatus)
	return cmdVoucher
}
