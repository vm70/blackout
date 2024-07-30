package cmd

import (
	"path/filepath"
	"testing"
)

func TestGetLengths(t *testing.T) {
	downloadErr := downloadPoemsJSON("testdata/poems.json")
	if downloadErr != nil {
		t.Fatal(downloadErr)
	}
	poems, readErr := readPoemsJSON("testdata/poems.json")
	if readErr != nil {
		t.Fatal(readErr)
	}
	splitErr := splitPoems(poems, "testdata/poems_folder")
	if splitErr != nil {
		t.Fatal(splitErr)
	}
	lengths, lengthsErr := getLengths("testdata/poems_folder")
	if lengthsErr != nil {
		t.Fatal(lengthsErr)
	}
	for idx, length := range lengths {
		poem, poemErr := json2poem(filepath.Join("testdata/poems_folder", poemFilename(idx)))
		if poemErr != nil {
			t.Fatal(poemErr)
		}
		if len(poem.Text) != length {
			t.Fatal("Poem text length does not match recorded length")
		}
	}
}
