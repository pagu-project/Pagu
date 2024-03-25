package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/pactus-project/pactus/crypto"
	robopac "github.com/robopac-project/RoboPac"
	"github.com/robopac-project/RoboPac/config"
	"github.com/robopac-project/RoboPac/engine"
	"github.com/robopac-project/RoboPac/engine/command"
	"github.com/robopac-project/RoboPac/log"
	cobra "github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func run(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Println("Provide your platform ID as the first argument, eg Discord-ID or Telegram-Username")
		cmd.Println("Usage: robopac-cli <platform-ID>")
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

		// Determine the platform based on the callerID
		platform := determinePlatform(callerID)

		var response command.CommandResult
		switch platform {
		case "telegram":
			response = botEngine.Run(command.AppIdCLI, callerID, inputs)
		case "discord":
			response = botEngine.Run(command.AppIdCLI, callerID, inputs)
		default:
			cmd.Println("Unsupported platform")
			return
		}

		cmd.Printf("%v\n%v", response.Title, response.Message)
	}
}

func determinePlatform(callerID string) string {
	// if it starts with "@" then its telegram because telegram usernames start with "@".
	if strings.HasPrefix(callerID, "@") {
		return "telegram"
	}
	//return discord as default .
	return "discord"
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
