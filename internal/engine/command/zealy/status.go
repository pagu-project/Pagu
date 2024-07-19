package zealy

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/amount"
)

func (z *Zealy) statusHandler(_ *entity.User, cmd *command.Command, _ map[string]string) command.CommandResult {
	allUsers, err := z.db.GetAllZealyUser()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	total := 0
	totalClaimed := 0
	totalNotClaimed := 0

	totalAmount := amount.Amount(0)
	totalClaimedAmount := amount.Amount(0)
	totalNotClaimedAmount := amount.Amount(0)

	for _, u := range allUsers {
		total++
		totalAmount += u.Amount

		if u.IsClaimed() {
			totalClaimed++
			totalClaimedAmount += u.Amount
		} else {
			totalNotClaimed++
			totalNotClaimedAmount += u.Amount
		}
	}

	return cmd.SuccessfulResult("Total Users: %v\nTotal Claims: %v\nTotal not remained claims: %v\nTotal Coins: %v PAC\n"+
		"Total claimed coins: %v PAC\nTotal not claimed coins: %v PAC\n", total, totalClaimed, totalNotClaimed,
		totalAmount.ToPAC(), totalClaimedAmount.ToPAC(), totalNotClaimedAmount.ToPAC(),
	)
}
