package main

import (
	"os"

	robopac "github.com/robopac-project/RoboPac"
	"github.com/robopac-project/RoboPac/log"
	"github.com/spf13/cobra"
)

func main() {
	log.SetLoggerLevel()

	rootCmd := &cobra.Command{
		Use:     "robopac-discord",
		Version: robopac.StringVersion(),
	}

	log.InitGlobalLogger()

	RunCommand(rootCmd)

	err := rootCmd.Execute()
	ExitOnError(rootCmd, err)
}

func ExitOnError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErr(err.Error())
		os.Exit(1)
	}
}
