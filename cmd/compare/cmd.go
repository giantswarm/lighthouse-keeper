// Package compare provides the `compare` command to perform a comparison
// of two lighthouse reports
package compare

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/giantswarm/microerror"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/giantswarm/lighthouse-keeper/service/commenter"
	"github.com/giantswarm/lighthouse-keeper/service/parser"
)

// Cmd is our cobra command
var Cmd = &cobra.Command{
	Use:     "compare",
	Short:   "Compare two lighthouse reports",
	PreRunE: validateFlags,
	Run:     compare,
	Example: `
  lighthouse-keeper compare \
    --input lighthouse-a.json --inputlabel before \
    --input lighthouse-b.json --inputlabel after

  lighthouse-keeper compare \
    --input lighthouse-a.json --inputlabel before \
    --input lighthouse-b.json --inputlabel after \
    --github-owner giantswarm \
    --github-repo lighthouse-keeper \
    --github-issue 11 \
    --github-token $(cat ~/.github-token) \
`,
}

func init() {
	Cmd.Flags().StringArrayP("input", "i", []string{}, "Input file path, to be used twice")
	Cmd.Flags().StringArrayP("inputlabel", "l", []string{}, "Input file label, to b used twice")
	Cmd.Flags().StringP("github-owner", "", "", "GitHub user or org owning the repo to post the result to as a comment")
	Cmd.Flags().StringP("github-repo", "", "", "GitHub repo to post the reult to as a comment")
	Cmd.Flags().IntP("github-issue", "", 0, "GitHub issue or PR ID to post this to as a comment")
	Cmd.Flags().StringP("github-token", "", "", "Personal GitHub auth token to submit the comparison as a comment")
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

	// table data that works in markdown, without ANSII escape sequences
	markdownData := [][]string{}

	// Compare main category scores
	for catID, catA := range reports[0].Categories {
		catB, ok := reports[1].Categories[catID]
		if !ok || catA.Score == catB.Score {
			continue
		}

		delta := fmt.Sprintf("%.0f", (catB.Score-catA.Score)*100)
		markdownDelta := delta

		if string(delta[0]) == "-" {
			delta = color.RedString(delta)
			markdownDelta = "❌  " + markdownDelta
		} else {
			delta = color.GreenString("+" + delta)
			markdownDelta = "✅  " + "+" + markdownDelta
		}

		row := []string{
			catA.Title,
			fmt.Sprintf("%.0f", catA.Score*100),
			fmt.Sprintf("%.0f", catB.Score*100),
			delta,
		}

		markdownRow := []string{
			"**" + catA.Title + "**",
			fmt.Sprintf("%.0f", catA.Score*100),
			fmt.Sprintf("%.0f", catB.Score*100),
			markdownDelta,
		}

		data = append(data, row)

		markdownData = append(markdownData, markdownRow)

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
			markdownDelta := delta

			if string(delta[0]) == "-" {
				delta = color.RedString(delta)
				markdownDelta = "❌  " + markdownDelta
			} else {
				delta = color.GreenString("+" + delta)
				markdownDelta = "✅  " + " +" + markdownDelta
			}

			row := []string{
				"- " + auditA.Title,
				fmt.Sprintf("%.0f", auditA.Score*100),
				fmt.Sprintf("%.0f", auditB.Score*100),
				delta,
			}

			markdownRow := []string{
				"- " + auditA.Title,
				fmt.Sprintf("%.0f", auditA.Score*100),
				fmt.Sprintf("%.0f", auditB.Score*100),
				markdownDelta,
			}

			data = append(data, row)

			markdownData = append(markdownData, markdownRow)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	labels := []string{"", inputLabel[0], inputLabel[1], "Delta"}
	table.SetHeader(labels)

	for _, v := range data {
		table.Append(v)
	}

	if len(data) > 0 {
		table.Render()
	} else if len(inputLabel) == 2 && inputLabel[0] != "" && inputLabel[1] != "" {
		fmt.Printf("The comparison of lighthouse reports between `%s` and `%s` showed no difference.\n", inputLabel[0], inputLabel[1])
	}

	// comment to Github
	var owner string
	var repo string
	var token string
	var issue int
	{
		owner, err = cmd.Flags().GetString("github-owner")
		if err != nil {
			fmt.Println("Error while reading --github-owner flag:")
			fmt.Println(err)
			os.Exit(1)
		}

		repo, err = cmd.Flags().GetString("github-repo")
		if err != nil {
			fmt.Println("Error while reading --github-repo flag:")
			fmt.Println(err)
			os.Exit(1)
		}

		token, err = cmd.Flags().GetString("github-token")
		if err != nil {
			fmt.Println("Error while reading --github-token flag:")
			fmt.Println(err)
			os.Exit(1)
		}

		issue, err = cmd.Flags().GetInt("github-issue")
		if err != nil {
			fmt.Println("Error while reading --github-issue flag:")
			fmt.Println(err)
			os.Exit(1)
		}

		if owner != "" && repo != "" && token != "" && issue != 0 {
			var body string
			if len(markdownData) > 0 {
				var buf bytes.Buffer
				markdownTable := tablewriter.NewWriter(&buf)
				markdownTable.SetHeader(labels)
				markdownTable.SetAutoWrapText(false)
				markdownTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
				markdownTable.SetCenterSeparator("|")
				markdownTable.AppendBulk(markdownData)
				markdownTable.Render()

				body = "Comparison of lighthouse reports:\n\n"
				body += buf.String()
			} else {
				body = fmt.Sprintf("Comparison of lighthouse reports between `%s` and `%s` showed no difference.", labels[0], labels[1])
			}

			err = commenter.AddComment(token, owner, repo, body, issue)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

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
