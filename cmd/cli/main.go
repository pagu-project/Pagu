package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/pactus-project/pactus/crypto"
	robopac "github.com/robopac-project/RoboPac"
	rpCmd "github.com/robopac-project/RoboPac/cmd"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/log"
	cobra "github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func run(cmd *cobra.Command, args []string) {
	envOpt := cmd.Flags().StringP("env", "e", ".env", "the env file path")
	config, err := config.Load(*envOpt)
	rpCmd.ExitOnError(cmd, err)

	log.InitGlobalLogger(config.LoggerConfig)

	if config.Network == "Localnet" {
		crypto.AddressHRP = "tpc"
	}

	botEngine, err := engine.NewBotEngine(config)
	rpCmd.ExitOnError(cmd, err)

	botEngine.RegisterAllCommands()
	botEngine.Start()

	reader := bufio.NewReader(os.Stdin)

	for {
		cmd.Print(PROMPT)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if strings.ToLower(input) == "exit" {
			cmd.Println("exiting from repl")

			return
		}

		inputs := strings.Split(input, " ")

		response := botEngine.Run(command.AppIdCLI, "0", inputs)

		cmd.Printf("%v\n%v", response.Title, response.Message)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "robopac-cli",
		Version: robopac.StringVersion(),
		Run:     run,
	}

	err := rootCmd.Execute()
	rpCmd.ExitOnError(rootCmd, err)
}
