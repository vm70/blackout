package cmd

import "testing"

func TestIsProfane(t *testing.T) {
	profanePoem, prErr := json2poem("testdata/profane_poem.json")
	if prErr != nil {
		t.Fatalf(prErr.Error())
	}
	nonProfanePoem, nonprErr := json2poem("testdata/non_profane_poem.json")
	if nonprErr != nil {
		t.Fatalf(nonprErr.Error())
	}
	if isProfane(nonProfanePoem) {
		t.Fatalf("non-profane poem should not be profane")
	}
	if !isProfane(profanePoem) {
		t.Fatalf("profane poem should not be profane")
	}
}

func TestPoemRoundTrip(t *testing.T) {
	testPoem := Poem{"Lorem", "Ipsum", "Dolor Sit Amet"}
	testPoemFilename := "testdata/test_poem.json"
	jsonErr := poem2json(testPoem, testPoemFilename)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}
	filePoem, poemErr := json2poem(testPoemFilename)
	if poemErr != nil {
		t.Fatal(poemErr)
	}
	if filePoem != testPoem {
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
