package main

import "testing"

func TestPoemRoundTrip(t *testing.T) {
	testPoem := Poem{"Test Poem", "Somebody", "something"}
	testPoemFilename := "test_poem.json"
	err := poem2json(testPoem, testPoemFilename)
	if err != nil {
		t.Fatalf("Could not read to %s: %e", testPoemFilename, err)
	}
	filePoem, err := json2poem(testPoemFilename)
	if filePoem != testPoem {
		t.Fatalf("Poems do not match")
	}
}
