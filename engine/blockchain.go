package engine

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kehiy/RoboPac/utils"
	"github.com/pactus-project/pactus/util"
)

const (
	BlockChainCommandName     = "blockchain"
	CalcRewardCommandName     = "calc-reward"
	CalcFeeCommandName        = "calc-fee"
	BlockChainHelpCommandName = "help"
)

func (be *BotEngine) RegisterBlockchainCommands() {
	subCmdCalcReward := Command{
		Name: CalcRewardCommandName,
		Desc: "Calculate how much PAC coins you will earn with your validator stakes",
		Help: "Provide an stake amount between 1 to 100, please avoid using float numbers like: 1.9 or PAC prefix",
		Args: []Args{
			{
				Name:     "stake-amount",
				Desc:     "amount of stake in your validator (1-1000)",
				Optional: false,
			},
			{
				Name:     "time-interval",
				Desc:     "after one: day | month | year",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		Handler:     be.calcRewardHandler,
	}

	subCmdHelp := Command{
		Name: BlockChainHelpCommandName,
		Desc: "Help for Blockchain commands",
		Help: "Provide the sub command name as parameter",
		Args: []Args{
			{
				Name:     "sub-command",
				Desc:     "the subcommand you want to see the related help of it. (optional)",
				Optional: true,
			},
		},
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: nil,
		Handler:     be.blockchainHelpHandler,
	}

	cmdBlockchain := Command{
		Name:        BlockChainCommandName,
		Desc:        "Blockchain information and tools",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&subCmdCalcReward, &subCmdHelp},
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

func (be *BotEngine) blockchainHelpHandler(source AppID, callerID string, args ...string) (*CommandResult, error) {
	if len(args) == 0 {
		return be.help(source, callerID, BlockChainCommandName)
	}
	return be.help(source, callerID, BlockChainCommandName, args[0])
}
