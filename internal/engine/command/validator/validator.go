package validator

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

const (
	CommandName       = "validator"
	ImportCommandName = "import"
	HelpCommandName   = "help"
)

type Validator struct {
	db            repository.Database
	wallet        wallet.IWallet
	clientManager client.Manager
}

func NewValidator(db repository.Database, wlt wallet.IWallet, cli client.Manager) *Validator {
	return &Validator{
		db:            db,
		wallet:        wlt,
		clientManager: cli,
	}
}

func (v *Validator) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(v.db, v.wallet)

	subCmdImport := &command.Command{
		Name: ImportCommandName,
		Help: "Import list of validator",
		Args: []command.Args{
			{
				Name:     "file",
				Desc:     "include list of validators",
				Type:     command.CommandParamTypeAttachment,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: []command.MiddlewareFunc{middlewareHandler.CreateUser, middlewareHandler.WalletBalance},
		Handler:     v.importHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdValidator := &command.Command{
		Name:        CommandName,
		Help:        "Validator Commands",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdValidator.AddSubCommand(subCmdImport)
	return cmdValidator
}
