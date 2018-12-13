// Package parser implements the parsing of lighthouse reports into Go types
package parser

import (
	"encoding/json"

	"github.com/giantswarm/microerror"
)

// ParseReportJSON consumes lighthouse report JSON and returns a pointer to
// a Report object.
func ParseReportJSON(jsonBlob []byte) (*Report, error) {
	var report *Report

	err := json.Unmarshal(jsonBlob, &report)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return report, nil
}
