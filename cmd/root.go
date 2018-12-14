package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/giantswarm/lighthouse-keeper/cmd/compare"
	"github.com/giantswarm/lighthouse-keeper/cmd/view"
)

var RootCmd = &cobra.Command{
	Use:   "lighthouse-keeper",
	Short: "Creates and pretty-prints lighthouse reports for a CI context",
}

func init() {
	RootCmd.AddCommand(compare.Cmd)
	RootCmd.AddCommand(view.Cmd)
}

// Execute is called by main to run the CLI
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
