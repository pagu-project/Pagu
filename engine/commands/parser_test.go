package commands_test

import (
	"testing"

	"github.com/kehiy/RoboPac/engine/commands"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	query := "claim tpc1pg663t0fyae0a8kvpg7mfj5nptjrgmjq930hz2q 1700842078"
	paredQuery := commands.ParseQuery(query)

	assert.Equal(t, paredQuery.Cmd, "claim")
	assert.Equal(t, paredQuery.Tokens[0], "tpc1pg663t0fyae0a8kvpg7mfj5nptjrgmjq930hz2q")
	assert.Equal(t, paredQuery.Tokens[1], "1700842078")
}
