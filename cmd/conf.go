package cmd

import (
	"github.com/spf13/cobra"
	"zft/ui"
)

var Confcmd = &cobra.Command{
	Use:   "conf",
	Short: "Set configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConf()
	},
}

func runConf() error {
	ui.Conf()
	return nil
}
