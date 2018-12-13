package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/giantswarm/lighthouse-keeper/cmd/compare"
)

var RootCmd = &cobra.Command{
	Use:   "lighthouse-keeper",
	Short: "Creates and pretty-prints lighthouse reports for a CI context",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented")
	},
}

func init() {
	RootCmd.AddCommand(compare.Cmd)
}

// Execute is called by main to run the CLI
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
