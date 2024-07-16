package command

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/entity"
)

func (h *MiddlewareHandler) OnlyAdmin(cmd *Command, _ entity.AppID, _ string, _ map[string]string) error {
	if cmd.User.Role != entity.Admin {
		return errors.New("this command is Only Admin")
	}

	return nil
}

func (h *MiddlewareHandler) OnlyModerator(cmd *Command, _ entity.AppID, _ string, _ map[string]string) error {
	if cmd.User.Role != entity.Moderator {
		return errors.New("this command is Only Moderator")
	}

	return nil
}
