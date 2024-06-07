package omitlint

import "testing"

func TestVersionFlag(t *testing.T) {
	vf := versionFlag{}
	assert(t, vf.IsBoolFlag(), "version flag is boolean")
	assert(t, vf.Get() == nil, "Get() returns nil")
	assert(t, vf.String() == "", "String() returns empty string")
}

func assert(t *testing.T, v bool, desc string) {
	t.Helper()
	if !v {
		t.Errorf("assertion failed - %s", desc)
	}
}
