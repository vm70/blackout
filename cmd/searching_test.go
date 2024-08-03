package cmd

import (
	"os"
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

func TestSearchingIsDeterministic(t *testing.T) {
	regexpString := msg2regex("a very long message")
	blackoutRegex := regexp.MustCompile(regexpString)
	setupErr := setupDataFolder()
	if setupErr != nil {
		t.Fatalf(setupErr.Error())
	}
	dir, dirErr := os.ReadDir("testdata/poems_folder")
	if dirErr != nil {
		t.Fatalf(dirErr.Error())
	}
	for nThreads := 1; nThreads < 10; nThreads++ {
		sp := SearchParams{dataFolderPoems, len(dir), nThreads, blackoutRegex, MaxLength, Profanities}
		poemID, searchErr := searchPoemsFolder(sp)
		if searchErr != nil {
			t.Fatalf(searchErr.Error())
		}
		for i := 0; i < 10; i++ {
			loopPoemID, searchErr := searchPoemsFolder(sp)
			if searchErr != nil {
				t.Fatalf(searchErr.Error())
			}
			if loopPoemID != poemID {
				t.Fatalf("%d != %d", loopPoemID, poemID)
			}
		}
	}
}
