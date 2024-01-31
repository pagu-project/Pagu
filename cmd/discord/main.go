package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "robopac-discord",
		Version: "0.0.1",
	}

	RunCommand(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		kill(rootCmd, err)
	}
}

func kill(cmd *cobra.Command, err error) {
	cmd.PrintErr(err.Error())
	os.Exit(1)
}
