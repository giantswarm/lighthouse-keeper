package parser

import (
	"io/ioutil"
	"strings"
	"testing"
)

// TestParseJSON reads all JSON files form the testdata directory
// and attempts to parse them.
func TestParseJSON(t *testing.T) {
	basePath := "testdata"
	files, _ := ioutil.ReadDir(basePath)
	for _, f := range files {
		path := basePath + "/" + f.Name()

		if !strings.HasSuffix(path, ".json") {
			continue
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Error(err)
		}

		report, err := ParseReportJSON(data)
		if err != nil {
			t.Error(err)
		}

		t.Logf("File %q: Lighthouse version %q", path, report.LighthouseVersion)
	}
}
