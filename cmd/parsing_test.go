package cmd

import (
	"testing"
)

func TestDownloadingPoems(t *testing.T) {
	err := downloadPoemsJSON("poems.json")
	if err != nil {
		t.Fail()
	}
}

func TestReadPoemDB(t *testing.T) {
	downloadPoemsJSON("poems.json")
	poems, err := readPoemsJSON("poems.json")
	if err != nil {
		t.Fail()
	}
	err = splitPoems(poems, "poem_folder")
	if err != nil {
		t.Fail()
	}
}
