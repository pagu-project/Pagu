package zealy

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (z *Zealy) statusHandler(cmd command.Command, _ entity.AppID, _ string, args ...string) command.CommandResult {
	allUsers, err := z.db.GetAllZealyUser()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	total := 0
	totalClaimed := 0
	totalNotClaimed := 0

	totalAmount := 0
	totalClaimedAmount := 0
	totalNotClaimedAmount := 0

	for _, u := range allUsers {
		total++
		totalAmount += int(u.Amount)

		if u.IsClaimed() {
			totalClaimed++
			totalClaimedAmount += int(u.Amount)
		} else {
			totalNotClaimed++
			totalNotClaimedAmount += int(u.Amount)
		}
	}

	return cmd.SuccessfulResult("Total Users: %v\nTotal Claims: %v\nTotal not remained claims: %v\nTotal Coins: %v PAC\n"+
		"Total claimed coins: %v PAC\nTotal not claimed coins: %v PAC\n", total, totalClaimed, totalNotClaimed,
		totalAmount, totalClaimedAmount, totalNotClaimedAmount,
	)
}
