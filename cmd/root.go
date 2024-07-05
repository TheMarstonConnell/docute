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

	return &c
}

func Execute(c *cobra.Command) {
	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
