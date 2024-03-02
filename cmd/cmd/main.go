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

func run(cmd *cobra.Command, args []string) {
	log.InitGlobalLogger()

	log.Info("initializing repl...")

	envOpt := cmd.Flags().StringP("env", "e", ".env", "the env file path")
	config, err := config.Load(*envOpt)
	if err != nil {
		log.Panic("can't load config env", "err", err, "path", *envOpt)
	}

	if config.Network == "Localnet" {
		crypto.AddressHRP = "tpc"
	}

	botEngine, err := engine.NewBotEngine(config)
	if err != nil {
		log.Panic("could not start discord bot", "err", err)
	}

	botEngine.RegisterCommands()

	botEngine.Start()

	log.Info("repl started")
	reader := bufio.NewReader(os.Stdin)

	for {
		cmd.Print(PROMPT)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if strings.ToLower(input) == "exit" {
			cmd.Println("exiting from repl")
		}

		callerID := args[0]
		inputs := strings.Split(input, " ")

		response, err := botEngine.Run(engine.AppIdCLI, callerID, inputs)
		if err != nil {
			cmd.PrintErr(err)
		}

		cmd.Print(response)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "robopac-cmd",
		Version: "0.0.1",
		Run:     run,
	}

	err := rootCmd.Execute()
	if err != nil {
		kill(rootCmd, err)
	}
}

func kill(cmd *cobra.Command, err error) {
	cmd.PrintErr(err.Error())
	os.Exit(1)
}
