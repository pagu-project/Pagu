package main

import (
	robopac "github.com/robopac-project/RoboPac"
	"github.com/robopac-project/RoboPac/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "robopac-telegram",
		Version: robopac.StringVersion(),
	}

	RunCommand(rootCmd)

	err := rootCmd.Execute()
	cmd.ExitOnError(rootCmd, err)
}
