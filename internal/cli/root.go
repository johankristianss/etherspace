package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const TimeLayout = "2006-01-02 15:04:05"

var ASCII bool
var Verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(versionCmd)
}

var rootCmd = &cobra.Command{
	Use:   "evrium",
	Short: "evrium",
	Long:  "CLI to interact with Etherspace DHT",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version",
	Long:  "Version",
	Run: func(cmd *cobra.Command, args []string) {
		ASCII = false
		ASCIIStr := os.Getenv("ETHERSPACE_CLI_ASCII")
		if ASCIIStr == "true" {
			ASCII = true
		}

		printVersionTable()
	},
}
