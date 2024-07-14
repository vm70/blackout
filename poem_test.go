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

func TestMessage2Regex(t *testing.T) {
	testMessage := "blackoutpoem"
	testRegex := `(?m)^(.*?)(b)(.*?)(l)(.*?)(a)(.*?)(c)(.*?)(k)(.*?)(o)(.*?)(u)(.*?)(t)(.*?)(p)(.*?)(o)(.*?)(e)(.*?)(m)(.*?)$`
	if msg2regex(testMessage) != testRegex {
		t.Fatalf("Regexes don't match: %s, %s", testRegex, msg2regex(testMessage))
	}
}
