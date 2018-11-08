package commands

import (
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "xxx",
	Short: "yyy",
	Long:  `zzz`,
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "config.json", "Configuration file to use.")
}
