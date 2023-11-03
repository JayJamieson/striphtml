package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.striphtml.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "striphtml",
	Short: "Remove html tags and output plain text",
	Long: `Run striphtml as a standalone command to remove html tags or

run as a HTTP server send HTML documents directly to to be stripped:

Run on port 8080
striphtml serve -p 8080

Strip html from a provided url
curl -X GET http://localhost:8080/strip -d 'url=https://www.google.com'

Send html directly and get plain text as a response
curl -X POST -H 'Content-Type: text/html' http://localhost:8080/strip -d '<div>Hello world!</div>'
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func main() {

	isttyc := isatty.IsCygwinTerminal(os.Stdin.Fd())
	istty := isatty.IsTerminal(os.Stdin.Fd())

	fmt.Printf("IsCygwinTerminal %v", isttyc)
	fmt.Printf("IsTerminal %v", istty)
}
