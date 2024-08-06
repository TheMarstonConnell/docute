package cmd

import (
	"fmt"
	"os"
)

import "github.com/spf13/cobra"

func RootCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "docute",
		Short: "Docute is an insanely fast static documentation site generator.",
		Long:  `Docute is an insanely fast static documentation site generator.`,
	}
	c.PersistentFlags().StringP("base", "b", "/", "choose which base to use for all links")
	c.PersistentFlags().StringP("title", "t", "Docute Docs", "Set the title of your page")

	return &c
}

func Execute(c *cobra.Command) {
	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
