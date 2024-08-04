package cmd

import (
	"testing"
)

func TestDownloadingPoems(t *testing.T) {
	downloadErr := downloadPoemsJSON("testdata/poems.json")
	if downloadErr != nil {
		t.Fail()
	}
}

func TestReadPoemDB(t *testing.T) {
	downloadErr := downloadPoemsJSON("testdata/poems.json")
	if downloadErr != nil {
		t.Fail()
	}
	poems, readErr := readPoemsJSON("testdata/poems.json")
	if readErr != nil {
		t.Fail()
	}
	parseErr := parsePoems(poems, "testdata/poems_folder")
	if parseErr != nil {
		t.Fail()
	}
}
