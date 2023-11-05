package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JayJamieson/striphtml"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	prettyTables bool
	omitLinks    bool
	textOnly     bool
	port         int
	httpClient   *http.Client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)

	rootCmd.Flags().BoolVarP(&prettyTables, "pretty-tables", "t", false, "Pretty print tables.")
	rootCmd.Flags().BoolVarP(&omitLinks, "omit-links", "l", false, "Omit links.")
	rootCmd.Flags().BoolVar(&textOnly, "pretty", false, "Print text in prettified format.")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run HTTP server on")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "striphtml",
	Short: "Remove html tags and output plain text",
	Long: `This application is a tool to remove html tags and output plain text
representation of the file.

striphtml can be run as a standalone command or HTTP Server.

Standalone:
$ cat index.html | striphtml

$ striphtml < index.html

HTTP Server:
$ striphtml serve -p 8080

Strip html from a provided url:
$ curl -X GET http://localhost:8080/strip?url=https://www.google.com

Send html directly:
$ curl -X POST -H 'Content-Type: text/html' http://localhost:8080/strip -d '<div>Hello world!</div>'
`,
	Run: runCli,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run striphtml in server mode",
	Long: `striphtml can be run in server mode and used as an HTTP serrvice
to strip html tags from URL supplied query parameter or send a webpage directly`,
	Run: runServer,
}

func main() {

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func runCli(cmd *cobra.Command, args []string) {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		cmd.Help()
		return
	}

	reader := bufio.NewReader(os.Stdin)

	opts := striphtml.Options{
		PrettyTables: prettyTables,
		OmitLinks:    omitLinks,
		TextOnly:     !textOnly,
	}

	out, err := striphtml.FromReader(reader, opts)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(out)
}

func handleStrip(w http.ResponseWriter, r io.ReadCloser) {
	reader := bufio.NewReader(r)

	opts := striphtml.Options{
		PrettyTables: prettyTables,
		OmitLinks:    omitLinks,
		TextOnly:     !textOnly,
	}

	out, err := striphtml.FromReader(reader, opts)

	w.Header().Add("Content-Type", "text/plain")

	if err != nil {
		msg := fmt.Sprintf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.Write([]byte(out))
}

func runServer(cmd *cobra.Command, args []string) {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/strip", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")

		switch url {
		case "":
			handleStrip(w, r.Body)
			return
		default:
			resp, err := httpClient.Get(url)

			if err != nil {
				msg := fmt.Sprintf("%v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(msg))
				return
			}

			handleStrip(w, resp.Body)
		}
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		if errShutdown := server.Shutdown(shutdownCtx); errShutdown != nil {
			fmt.Printf("error stopping: %v\n", errShutdown)
			os.Exit(1)
		}
	}()

	fmt.Printf("Starting server on http://localhost:%v", port)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error starting server: %v\n", err)
	}
}
