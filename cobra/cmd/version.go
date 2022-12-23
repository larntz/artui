package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version = "asdf"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print artui version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Artui version %s\n", version)
		os.Exit(0)
	},
}
