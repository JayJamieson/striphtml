package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/JayJamieson/striphtml"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var prettyTables bool
var omitLinks bool
var textOnly bool

func init() {
	rootCmd.Flags().BoolVarP(&prettyTables, "pretty-tables", "t", false, "Pretty print tables.")
	rootCmd.Flags().BoolVarP(&omitLinks, "omit-links", "l", false, "Omit links.")
	rootCmd.Flags().BoolVarP(&textOnly, "pretty", "p", false, "Print text in prettified format.")
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
$ curl -X GET http://localhost:8080/strip -d 'url=https://www.google.com'

Send html directly:
$ curl -X POST -H 'Content-Type: text/html' http://localhost:8080/strip -d '<div>Hello world!</div>'
`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
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

func main() {

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
