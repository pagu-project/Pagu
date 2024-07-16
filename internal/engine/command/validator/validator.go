package validator

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
)

const (
	CommandName       = "validator"
	ImportCommandName = "import"
	HelpCommandName   = "help"
)

type Validator struct {
	db repository.Database
}

func NewValidator(db repository.Database) *Validator {
	return &Validator{
		db: db,
	}
}

func (v *Validator) GetCommand() *command.Command {
	subCmdImport := &command.Command{
		Name: ImportCommandName,
		Help: "Import list of validator",
		Args: []command.Args{
			{
				Name:     "file",
				Desc:     "include list of validators",
				InputBox: command.InputBoxFile,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: nil,
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
