package phoenix

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() (Phoenix, command.Command) {
	db, _ := repository.NewDB("root:ns1294password@tcp(127.0.0.1:4417)/pagu?parseTime=true")
	p := NewPhoenix(nil, 5, nil, *db)

	return p, command.Command{
		Name:        FaucetCommandName,
		Desc:        "",
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		Middlewares: nil,
		AppIDs:      entity.AllAppIDs(),
	}
}

func TestPhoenixFaucet(t *testing.T) {
	phoenix, cmd := setup()

	t.Run("empty args", func(t *testing.T) {
		result := phoenix.faucetHandler(cmd, entity.AppIdDiscord, "")
		assert.Equal(t, "An error occurred: invalid wallet address", result.Message)
	})
}
