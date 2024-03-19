package main

import (
	"bufio"
	"os"
	"strings"

	robopac "github.com/kehiy/RoboPac"
	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/engine/command"
	"github.com/kehiy/RoboPac/log"
	"github.com/pactus-project/pactus/crypto"
	cobra "github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Println("Provide your Discord ID as the first argument.")
		cmd.Println("Usage: robopac-cli <Discord-ID>")
		return
	}
	log.InitGlobalLogger()

	envOpt := cmd.Flags().StringP("env", "e", ".env", "the env file path")
	config, err := config.Load(*envOpt)
	ExitOnError(cmd, err)

	if config.Network == "Localnet" {
		crypto.AddressHRP = "tpc"
	}

	botEngine, err := engine.NewBotEngine(config)
	ExitOnError(cmd, err)

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

		callerID := args[0]
		inputs := strings.Split(input, " ")

		response := botEngine.Run(command.AppIdCLI, callerID, inputs)

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
	ExitOnError(rootCmd, err)
}

func ExitOnError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErr(err.Error())
		os.Exit(1)
	}
}
