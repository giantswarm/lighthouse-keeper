// Package audit provides the `audit` command to create a lighthouse report.
package audit

import (
	"fmt"
	"os"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"

	"github.com/giantswarm/lighthouse-keeper/service/lighthouse"
)

// Cmd is our cobra command
var Cmd = &cobra.Command{
	Use:     "audit",
	Short:   "Create a Lighthouse report",
	PreRunE: validateFlags,
	Run:     audit,
	Example: `
  lighthouse-keeper audit --url https://example.com/

  lighthouse-keeper audit --form-factor mobile --url https://example.com/

  lighthouse-keeper audit --name mysite --form-factor mobile --url https://example.com/

  lighthouse-keeper audit --url https://container:5000/ --docker-link container:container

  lighthouse-keeper audit \
    --name first-name --url http://first-url \
    --name second-name --url http://second-url`,
}

func init() {
	Cmd.Flags().StringArrayP("url", "u", []string{}, "URL to audit, can be used multiple times")
	Cmd.Flags().StringArrayP("name", "n", []string{}, "Output file name prefix, can be used multiple times")
	Cmd.Flags().StringP("form-factor", "f", "desktop", "Either 'desktop' or 'mobile")
	Cmd.Flags().StringArrayP("docker-link", "l", []string{}, "Link the lighthouse docker container to these named links")
}

func audit(cmd *cobra.Command, args []string) {
	urls, err := cmd.Flags().GetStringArray("url")
	if err != nil {
		fmt.Println("Error while reading --url flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	names, err := cmd.Flags().GetStringArray("name")
	if err != nil {
		fmt.Println("Error while reading --name flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	dockerLinks, err := cmd.Flags().GetStringArray("docker-link")
	if err != nil {
		fmt.Println("Error while reading --docker-link flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	formFactor, err := cmd.Flags().GetString("form-factor")
	if err != nil {
		fmt.Println("Error while reading --name flag:")
		fmt.Println(err)
		os.Exit(1)
	}

	for index, url := range urls {
		// set automatic output name if none given
		if len(names) < (index + 1) {
			t := time.Now()
			names = append(names, t.Format("20060102-150405")+fmt.Sprintf("-%s-%d", formFactor, index+1))
		}

		_, err := lighthouse.AuditURL(url, names[index], formFactor, dockerLinks)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

func validateFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flag("url") == nil {
		return microerror.Maskf(invalidFlagsError, "please specify at least one URL to audit using the --url/-u flag")
	}

	inputs, err := cmd.Flags().GetStringArray("url")
	if err != nil {
		return microerror.Maskf(invalidFlagsError, "could not read values for --url/-u flag")
	}
	if len(inputs) < 1 {
		return microerror.Maskf(invalidFlagsError, "please specify at least one URL to audit via the --url/-u flag")
	}

	return nil
}
