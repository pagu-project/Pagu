package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/pactus-project/pactus/crypto"
	cobra "github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func REPLCommand(parentCmd *cobra.Command) {
	connect := &cobra.Command{
		Use:   "repl",
		Short: "Runs a local RoboPac instance which is connects to a local-net Pactus node",
	}
	parentCmd.AddCommand(connect)

	envOpt := connect.Flags().StringP("env", "e", ".env", "the env file path")

	connect.Run = func(cmd *cobra.Command, args []string) {
		// initializing logger global instance.
		log.InitGlobalLogger()

		log.Info("initializing repl...")

		config, err := config.Load(*envOpt)
		if err != nil {
			log.Panic("can't load config env", "err", err, "path", *envOpt)
		}

		if config.Network == "Localnet" {
			crypto.AddressHRP = "tpc"
		}

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(config)
		if err != nil {
			log.Panic("could not start discord bot", "err", err)
		}
		botEngine.Start()

		log.Info("repl started")
		reader := bufio.NewReader(os.Stdin)

		for {
			cmd.Print(PROMPT)

			input, _ := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n")

			if strings.ToLower(input) == "exit" {
				cmd.Println("exiting from repl")

				return
			}

			response, err := botEngine.Run(input)
			if err != nil {
				cmd.PrintErr(err)
			}

			cmd.Print(response)
		}
	}
}
