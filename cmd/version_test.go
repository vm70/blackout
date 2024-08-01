package cmd

import (
	"regexp"
	"testing"
)

// Semantic versioning regex, including the `v` in Go's variant.
var semverRegex = regexp.MustCompile(`^v(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func TestVersionCompliance(t *testing.T) {
	if !semverRegex.Match([]byte(Version)) {
		t.Fail()
	}
}
