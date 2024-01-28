package engine

import (
	"fmt"
	"strings"
)

const (
	botCmdClaim = "claim"
)

// The input is always string.
//
//	The input format is like: [Command] <Arguments ...>
//
// The output is always string, but format might be JSON. ???
func (be *BotEngine) Run(input string) (string, error) {
	cmd, args := be.parseQuery(input)

	switch cmd {
	case botCmdClaim:
		if len(args) != 3 {
			return "", fmt.Errorf("expected to have 3 arguments, but it received %d", len(args))
		}

		_, err := be.Claim(args[0], args[1], args[2])
		if err != nil {
			return "", err
		}
		return "", nil

	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

func (be *BotEngine) parseQuery(query string) (string, []string) {
	subs := strings.Split(query, " ")
	if len(subs) == 0 {
		return "", nil
	}

	return subs[0], subs[1:]
}
