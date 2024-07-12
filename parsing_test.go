package main

import (
	"runtime"
	"testing"
)

func TestDownloadingPoems(t *testing.T) {
	downloadPoems("poems.json")
}

func TestPoemsLocation(t *testing.T) {
	if (runtime.GOOS[0:4] == "linux") && (dataFolder != "~/.local/share/blackout/") {
		t.Fail()
	}
}

func TestReadPoemDB(t *testing.T) {
	downloadPoems("poems.json")
	poems, err := readPoemDB("poems.json")
	if err != nil {
		t.Fail()
	}
	err = splitPoems(poems, "poem_folder")
	if err != nil {
		t.Fail()
	}
}
