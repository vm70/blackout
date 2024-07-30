package cmd

import (
	"regexp"
	"testing"
)

const (
	MaxInt = int(^uint(0) >> 1)
)

func TestCanBlackout(t *testing.T) {
	goodRegexP, _ := regexp.Compile("e")
	badRegexP, _ := regexp.Compile("xxxxxx")
	lengths, _ := getLengths("testdata/poems_folder")
	// Check good regex with good length -> false
	goodRgoodL, err := canBlackout(goodRegexP, "testdata/poems_folder", 0, lengths[0], MaxInt, true)
	if !(err == nil && goodRgoodL == true) {
		t.Logf("Error: %s", err.Error())
		t.Logf("Good regex with good length: %t", goodRgoodL)
		t.Fail()
	}
	// Check good regex with bad length -> false
	goodRbadL, err := canBlackout(goodRegexP, "testdata/poems_folder", 0, lengths[0], 0, true)
	if !(err == nil && goodRbadL == false) {
		t.Logf("Error: %s", err.Error())
		t.Logf("Good regex with bad length: %t", goodRbadL)
		t.Fail()
	}
	// Check bad regex with good length -> false
	badRgoodL, err := canBlackout(badRegexP, "testdata/poems_folder", 0, lengths[0], MaxInt, true)
	if !(err == nil && badRgoodL == false) {
		t.Logf("Error: %s", err.Error())
		t.Logf("Bad regex with good length: %t", badRgoodL)
		t.Fail()
	}
	// Check bad regex with bad length -> false
	badRbadL, err := canBlackout(badRegexP, "testdata/poems_folder", 0, lengths[0], 0, true)
	if !(err == nil && badRbadL == false) {
		t.Logf("Error: %s", err.Error())
		t.Logf("Bad regex with bad length: %t", badRbadL)
		t.Fail()
	}
}
