// Package compare provides the `compare` command to perform a comparison
// of two lighthouse reports
package compare

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/giantswarm/microerror"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/giantswarm/lighthouse-keeper/service/parser"
)

// Cmd is our cobra command
var Cmd = &cobra.Command{
	Use:     "compare",
	Short:   "Compare two lighthouse reports",
	PreRunE: validateFlags,
	Run:     compare,
}

func init() {
	Cmd.Flags().StringArrayP("input", "i", []string{}, "Input file path, to be used twice")
	Cmd.Flags().StringArrayP("inputlabel", "l", []string{}, "Input file label, to b used twice")
}

func compare(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetStringArray("input")
	if err != nil {
		fmt.Println("Error while reading --input flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	inputLabel, err := cmd.Flags().GetStringArray("inputlabel")
	if err != nil {
		fmt.Println("Error while reading --inputlabel flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(inputLabel) == 0 {
		inputLabel = append(inputLabel, "A")
	}
	if len(inputLabel) == 1 {
		inputLabel = append(inputLabel, "B")
	}

	reports := []*parser.Report{}
	{
		for _, inputItem := range input {

			data, err := ioutil.ReadFile(inputItem)
			if err != nil {
				fmt.Printf("Error while reading file %q:\n", inputItem)
				fmt.Println(err)
				os.Exit(1)
			}

			report, err := parser.ParseReportJSON(data)
			if err != nil {
				fmt.Printf("Error while parsing report %q:\n", inputItem)
				fmt.Println(err)
				os.Exit(1)
			}

			reports = append(reports, report)
		}
	}

	// output table data
	data := [][]string{}

	// Compare main category scores
	for catID, catA := range reports[0].Categories {
		catB, ok := reports[1].Categories[catID]
		if !ok || catA.Score == catB.Score {
			continue
		}

		delta := fmt.Sprintf("%.0f", (catB.Score-catA.Score)*100)
		if string(delta[0]) == "-" {
			delta = color.RedString(delta)
		} else {
			delta = color.GreenString("+" + delta)
		}

		row := []string{
			catA.Title,
			fmt.Sprintf("%.0f", catA.Score*100),
			fmt.Sprintf("%.0f", catB.Score*100),
			delta,
		}

		data = append(data, row)

		// Compare individual audits
		for _, auditRef := range catA.AuditRefs {
			auditA, ok := reports[0].Audits[auditRef.ID]
			if !ok {
				continue
			}

			auditB, ok := reports[1].Audits[auditRef.ID]
			if !ok {
				continue
			}

			if auditA.Score == auditB.Score {
				continue
			}

			delta := fmt.Sprintf("%.0f", (auditB.Score-auditA.Score)*100)
			if string(delta[0]) == "-" {
				delta = color.RedString(delta)
			} else {
				delta = color.GreenString("+" + delta)
			}

			row := []string{
				"- " + auditA.Title,
				fmt.Sprintf("%.0f", auditA.Score*100),
				fmt.Sprintf("%.0f", auditB.Score*100),
				delta,
			}

			data = append(data, row)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	labels := []string{"", inputLabel[0], inputLabel[1], "Delta"}
	table.SetHeader(labels)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flag("input") == nil {
		return microerror.Maskf(invalidFlagsError, "please specify two reports to compare using --input/-i flags")
	}

	inputs, err := cmd.Flags().GetStringArray("input")
	if err != nil {
		return microerror.Maskf(invalidFlagsError, "could not read values for --input/-i flags")
	}
	if len(inputs) != 2 {
		return microerror.Maskf(invalidFlagsError, "please specify exactly two --input/-i flags")
	}

	return nil
}
