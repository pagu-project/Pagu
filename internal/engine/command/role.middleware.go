package command

import (
	"errors"

	"github.com/pagu-project/Pagu/internal/entity"
)

func (h *MiddlewareHandler) OnlyAdmin(caller *entity.User, _ *Command, _ map[string]string) error {
	if caller.Role != entity.Admin {
		return errors.New("this command is Only Admin")
	}

	return nil
}

func (h *MiddlewareHandler) OnlyModerator(caller *entity.User, _ *Command, _ map[string]string) error {
	if caller.Role != entity.Moderator {
		return errors.New("this command is Only Moderator")
	}

	return nil
}
