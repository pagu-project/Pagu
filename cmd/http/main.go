package main

import (
	pagu "github.com/pagu-project/Pagu"
	"github.com/pagu-project/Pagu/cmd"
	"github.com/spf13/cobra"
)

var configPath string

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu-http",
		Version: pagu.StringVersion(),
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yml", "config path ./config.yml")
	runCommand(rootCmd)
	err := rootCmd.Execute()
	cmd.ExitOnError(rootCmd, err)
}
