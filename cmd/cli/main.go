package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/pkg/log"

	"github.com/pactus-project/pactus/crypto"
	pagu "github.com/pagu-project/Pagu"
	pCmd "github.com/pagu-project/Pagu/cmd"
	"github.com/pagu-project/Pagu/config"
	cobra "github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func run(cmd *cobra.Command, args []string) {
	envOpt := cmd.Flags().StringP("env", "e", ".env", "the env file path")
	config, err := config.Load(*envOpt)
	pCmd.ExitOnError(cmd, err)

	log.InitGlobalLogger(config.Logger)

	if config.Network == "Localnet" {
		crypto.AddressHRP = "tpc"
	}

	botEngine, err := engine.NewBotEngine(config)
	pCmd.ExitOnError(cmd, err)

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
		Use:     "pagu-cli",
		Version: pagu.StringVersion(),
		Run:     run,
	}

	err := rootCmd.Execute()
	pCmd.ExitOnError(rootCmd, err)
}
