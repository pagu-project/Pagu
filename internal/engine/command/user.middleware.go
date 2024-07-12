package command

import (
	"github.com/pagu-project/Pagu/internal/entity"
)

func (h *MiddlewareHandler) CreateUser(cmd *Command, appID entity.AppID, callerID string, _ map[string]any) error {
	if user, _ := h.db.GetUserInApp(appID, callerID); user != nil {
		cmd.User = user
		return nil
	}

	user := &entity.User{ApplicationID: appID, CallerID: callerID}
	if err := h.db.AddUser(user); err != nil {
		return err
	}
	cmd.User = user
	return nil
}
