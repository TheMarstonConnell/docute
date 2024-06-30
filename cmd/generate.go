package cmd

import (
	"docute/generator"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func GenerateCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "generate",
		Short: "Generate a fully static site",
		Long:  `Generate the custom documentation site from your configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			target, _ := cmd.Flags().GetString("target")

			g := generator.NewGenerator(target)
			g.Start()
			g.PrintTree()
			g.Walk()
			return nil
		},
	}

	return &c
}

func WatchCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "watch",
		Short: "Watch for changes and host them as a static test site.",
		RunE: func(cmd *cobra.Command, args []string) error {
			target, _ := cmd.Flags().GetString("target")

			g := generator.NewGenerator(target)
			g.Start()
			g.PrintTree()
			g.Walk()

			fs := http.FileServer(http.Dir("./out"))

			// Serve static files from the /static URL path
			http.Handle("/", fs)

			// Start the server on port 9797
			log.Println("Listening on :9797...")
			err := http.ListenAndServe(":9797", nil)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	return &c
}
