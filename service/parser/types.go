package parser

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
)

type ScoreDisplayMode int
type DisplayValue string

const (
	ScoreDisplayModeNotApplicable ScoreDisplayMode = 0
	ScoreDisplayModeBinary        ScoreDisplayMode = 1
	ScoreDisplayModeInformative   ScoreDisplayMode = 2
	ScoreDisplayModeManual        ScoreDisplayMode = 3
	ScoreDisplayModeNumeric       ScoreDisplayMode = 4
)

// Report represents the root structure of a lighthouse report
type Report struct {
	FetchTime         time.Time `json:"fetchTime"`
	LighthouseVersion string    `json:"lighthouseVersion"`
	FinalURL          string    `json:"finalUrl"`
	RequestedURL      string    `json:"requestedUrl"`

	Audits map[string]Audit `json:"audits"`

	Categories map[string]Category `json:"categories"`
}

// Audit represents an audit block a lighthouse report
type Audit struct {
	ID               string           `json:"id"`
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	Score            float32          `json:"score"`
	ScoreDisplayMode ScoreDisplayMode `json:"scoreDisplayMode"`
	RawValue         RawValue         `json:"rawValue"`
	DisplayValue     DisplayValue     `json:"displayValue"`
}

type Category struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Score     float32    `json:"score"`
	AuditRefs []AuditRef `json:"auditRefs"`
}

type AuditRef struct {
	ID     string `json:"id"`
	Weight int    `json:"weight"`
	Group  string `json:"group"`
}

// RawValue is the type representing
type RawValue interface{}

// UnmarshalJSON converts the various values of the displayValue field
// into a common string representation.
// Note: We are losing some iformation here. For example, from
//   ["%.1d% legible text", 95.06061966771442]
// we only return
//   "95.06"
func (dv *DisplayValue) UnmarshalJSON(b []byte) error {
	// if b is an array, try to parse into [formatString, number]
	if string(b[0]) == "[" {
		var arr = make([]interface{}, 2)
		err := json.Unmarshal(b, &arr)
		if err != nil {
			return microerror.Mask(err)
		}

		floatVal, floatOK := arr[1].(float64)
		_, stringOK := arr[0].(string)
		if floatOK && stringOK {
			*dv = DisplayValue(fmt.Sprintf("%.2f", floatVal))
		} else {
			*dv = DisplayValue("ERROR")
		}
	} else if string(b) == "" || string(b) == "\"\"" {
		dv = nil
	} else {
		*dv = DisplayValue(string(b))
	}

	return nil
}

func (sdm *ScoreDisplayMode) UnmarshalJSON(b []byte) error {
	str := string(b)

	switch str {
	case "\"manual\"":
		*sdm = ScoreDisplayModeManual
	case "\"numeric\"":
		*sdm = ScoreDisplayModeNumeric
	case "\"binary\"":
		*sdm = ScoreDisplayModeBinary
	case "\"informative\"":
		*sdm = ScoreDisplayModeInformative
	case "\"not-applicable\"":
		*sdm = ScoreDisplayModeNotApplicable
	}

	return nil
}
