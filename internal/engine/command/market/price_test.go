package market

import (
	"testing"
	"time"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/job"
	"github.com/pagu-project/Pagu/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func setup() (Market, command.Command) {
	priceCache := cache.NewBasic[string, entity.Price](1 * time.Second)
	priceJob := job.NewPrice(priceCache)
	priceJobSched := job.NewScheduler()
	priceJobSched.Submit(priceJob)
	go priceJobSched.Run()
	m := NewMarket(nil, priceCache)

	return m, command.Command{
		Name:        PriceCommandName,
		Desc:        "Shows the last price of PAC coin on the markets",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      command.AllAppIDs(),
	}
}

func TestGetPrice(t *testing.T) {
	market, cmd := setup()
	time.Sleep(10 * time.Second)
	result := market.getPrice(cmd, command.AppIdDiscord, "")
	assert.Equal(t, result.Successful, true)
}
