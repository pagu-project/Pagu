package market

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
)

func (m *Market) getPrice(_ *entity.User, cmd *command.Command, _ map[string]string) command.CommandResult {
	priceData, ok := m.priceCache.Get(config.PriceCacheKey)
	if !ok {
		return cmd.ErrorResult(fmt.Errorf("failed to get price from markets. please try again later"))
	}

	sb := strings.Builder{}
	xeggexPrice, err := strconv.ParseFloat(priceData.XeggexPacToUSDT.LastPrice, 64)
	if err == nil {
		sb.WriteString(fmt.Sprintf("Xeggex Price: %f	USDT\n https://xeggex.com/market/PACTUS_USDT \n\n",
			xeggexPrice))
	}

	if priceData.AzbitPacToUSDT.Price > 0 {
		sb.WriteString(fmt.Sprintf("Azbit Price: %f	USDT\n https://azbit.com/exchange/PAC_USDT \n\n",
			priceData.AzbitPacToUSDT.Price))
	}

	return cmd.SuccessfulResult(sb.String()) //nolint
}
