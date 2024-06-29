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
	targetMask int
}

func NewMarket(clientMgr *client.Mgr, priceCache cache.Cache[string, entity.Price], target int) Market {
	return Market{
		clientMgr:  clientMgr,
		priceCache: priceCache,
		targetMask: target,
	}
}

func (m *Market) GetCommand() command.Command {
	subCmdPrice := command.Command{
		Name:        PriceCommandName,
		Desc:        "Shows the last price of PAC coin on the markets",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     m.getPrice,
	}

	cmdMarket := command.Command{
		Name:        CommandName,
		Desc:        "Blockchain data and information",
		Help:        "",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]command.Command, 0),
		Handler:     nil,
		TargetMask:  m.targetMask,
	}

	cmdMarket.AddSubCommand(subCmdPrice)

	return cmdMarket
}
