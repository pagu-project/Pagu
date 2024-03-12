package blockchain

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/engine/command"
	"github.com/kehiy/RoboPac/nowpayments"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/twitter_api"
	"github.com/kehiy/RoboPac/utils"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/pactus-project/pactus/util"
)

const (
	BlockChainCommandName     = "blockchain"
	CalcRewardCommandName     = "calc-reward"
	CalcFeeCommandName        = "calc-fee"
	BlockChainHelpCommandName = "help"
)

type Blockchain struct {
	sync.RWMutex //! remove this.

	ctx           context.Context
	AdminIDs      []string
	store         store.IStore
	wallet        wallet.IWallet
	nowpayments   nowpayments.INowpayment
	clientMgr     *client.Mgr
	twitterClient twitter_api.IClient
}

func NewBlockchain(ctx context.Context,
	adminIDs []string,
	store store.IStore,
	wallet wallet.IWallet,
	nowpayments nowpayments.INowpayment,
	clientMgr *client.Mgr,
	twitterClient twitter_api.IClient,
) *Blockchain {
	return &Blockchain{
		ctx:           ctx,
		AdminIDs:      adminIDs,
		store:         store,
		wallet:        wallet,
		nowpayments:   nowpayments,
		clientMgr:     clientMgr,
		twitterClient: twitterClient,
	}
}

func (bc *Blockchain) GetCommand() *command.Command {
	subCmdCalcReward := &command.Command{
		Name: CalcRewardCommandName,
		Desc: "Calculate how many PAC coins you will earn with your validator stake",
		Help: "Provide a stake amount between 1 to 100, please avoid using float numbers like: 1.9 or PAC prefix",
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
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler:     bc.calcRewardHandler,
	}

	cmdBlockchain := command.Command{
		Name:        BlockChainCommandName,
		Desc:        "Blockchain information and tools",
		Help:        "",
		Args:        nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		SubCommands: []*command.Command{subCmdCalcReward},
		Handler:     nil,
	}

	cmdBlockchain.AddSubCommand(subCmdCalcReward)

	return &cmdBlockchain
}

func (bc *Blockchain) calcRewardHandler(cmd *command.Command, _ command.AppID, _ string, args ...string) *command.CommandResult {
	stake, err := strconv.Atoi(args[0])
	if err != nil {
		return &command.CommandResult{
			Error:      err.Error(),
			Successful: false,
		}
	}

	time := args[1]

	if stake < 1 || stake > 1_000 {
		return &command.CommandResult{
			Error:      "minimum of stake is 1 PAC and maximum is 1,000 PAC",
			Successful: false,
		}
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
		return &command.CommandResult{
			Error:      err.Error(),
			Successful: false,
		}
	}

	reward := int64(stake*blocks) / int64(util.ChangeToCoin(bi.TotalPower))

	result := fmt.Sprintf("Approximately you earn %v PAC reward, with %v PAC stake 🔒 on your validator in one %s ⏰ with %v PAC total power ⚡ of committee."+
		"\n\n> Note📝: This number is just an estimation. It will vary depending on your stake amount and total network power.",
		utils.FormatNumber(reward), utils.FormatNumber(int64(stake)), time, utils.FormatNumber(bi.TotalPower))

	return &command.CommandResult{
		Successful: true,
		Message:    result,
	}
}
