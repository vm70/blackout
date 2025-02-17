package cmd

import "testing"

var (
	// Example profane poem.
	profanePoem = Poem{"Lorem", "Ipsum", "Dolor Sit Amet Fuck"}
	// Example non-profane poem.
	nonProfanePoem = Poem{"Lorem", "Ipsum", "Dolor Sit Amet"}
	// Example non-profane parsed poem.
	nonProfaneParsedPoem = NewParsedPoem(nonProfanePoem)
)

func TestIsProfane(t *testing.T) {
	if isProfane(nonProfanePoem) {
		t.Fatalf("non-profane poem should not be profane")
	}
	if !isProfane(profanePoem) {
		t.Fatalf("profane poem should not be profane")
	}
}

func TestParsedPoemRoundTrip(t *testing.T) {
	testPoemFilename := "testdata/test_poem.json"
	jsonErr := parsedPoem2json(nonProfaneParsedPoem, testPoemFilename)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}
	filePoem, poemErr := json2parsedPoem(testPoemFilename)
	if poemErr != nil {
		t.Fatal(poemErr)
	}
	if filePoem != nonProfaneParsedPoem {
		t.Fatal("Original poem and file-read poem do not match")
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
