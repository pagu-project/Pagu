package commands

import (
	"bufio"
	"os"
	"strings"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/engine/commands"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
	cobra "github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func REPLCommand(parentCmd *cobra.Command) {
	connect := &cobra.Command{
		Use:   "repl",
		Short: "Runs a local RoboPac instance which is connects to a local-net Pactus node",
	}
	parentCmd.AddCommand(connect)

	localnetOpt := connect.Flags().StringP("localnet", "l", "localhost:8080", "your local-net node address")
	envOpt := connect.Flags().StringP("env", "e", ".env.local", "your local/test env file for config")

	connect.Run = func(cmd *cobra.Command, args []string) {
		log.Info("initializing repl...")

		config, err := config.Load(*envOpt)
		if err != nil {
			log.Panic("can't load config env", "err", err, "path", *envOpt)
		}

		cm := client.NewClientMgr()
		c, err := client.NewClient(*localnetOpt)
		if err != nil {
			log.Panic("can't make a new local-net client", "err", err, "addr", *localnetOpt)
		}

		cm.AddClient("local-net", c)

		// initializing logger global instance.
		log.InitGlobalLogger()

		// new subLogger for engine.
		eSl := log.NewSubLogger("engine")

		// new subLogger for store.
		sSl := log.NewSubLogger("store")

		// new subLogger for store.
		wSl := log.NewSubLogger("wallet")

		// load or create wallet.
		wallet := wallet.Open(config, wSl)
		if wallet == nil {
			log.Panic("wallet could not be opened, wallet is nil", "path", config.WalletPath)
		}

		log.Info("wallet opened successfully", "address", wallet.Address())

		// load store.
		store, err := store.LoadStore(config, sSl)
		if err != nil {
			log.Panic("could not load store", "err", err, "path", config.StorePath)
		}

		log.Info("store loaded successfully", "path", config.StorePath)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(eSl, cm, wallet, store)
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

			q := commands.ParseQuery(input)
			response, err := commands.Execute(q, botEngine)
			if err != nil {
				cmd.PrintErr(err)
			}

			cmd.Print(response)
		}
	}
}
