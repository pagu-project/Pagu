package main

import (
	"os"

	"github.com/robopac-project/RoboPac/log"
	robopac "github.com/robopac-project/RoboPac"
	"github.com/spf13/cobra"
)

func main() {
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
