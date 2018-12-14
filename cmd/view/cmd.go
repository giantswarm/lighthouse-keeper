// Package view provides the `view` command to print a lighthouse report
package view

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/giantswarm/microerror"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/giantswarm/lighthouse-keeper/service/parser"
)

// Cmd is our cobra command
var Cmd = &cobra.Command{
	Use:     "view",
	Short:   "Print a Lighthouse report",
	PreRunE: validateFlags,
	Run:     view,
}

func init() {
	Cmd.Flags().StringP("input", "i", "", "Input file path")
	Cmd.Flags().BoolP("omit-done", "o", false, "Avoid praising yourself, hide audit rows showing full score")
}

func view(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		fmt.Println("Error while reading --input flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	omitDone, err := cmd.Flags().GetBool("omit-done")
	if err != nil {
		fmt.Println("Error while reading --omit-done flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	var report *parser.Report
	{
		data, err := ioutil.ReadFile(input)
		if err != nil {
			fmt.Printf("Error while reading file %q:\n", input)
			fmt.Println(err)
			os.Exit(1)
		}

		report, err = parser.ParseReportJSON(data)
		if err != nil {
			fmt.Printf("Error while parsing report %q:\n", input)
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// output table data
	data := [][]string{}

	// Print by category, then by audit
	for _, cat := range report.Categories {
		row := []string{
			cat.Title,
			fmt.Sprintf("%.0f", cat.Score*100),
			"",
		}

		data = append(data, row)

		// individual audits
		for _, auditRef := range cat.AuditRefs {
			audit, ok := report.Audits[auditRef.ID]
			if !ok {
				continue
			}

			score := fmt.Sprintf("%.0f", audit.Score*100)

			if omitDone && score == "100" {
				continue
			}

			row := []string{
				"- " + audit.Title,
				score,
				fmt.Sprintf("%d", auditRef.Weight),
			}

			data = append(data, row)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	labels := []string{"Metric", "Score", "Weight"}
	table.SetHeader(labels)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flag("input") == nil {
		return microerror.Maskf(invalidFlagsError, "please specify a reports to compare using the --input/-i flag")
	}

	input, err := cmd.Flags().GetString("input")
	if err != nil {
		return microerror.Maskf(invalidFlagsError, "could not read value for --input/-i flag")
	}
	if input == "" {
		return microerror.Maskf(invalidFlagsError, "please specify a reports to compare using the --input/-i flag")
	}

	return nil
}
