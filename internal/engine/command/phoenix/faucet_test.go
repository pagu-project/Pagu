package phoenix

import (
	"testing"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/stretchr/testify/assert"
)

func setup() (Phoenix, command.Command) {
	db, _ := repository.NewDB("root:ns1294password@tcp(127.0.0.1:4417)/pagu?parseTime=true")
	p := NewPhoenix(nil, 5, nil, *db)

	return p, command.Command{
		Name:        FaucetCommandName,
		Help:        "",
		Args:        []command.Args{},
		SubCommands: nil,
		Middlewares: nil,
		AppIDs:      entity.AllAppIDs(),
		User:        &entity.User{ID: 1},
	}
}

func TestPhoenixFaucet(t *testing.T) {
	t.Run("empty args", func(t *testing.T) {
		phoenix, cmd := setup()
		result := phoenix.faucetHandler(cmd, entity.AppIdDiscord, "")
		assert.Equal(t, "An error occurred: invalid wallet address", result.Message)
	})
}
