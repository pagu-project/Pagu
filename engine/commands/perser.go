package commands

import (
	"strings"
)

type Query struct {
	Cmd    string
	Tokens []string
}

func ParseQuery(query string) Query {
	command := ""
	args := []string{}

	for _, word := range strings.Split(query, " ") {
		if word == "" {
			continue
		}

		if command != "" {
			args = append(args, word)
		} else {
			command = word
		}
	}

	return Query{Cmd: command, Tokens: args}
}
