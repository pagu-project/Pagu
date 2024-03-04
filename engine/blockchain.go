package engine

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
)

const (
	BlockChainCommandName = "blockchain"
	CalcRewardCommandName = "calc-reward"
	CalcFeeCommandName    = "calc-fee"
)

func (be *BotEngine) RegisterBlockchainCommands() {
	cmdCalcReward := Command{
		Name: CalcRewardCommandName,
		Desc: "claculate how much PAC coins you will earn with your validator stakes",
		Help: "",
		Args: []Args{
			{
				Name:     "stake-amount",
				Desc:     "amount of stake in your validator (1-1000)",
				Optional: false,
			},
			{
				Name:     "time-interval",
				Desc:     "after one: day | month | year",
				Optional: true,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.calcRewardHandler,
	}

	cmdBlockchain := Command{
		Name:        BlockChainCommandName,
		Desc:        "Blockchain information and tools",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&cmdCalcReward},
		Handler:     nil,
	}

	be.Cmds = append(be.Cmds, cmdBlockchain)
}

func (be *BotEngine) calcRewardHandler(_ AppID, _ string, args ...string) (*CommandResult, error) {
	stake, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	time := args[1]

	if stake < 1 || stake > 1_000 {
		return nil, errors.New("minimum of stake is 1 PAC and maximum is 1,000 PAC")
	}

	var blocks int
	switch time {
	case "day":
		blocks = 8640
	case "month":
		blocks = 259200
	case "year":
		blocks = 3110400
	default:
		blocks = 8640
		time = "day"
	}

	bi, err := be.clientMgr.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}

	reward := int64(stake*blocks) / int64(util.ChangeToCoin(bi.TotalPower))

	result := fmt.Sprintf("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator in one %s ⏰ with %v PAC total power ⚡ of committee."+
		"\n\n> Note📝: This is an estimation and the number can get changed by changes of your stake amount, total power and ...",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), time, utils.FormatNumber(bi.TotalPower))

	return &CommandResult{
		Successful: true,
		Message:    result,
	}, nil
}
