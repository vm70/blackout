package cmd

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
	testMessage1 := "blackoutpoem"
	testMessage2 := "blackout poem"
	testRegex1 := `(?s)\A(.*?)(b)(.*?)(l)(.*?)(a)(.*?)(c)(.*?)(k)(.*?)(o)(.*?)(u)(.*?)(t)(.*?)(p)(.*?)(o)(.*?)(e)(.*?)(m)(.*?)\z`
	testRegex2 := `(?s)\A(.*?)(b)(.*?)(l)(.*?)(a)(.*?)(c)(.*?)(k)(.*?)(o)(.*?)(u)(.*?)(t)(.*?)(p)(.*?)(o)(.*?)(e)(.*?)(m)(.*?)\z`
	if msg2regex(testMessage1) != testRegex1 {
		t.Fatalf("Regexes don't match: %s, %s", testRegex1, msg2regex(testMessage1))
	}
	if msg2regex(testMessage2) != testRegex2 {
		t.Fatalf("Regexes don't match: %s, %s", testRegex2, msg2regex(testMessage2))
	}
}
