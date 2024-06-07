package omitlint

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// In order to test that the assertions we make about the "Good" type in
// internal/testdata hold true, we have to actually marshal the struct; however
// if we use the normal analysistest strategy of creating a gopath-style tree,
// we have no access to import the code contained therein from tests in the root
// directory. Luckily, analysistest also supports testing modules, so we can
// turn testdata into a module and `go run` a main package inside the module.
// This allows us access to the testdata type so that it can be encoded and
// checked.
func TestOmitemptyAssertion(t *testing.T) {
	// these two fields are the only two that should encode from a zero value of
	// the Good type. The first does not have the omitempty option set, while
	// the second is actually named "omitempty".
	const want = `{"notOmitempty":{},"omitempty":{}}`

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal("couldn't get working directory", err)
	}
	if err := os.Chdir(filepath.Join(pwd, "internal", "testdata")); err != nil {
		t.Fatal("couldn't chdir to internal/testdata")
	}
	got, err := exec.Command("go", "run", "assert_omitempty.go").CombinedOutput()
	if err != nil {
		t.Fatal("unable to run testassertion.go", err)
	}
	if string(got) != want {
		t.Fatalf("omit assertion failed\n\twant: %s\n\tgot: %s", want, got)
	}
}
