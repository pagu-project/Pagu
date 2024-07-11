package phoenix

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

//nolint:nolintlint // remove me after I am used
func (pt *Phoenix) walletHandler(cmd command.Command, _ entity.AppID, _ string, _ ...string) command.CommandResult {
	return cmd.SuccessfulResult("Pagu Phoenix Address: %s\nBalance: %d", pt.wallet.Address(), pt.wallet.Balance())
}
