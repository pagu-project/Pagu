package blockchain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pactus-project/pactus/util"
	"github.com/robopac-project/RoboPac/client"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/utils"
)

const (
	BlockChainCommandName     = "blockchain"
	CalcRewardCommandName     = "calc-reward"
	CalcFeeCommandName        = "calc-fee"
	BlockChainHelpCommandName = "help"
)

type Blockchain struct {
	ctx       context.Context
	clientMgr *client.Mgr
}

func NewBlockchain(ctx context.Context,
	clientMgr *client.Mgr,
) Blockchain {
	return Blockchain{
		ctx:       ctx,
		clientMgr: clientMgr,
	}
}

func (bc *Blockchain) GetCommand() command.Command {
	subCmdCalcReward := command.Command{
		Name: CalcRewardCommandName,
		Desc: "Calculate how many PAC coins you will earn with your validator stake",
		Help: "Provide a stake amount between 1 to 100, please avoid using float numbers like: 1.9 or PAC suffix",
		Args: []command.Args{
			{
				Name:     "stake-amount",
				Desc:     "Amount of stake in your validator (1-1000)",
				Optional: false,
			},
			{
				Name:     "time-interval",
				Desc:     "After one: day | month | year",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord, command.AppIdgRPC},
		Handler:     bc.calcRewardHandler,
	}

	subCmdCalcFee := command.Command{
		Name: CalcFeeCommandName,
		Desc: "Calculate fee of a transaction with providing amount",
		Help: "Provide your amount in PAC, please avoid using float numbers like: 1.9 or PAC suffix",
		Args: []command.Args{
			{
				Name:     "amount",
				Desc:     "Amount of transaction",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord, command.AppIdgRPC},
		Handler:     bc.calcFeeHandler,
	}

	cmdBlockchain := command.Command{
		Name:        BlockChainCommandName,
		Desc:        "Blockchain information and tools",
		Help:        "",
		Args:        nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord, command.AppIdgRPC},
		SubCommands: []command.Command{subCmdCalcReward, subCmdCalcFee},
		Handler:     nil,
	}

	cmdBlockchain.AddSubCommand(subCmdCalcReward)
	cmdBlockchain.AddSubCommand(subCmdCalcFee)

	cmdBlockchain.AddHelpSubCommand()

	return cmdBlockchain
}

func (bc *Blockchain) calcRewardHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	stake, err := strconv.Atoi(args[0])
	if err != nil {
		return cmd.ErrorResult(err)
	}

	time := args[1]

	if stake < 1 || stake > 1_000 {
		return cmd.ErrorResult(fmt.Errorf("%v is invalid amount; minimum stake amount is 1 PAC and maximum is 1,000 PAC", stake))
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

	bi, err := bc.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	reward := int64(stake*blocks) / int64(util.ChangeToCoin(bi.TotalPower))

	return cmd.SuccessfulResult("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator in one %s ⏰ with %v PAC total power ⚡ of committee."+
		"\n\n> Note📝: This number is just an estimation. It will vary depending on your stake amount and total network power.",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), time, utils.FormatNumber(bi.TotalPower))
}

func (bc *Blockchain) calcFeeHandler(cmd command.Command, _ command.AppID, _ string, args ...string) command.CommandResult {
	amt, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	fee, err := bc.clientMgr.GetFee(util.CoinToChange(float64(amt)))
	if err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResult("Sending %v PAC will cost %v PAC with current fee percentage."+
		"\n> Note: Consider unbond and sortition transaction fee is 0 PAC always.", amt, util.ChangeToString(fee))
}
