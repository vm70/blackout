package main

import (
	"testing"
)

func TestDownloadingPoems(t *testing.T) {
	downloadPoems("poems.json")
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
