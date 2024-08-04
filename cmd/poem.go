/*
Package cmd contains the necessary functions to execute the code for `blackout`.

Copyright © 2024 Vincent Mercator <vmercator@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/TwiN/go-away"
)

// regexEscapes contains the characters that need to be escaped in regexes.
const regexEscapes = `.+*?()|[]{}^$`

var (
	// blackoutRP is the regular expression pointer that matches every non-whitespace charater for blacking out.
	blackoutRP = regexp.MustCompile(`[^\t\f\r\n\ ]`)
	// Regular expression pointer that matches every non-ASCII character.
	ASCIIRP = regexp.MustCompile("[[:^ascii:]]")
)

// A Poem in the database has a title, an author, and text.
type Poem struct {
	Title  string // The title of the poem.
	Author string // The author of the poem.
	Text   string // The poem's text itself. Poem lines are delineated with the digraph "\n".
}

// A ParsedPoem has its length and level of profanity pre-computed.
type ParsedPoem struct {
	Title     string // The title of the poem.
	Author    string // The author of the poem.
	Text      string // The poem's text itself. Poem lines are delineated with the digraph "\n".
	Length    int    // The poem's length [in characters].
	IsProfane bool   // Whether the poem's text contains profane language.
}

// isProfane signals whether a poem's text contains profane language.
func isProfane(poem Poem) bool {
	filteredText := ASCIIRP.ReplaceAllLiteralString(poem.Text, "")
	return goaway.IsProfane(filteredText)
}

// NewParsedPoem creates a new parsed poem from a poem in the dataset.
func NewParsedPoem(poem Poem) ParsedPoem {
	length := len(poem.Text)
	isProfane := isProfane(poem)
	return ParsedPoem{poem.Title, poem.Author, poem.Text, length, isProfane}
}

// parsedPoem2json writes the ParsedPoem struct to the given JSON file path.
func parsedPoem2json(parsedPoem ParsedPoem, jsonFile string) error {
	// Marshal to JSON bytes
	poemBytes, err := json.Marshal(parsedPoem)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Write file
	err = os.WriteFile(jsonFile, poemBytes, 0o666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// json2parsedPoem extracts the ParsedPoem object from the given JSON file path.
func json2parsedPoem(jsonFile string) (ParsedPoem, error) {
	// Read the file name
	var parsedPoem ParsedPoem
	fileBytes, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
		return parsedPoem, err
	}
	// parse JSON and return poem
	err = json.Unmarshal(fileBytes, &parsedPoem)
	if err != nil {
		log.Fatal(err)
		return parsedPoem, err
	}
	return parsedPoem, err
}

// delineate returns the poem's text with escaped line-break characters replaced with actual line breaks.
func delineate(parsedPoem ParsedPoem) string {
	return strings.Replace(parsedPoem.Text, "\\n", "\n", -1)
}

// buildBlackout takes a poem and the blackout regex that matches it, and returns the blacked out poem as a string.
func buildBlackout(parsedPoem ParsedPoem, rp *regexp.Regexp) (string, error) {
	delinatedPoem := delineate(parsedPoem)
	if !rp.MatchString(delinatedPoem) {
		err := errors.New("Regex does not match blackout poem")
		return "", err
	}
	groups := rp.FindStringSubmatch(delinatedPoem)
	rebuiltPoem := ""
	for idx, group := range groups[1:] {
		if idx%2 == 0 {
			blackedGroup := blackoutRP.ReplaceAllString(group, "█")
			rebuiltPoem += blackedGroup
		} else {
			rebuiltPoem += group
		}
	}
	return rebuiltPoem, nil
}

// msg2regex converts a blackout poem's message into a regex string for searching poems.
func msg2regex(message string) string {
	regexString := `(?s)\A`
	for _, msgChar := range strings.Split(message, "") {
		if unicode.IsSpace(rune(msgChar[0])) {
			continue
		}
		if strings.Contains(regexEscapes, msgChar) {
			regexString += `(.*?)(\` + msgChar + `)`
		} else {
			regexString += `(.*?)(` + msgChar + `)`
		}
	}
	regexString += `(.*?)\z`
	log.Printf("Message = %s\n", message)
	log.Printf("Regex = %s\n", regexString)
	return regexString
}

// PrintParsedPoem prints the given (un-blacked-out) poem.
func PrintParsedPoem(parsedPoem ParsedPoem) {
	// print title & author
	fmt.Printf("\"%s\" by %s\n\n", parsedPoem.Title, parsedPoem.Author)
	// print lines
	lines := strings.Split(parsedPoem.Text, "\\n")
	for _, line := range lines {
		fmt.Println(line)
	}
}

// PrintBlackoutPoem prints the given blackout poem from the given poem and hidden message.
func PrintBlackoutPoem(parsedPoem ParsedPoem, message string) error {
	// convert message to regex
	regexString := msg2regex(message)
	rp := regexp.MustCompile(regexString)
	// build the blackout poem
	bp, err := buildBlackout(parsedPoem, rp)
	if err != nil {
		return err
	}
	// print the blackout poem's lines
	lines := strings.Split(bp, "\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// print the message
	fmt.Println("\n" + message)
	// print the title & author
	fmt.Printf("Excerpt of \"%s\" by %s\n\n", parsedPoem.Title, parsedPoem.Author)
	return nil
}
