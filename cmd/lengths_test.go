package cmd

import (
	"path/filepath"
	"testing"
)

func TestGetLengths(t *testing.T) {
	downloadPoems("poems.json")
	poems, _ := readPoemsJSON("poems.json")
	splitPoems(poems, "poem_folder")
	lengths, lengthsErr := getLengths("poem_folder")
	if lengthsErr != nil {
		t.Fail()
	}
	for idx, length := range lengths {
		poem, poemErr := json2poem(filepath.Join("poem_folder", poemFilename(idx)))
		if poemErr != nil {
			t.Fail()
		}
		if len(poem.Text) != length {
			t.Fail()
		}
	}
}
