package cmd

import (
	"regexp"
	"testing"
)

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func TestCanBlackout(t *testing.T) {
	goodRegexP, _ := regexp.Compile("e")
	badRegexP, _ := regexp.Compile("xxxxxx")
	lengths, _ := getLengths("poem_folder")
	// Check good regex with good length -> false
	goodRgoodL, err := canBlackout(goodRegexP, "poem_folder", 0, lengths[0], MaxInt)
	if !(err == nil && goodRgoodL == true) {
		t.Fail()
	}
	// Check good regex with bad length -> false
	goodRbadL, err := canBlackout(goodRegexP, "poem_folder", 0, lengths[0], 0)
	if !(err == nil && goodRbadL == false) {
		t.Fail()
	}
	// Check bad regex with good length -> false
	badRgoodL, err := canBlackout(badRegexP, "poem_folder", 0, lengths[0], MaxInt)
	if !(err == nil && badRgoodL == false) {
		t.Fail()
	}
	// Check bad regex with bad length -> false
	badRbadL, err := canBlackout(badRegexP, "poem_folder", 0, lengths[0], 0)
	if !(err == nil && badRbadL == false) {
		t.Fail()
	}
}
