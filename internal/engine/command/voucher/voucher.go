package voucher

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

const (
	CommandName       = "voucher"
	ClaimCommandName  = "claim"
	CreateCommandName = "create"
	HelpCommandName   = "help"
)

type Voucher struct {
	db            *repository.DB
	wallet        *wallet.Wallet
	clientManager *client.Mgr
}

func NewVoucher(db *repository.DB, wallet *wallet.Wallet, cli *client.Mgr) Voucher {
	return Voucher{
		db:            db,
		wallet:        wallet,
		clientManager: cli,
	}
}

func (v *Voucher) GetCommand() command.Command {
	subCmdClaim := command.Command{
		Name: ClaimCommandName,
		Help: "Claim your voucher coins and bond to validator",
		Args: []command.Args{
			{
				Name:     "code",
				Desc:     "voucher code",
				Optional: false,
			},
			{
				Name:     "address",
				Desc:     "your pactus validator address",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     v.claimHandler,
	}

	subCmdCreate := command.Command{
		Name: ClaimCommandName,
		Help: "Add a new voucher to database",
		Args: []command.Args{
			{
				Name:     "recipient",
				Desc:     "Indicates the name of the recipient of the voucher",
				Optional: false,
			},
			{
				Name:     "description",
				Desc:     "Describes the reason for issuing the voucher",
				Optional: false,
			},
			{
				Name:     "valid-months",
				Desc:     "Indicates how many months the voucher is valid after it is issued",
				Optional: false,
			},
			{
				Name:     "amount",
				Desc:     "Amount of PAC to bond",
				Optional: false,
			},
			{
				Name:     "discord-id",
				Desc:     "Recipient Discord ID",
				Optional: false,
			},
			{
				Name:     "code",
				Desc:     "The voucher code",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     v.createHandler,
	}

	cmdVoucher := command.Command{
		Name:        CommandName,
		Help:        "Voucher Commands",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMain,
	}

	cmdVoucher.AddSubCommand(subCmdClaim)
	cmdVoucher.AddSubCommand(subCmdCreate)
	return cmdVoucher
}
