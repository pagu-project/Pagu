package calculator

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/client"
)

const (
	CommandName           = "calculate"
	CalcRewardCommandName = "reward"
	CalcFeeCommandName    = "fee"
	HelpCommandName       = "help"
)

type Calculator struct {
	clientMgr *client.Mgr
}

func NewCalculator(clientMgr *client.Mgr) Calculator {
	return Calculator{
		clientMgr: clientMgr,
	}
}

func (bc *Calculator) GetCommand() command.Command {
	subCmdCalcReward := command.Command{
		Name: CalcRewardCommandName,
		Help: "Calculate how many PAC coins you will earn with your validator stake",
		Args: []command.Args{
			{
				Name:     "stake",
				Desc:     "Amount of stake in your validator (1-1000)",
				Optional: false,
			},
			{
				Name:     "days",
				Desc:     "Number of days (1-365)",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     bc.calcRewardHandler,
	}

	subCmdCalcFee := command.Command{
		Name: CalcFeeCommandName,
		Help: "Calculate fee of a transaction with providing amount",
		Args: []command.Args{
			{
				Name:     "amount",
				Desc:     "Amount of transaction",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     bc.calcFeeHandler,
	}

	cmdBlockchain := command.Command{
		Name:        CommandName,
		Help:        "Calculator information and tools",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMain,
	}

	cmdBlockchain.AddSubCommand(subCmdCalcReward)
	cmdBlockchain.AddSubCommand(subCmdCalcFee)

	return cmdBlockchain
}
