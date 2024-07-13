package main

import (
	"regexp"
	"testing"
)

func TestCanBlackout(t *testing.T) {
	testPoem := Poem{"Test Poem", "Somebody", "something"}
	testPoemLength := 9
	goodRegexP, _ := regexp.Compile("thing")
	badRegexP, _ := regexp.Compile("bad")
	// Check good regex with good length -> false
	if canBlackout(goodRegexP, testPoem, testPoemLength, 999) != true {
		t.Fail()
	}
	// Check good regex with bad length -> false
	if canBlackout(goodRegexP, testPoem, testPoemLength, 1) != false {
		t.Fail()
	}
	// Check bad regex with good length -> false
	if canBlackout(badRegexP, testPoem, testPoemLength, 999) != false {
		t.Fail()
	}
	// Check bad regex with bad length -> false
	if canBlackout(badRegexP, testPoem, testPoemLength, 1) != false {
		t.Fail()
	}
}

// func TestSearchingRoutine(t *testing.T) {
//   poems :=
//   testRegexP := regexp.Compile("(.*?)(a)(.*?)(b)(.*?)(c)(.*?)")
//   searchingRoutine(0, 1, testRegexP, poems, )
// }
