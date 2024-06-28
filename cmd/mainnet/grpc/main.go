package main

import (
	pagu "github.com/pagu-project/Pagu"
	"github.com/pagu-project/Pagu/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu-grpc",
		Version: pagu.StringVersion(),
	}

	runCommand(rootCmd)

	err := rootCmd.Execute()
	cmd.ExitOnError(rootCmd, err)
}
