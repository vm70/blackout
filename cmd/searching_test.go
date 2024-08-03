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

	goodR, err := canBlackout(goodRegexP, nonProfaneParsedPoem)
	if !(err == nil && goodR == true) {
		t.Logf("Error: %s", err.Error())
		t.Logf("Good regex with good length: %t", goodR)
		t.Fail()
	}
	// Check bad regex with good length -> false
	badR, err := canBlackout(badRegexP, nonProfaneParsedPoem)
	if !(err == nil && badR == false) {
		t.Logf("Error: %s", err.Error())
		t.Logf("Bad regex with good length: %t", badR)
		t.Fail()
	}
}
