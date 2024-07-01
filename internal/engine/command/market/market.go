package market

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/cache"
	"github.com/pagu-project/Pagu/pkg/client"
)

const (
	CommandName      = "market"
	PriceCommandName = "price"
	HelpCommandName  = "help"
)

type Market struct {
	clientMgr  *client.Mgr
	priceCache cache.Cache[string, entity.Price]
}

func NewMarket(clientMgr *client.Mgr, priceCache cache.Cache[string, entity.Price]) Market {
	return Market{
		clientMgr:  clientMgr,
		priceCache: priceCache,
	}
}

func (m *Market) GetCommand() command.Command {
	subCmdPrice := command.Command{
		Name:        PriceCommandName,
		Help:        "Shows the last price of PAC coin on the markets",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     m.getPrice,
	}

	cmdMarket := command.Command{
		Name:        CommandName,
		Help:        "Blockchain data and information",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMain,
	}

	cmdMarket.AddSubCommand(subCmdPrice)

	return cmdMarket
}
