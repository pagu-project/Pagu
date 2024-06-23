package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pagu-project/Pagu/internal/entity"
)

func commandTestSetup() Command {
	return Command{
		Name: "Help",
		Desc: "",
		Help: "",
		Args: []Args{},
		SubCommands: []Command{
			{
				Name: "A",
				Desc: "some description for command A",
				Help: "",
			},
			{
				Name: "B",
				Desc: "some description for command B",
				Help: "",
			},
			{
				Name: "C",
				Desc: "some description for command C",
				Help: "",
			},
		},
		Middlewares: nil,
		AppIDs:      entity.AllAppIDs(),
		User:        &entity.User{ID: 1},
	}
}

func TestCommand_HelpMessage(t *testing.T) {
	command := commandTestSetup()
	message := command.HelpMessage()
	assert.Equal(t, "\n\nAvailable commands:\n<table><tr><td>A</td><td>some description for command A</td></tr><tr><td>B</td><td>some description for command B</td></tr><tr><td>C</td><td>some description for command C</td></tr></table>", message)
}
