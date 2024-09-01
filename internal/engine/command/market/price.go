package market

import (
	"fmt"
	"strconv"

	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (m *Market) getPrice(_ *entity.User, cmd *command.Command, _ map[string]string) command.CommandResult {
	priceData, ok := m.priceCache.Get(config.PriceCacheKey)
	if !ok {
		return cmd.ErrorResult(fmt.Errorf("failed to get price from markets. please try again later"))
	}

	xeggexPrice, xeggexErr := strconv.ParseFloat(priceData.XeggexPacToUSDT.LastPrice, 64)
	p2bPrice, p2bErr := strconv.ParseFloat(priceData.P2BPacToUSDT.LastPrice, 64)

	if xeggexErr != nil && p2bErr != nil {
		return cmd.ErrorResult(fmt.Errorf("pagu can not calculate the price. Please try agin later"))
	}

	return cmd.SuccessfulResult("Xeggex Price: %f	USDT\n P2B Price: %f	USDT"+
		"\n\n\n See below markets link for more details: \n xeggex: https://xeggex.com/market/PACTUS_USDT \n "+
		"exbitron: https://exbitron.com/trade?market=PAC-USDT", xeggexPrice, p2bPrice)
}
