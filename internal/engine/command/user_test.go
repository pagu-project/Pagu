package command

import (
	"testing"

	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareHandler_CreateUser(t *testing.T) {
	db, _ := repository.NewDB("root:ns1294password@tcp(127.0.0.1:4417)/pagu?parseTime=true")
	middlewareHandler := NewMiddlewareHandler(db, nil)
	cmd := Command{}
	t.Run("success creation", func(t *testing.T) {
		err := middlewareHandler.CreateUser(cmd, entity.AppIdDiscord, "ABCD")
		assert.Equal(t, nil, err)
	})
}
