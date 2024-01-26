package commands

import (
	"errors"

	"github.com/kehiy/RoboPac/engine"
)

func Execute(q Query, be engine.Engine) (interface{}, error) {
	switch q.Cmd {
	case "health":
		return be.NetworkHealth(q.Tokens)
	case "node-info":
		return be.NodeInfo(q.Tokens)
	case "claim":
		return be.Claim(q.Tokens)
	case "me":
		return be.ClaimerInfo(q.Tokens)
	case "network":
		return be.NetworkStatus(q.Tokens)
	default:
		return nil, errors.New("invalid command")
	}
}
