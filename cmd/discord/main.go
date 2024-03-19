package main

import (
	"os"

	robopac "github.com/kehiy/RoboPac"
	"github.com/kehiy/RoboPac/log"
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
