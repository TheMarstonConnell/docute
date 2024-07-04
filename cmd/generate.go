package cmd

import (
	"docute/gen"
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

			gen.Gen(target, "out")

			return nil
		},
	}

	return &c
}

func HostCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "host",
		Short: "host the static site.",
		RunE: func(cmd *cobra.Command, args []string) error {

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

func WatchCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "watch",
		Short: "Watch for changes and host them as a static test site.",
		RunE: func(cmd *cobra.Command, args []string) error {
			target, _ := cmd.Flags().GetString("target")

			gen.Gen(target, "out")

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
