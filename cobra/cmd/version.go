package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version string
var date string
var hash string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print artui version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Artui version: %s, build_date: %s, commit_hash: %s\n", version, date, hash)
		os.Exit(0)
	},
}
