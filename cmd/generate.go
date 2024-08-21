package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JackalLabs/docute/gen"
	"github.com/gorilla/websocket"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func GenerateCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "generate",
		Short: "Generate a fully static site",
		Long:  `Generate the custom documentation site from your configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := cmd.Flags().GetString("base")
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString("title")
			if err != nil {
				return err
			}
			err = gen.Gen(".out", base, title, false)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return &c
}

func GenColorFile() *cobra.Command {
	c := cobra.Command{
		Use:   "colors",
		Short: "Generate a color file",
		RunE: func(cmd *cobra.Command, args []string) error {
			col := gen.DefaultColors()

			data, err := yaml.Marshal(col)
			if err != nil {
				return err
			}

			err = os.WriteFile("colors.yaml", data, os.ModePerm)
			if err != nil {
				return err
			}
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
			fs := http.FileServer(http.Dir(".out"))

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

func regen(cmd *cobra.Command, ws *websocket.Conn) error {
	base, err := cmd.Flags().GetString("base")
	if err != nil {
		return err
	}

	title, err := cmd.Flags().GetString("title")
	if err != nil {
		return err
	}
	err = gen.Gen(".out", base, title, true)
	if err != nil {
		return err
	}

	if ws != nil {
		_ = ws.WriteMessage(1, []byte("refresh"))
	}

	return nil
}

func WatchCMD() *cobra.Command {
	c := cobra.Command{
		Use:   "watch",
		Short: "Watch for changes and host them as a static test site.",
		RunE: func(cmd *cobra.Command, args []string) error {
			upgrader := websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					return true // Allow all connections by default
				},
			}
			var ws *websocket.Conn
			websocketHandler := func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					return
				}
				defer conn.Close()
				ws = conn

				for {
					_, _, err := conn.ReadMessage()
					if err != nil {
						break
					}
				}
			}

			http.HandleFunc("/ws", websocketHandler)
			fs := http.FileServer(http.Dir(".out"))
			http.Handle("/", fs)

			w := watcher.New()
			w.SetMaxEvents(1)

			go func() {
				for {
					select {
					case event := <-w.Event:
						fmt.Println(event)
						err := regen(cmd, ws)
						if err != nil {
							log.Fatalln(err)
						}
					case err := <-w.Error:
						log.Fatalln(err)
					case <-w.Closed:
						return
					}
				}
			}()

			if err := w.Ignore(".out"); err != nil {
				log.Fatalln(err)
			}

			w.IgnoreHiddenFiles(true)

			allFiles, err := os.ReadDir(".")
			if err != nil {
				log.Fatalln(err)
				return err
			}
			for _, file := range allFiles {
				// Watch this folder for changes.
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}
				if err := w.AddRecursive(file.Name()); err != nil {
					log.Fatalln(err)
				}
			}

			// Print a list of all of the files and folders currently
			// being watched and their paths.
			for path, f := range w.WatchedFiles() {
				fmt.Printf("%s: %s\n", path, f.Name())
			}

			go func() {
				// Start the watching process - it'll check for changes every 100ms.
				if err := w.Start(time.Millisecond * 100); err != nil {
					log.Fatalln(err)
				}
			}()

			err = regen(cmd, ws)
			if err != nil {
				log.Fatal(err)
			}

			// Start the server on port 9797
			log.Println("Listening on :9797...")
			err = http.ListenAndServe(":9797", nil)
			if err != nil {
				log.Fatal(err)
			}

			w.Close()

			return nil
		},
	}

	return &c
}
