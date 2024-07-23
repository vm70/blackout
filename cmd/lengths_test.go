package cmd

import (
	"path/filepath"
	"testing"
)

func TestGetLengths(t *testing.T) {
	downloadErr := downloadPoemsJSON("poems.json")
	if downloadErr != nil {
		t.Fatal(downloadErr)
	}
	poems, readErr := readPoemsJSON("poems.json")
	if readErr != nil {
		t.Fatal(readErr)
	}
	splitErr := splitPoems(poems, "poem_folder")
	if splitErr != nil {
		t.Fatal(splitErr)
	}
	lengths, lengthsErr := getLengths("poem_folder")
	if lengthsErr != nil {
		t.Fatal(lengthsErr)
	}
	for idx, length := range lengths {
		poem, poemErr := json2poem(filepath.Join("poem_folder", poemFilename(idx)))
		if poemErr != nil {
			t.Fatal(poemErr)
		}
		if len(poem.Text) != length {
			t.Fatal("Poem text length does not match recorded length")
		}
	}
}
