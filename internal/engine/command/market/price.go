package market

import (
	"fmt"
	"strconv"

	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (m *Market) getPrice(cmd *command.Command, _ entity.AppID, _ string, _ map[string]any) command.CommandResult {
	priceData, ok := m.priceCache.Get(config.PriceCacheKey)
	if !ok {
		return cmd.ErrorResult(fmt.Errorf("failed to get price from markets. please try again later"))
	}

	lastPrice, err := strconv.ParseFloat(priceData.XeggexPacToUSDT.LastPrice, 64)
	if err != nil {
		return cmd.ErrorResult(fmt.Errorf("pagu can not calculate the price. please try again later"))
	}

	return cmd.SuccessfulResult("PAC Price: %f	USDT"+
		"\n\n\n See below markets link for more details: \n xeggex: https://xeggex.com/market/PACTUS_USDT \n "+
		"exbitron: https://exbitron.com/trade?market=PAC-USDT", lastPrice)
}
