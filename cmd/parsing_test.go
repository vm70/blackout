package cmd

import (
	"testing"
)

func TestDownloadingPoems(t *testing.T) {
	downloadErr := downloadPoemsJSON("poems.json")
	if downloadErr != nil {
		t.Fail()
	}
}

func TestReadPoemDB(t *testing.T) {
	downloadErr := downloadPoemsJSON("poems.json")
	if downloadErr != nil {
		t.Fail()
	}
	poems, readErr := readPoemsJSON("poems.json")
	if readErr != nil {
		t.Fail()
	}
	splitErr := splitPoems(poems, "poem_folder")
	if splitErr != nil {
		t.Fail()
	}
}
