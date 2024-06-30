package cmd

import (
	"docute/generator/config"
	"github.com/spf13/cobra"
)

func InitCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "init",
		Short: "init a brand new doc",
		RunE: func(cmd *cobra.Command, args []string) error {

			target, _ := cmd.Flags().GetString("target")

			return config.InitConfig(target)
		},
	}

	return &c
}
