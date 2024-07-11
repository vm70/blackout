package main

import "testing"

func TestPoemRoundTrip(t *testing.T) {
	test_poem := Poem{"Test Poem", "Somebody", "something", 9}
	test_poem_filename := "test_poem.json"
	err := poem2json(test_poem, test_poem_filename)
	if err != nil {
		t.Fatalf("Could not read to %s: %e", test_poem_filename, err)
	}
	file_poem, err := json2poem(test_poem_filename)
	if file_poem != test_poem {
		t.Fatalf("Poems do not match")
	}
}
