package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func ExitOnError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErr(err.Error())
		os.Exit(1)
	}
}
