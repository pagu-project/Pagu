package phoenix

import "github.com/pagu-project/Pagu/internal/engine/command"

func (pt *Phoenix) walletHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	return cmd.SuccessfulResult("Pagu Phoenix Address: %s\nBalance: %d", pt.wallet.Address(), pt.wallet.Balance())
}
