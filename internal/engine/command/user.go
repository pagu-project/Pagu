package command

import (
	"github.com/pagu-project/Pagu/internal/entity"
)

func (h *MiddlewareHandler) CreateUser(cmd *Command, appID entity.AppID, callerID string, _ ...string) error {
	if !h.db.HasUserInApp(appID, callerID) {
		user := &entity.User{ApplicationID: appID, CallerID: callerID}
		if err := h.db.AddUser(user); err != nil {
			return err
		}

		cmd.User = user
	}

	return nil
}
