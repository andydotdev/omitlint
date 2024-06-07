package omitlint_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"andy.dev/omitlint"
)

func TestOmitLint(t *testing.T) {
	testdata, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("couldn't find testdata")
	}
	// testdata = analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, omitlint.NewAnalyzer(), "testdata/a")
}
