package validator

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (v *Validator) importHandler(cmd *command.Command, _ entity.AppID, _ string, _ map[string]any,
) command.CommandResult {
	return cmd.SuccessfulResult("Validator created successfully!")
}
