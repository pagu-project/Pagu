package phoenix

import (
	"time"

	"github.com/pagu-project/Pagu/engine/command"
	"github.com/pagu-project/Pagu/utils"
)

func (pt *Phoenix) networkHealthHandler(cmd command.Command, _ command.AppID, _ string, _ ...string) command.CommandResult {
	lastBlockTime, lastBlockHeight := pt.clientMgr.GetLastBlockTime()
	lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0).Format("02/01/2006, 15:04:05")
	currentTime := time.Now()

	timeDiff := (currentTime.Unix() - int64(lastBlockTime))

	healthStatus := true
	if timeDiff > 15 {
		healthStatus = false
	}

	var status string
	if healthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return cmd.SuccessfulResult("Network is %s\nCurrentTime: %v\nLastBlockTime: %v\nTime Diff: %v\nLast Block Height: %v",
		status, currentTime.Format("02/01/2006, 15:04:05"), lastBlockTimeFormatted, timeDiff, utils.FormatNumber(int64(lastBlockHeight)))
}
