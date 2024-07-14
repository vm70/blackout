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
)

// regexEscapes contains the characters that need to be escaped in regexes.
const regexEscapes = `.+*?()|[]{}^$`

// blackoutRP is the regular expression pointer that matches every non-whitespace charater for blacking out.
var blackoutRP = regexp.MustCompile(`[^\t\f\r\n\ ]`)

// A Poem in the database has a title, an author, and text.
type Poem struct {
	// The title of the poem.
	Title string
	// The author of the poem.
	Author string
	// The poem's text itself. Poem lines are delineated with an escaped line break character "\n".
	Text string
}

// poem2json writes the Poem struct to the given JSON file path.
func poem2json(poem Poem, jsonFile string) error {
  // TODO Turn this into a method for `Poem`.
	// Marshal to JSON bytes
	poemBytes, err := json.Marshal(poem)
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

// json2poem extracts the Poem object from the given JSON file path.
func json2poem(jsonFile string) (Poem, error) {
	// Read the file name
	var poem Poem
	fileBytes, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
		return poem, err
	}
	// parse JSON and return poem
	err = json.Unmarshal(fileBytes, &poem)
	if err != nil {
		log.Fatal(err)
		return poem, err
	}
	return poem, err
}

// delineate returns the poem's text with escaped line-break characters replaced with actual line breaks.
func delineate(poem Poem) string {
	return strings.Replace(poem.Text, "\\n", "\n", -1)
}

// buildBlackout takes a poem and the blackout regex that matches it, and returns the blacked out poem as a string.
func buildBlackout(poem Poem, rp *regexp.Regexp) (string, error) {
	delinatedPoem := delineate(poem)
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
	regexString := "(?ms)^"

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
	regexString += "(.*?)$"
	return regexString
}

// PrintPoem prints the given (un-blacked-out) poem.
func PrintPoem(poem Poem) {
  // print title & author
	fmt.Printf("\"%s\" by %s\n\n", poem.Title, poem.Author)
  // print lines
	lines := strings.Split(poem.Text, "\\n")
	for _, line := range lines {
		fmt.Println(line)
	}
}

// PrintBlackoutPoem prints the given blackout poem from the given poem and hidden message.
func PrintBlackoutPoem(poem Poem, message string) error {
  // convert message to regex
  regexString := msg2regex(message)
  rp := regexp.MustCompile(regexString)
  // build the blackout poem
	bp, err := buildBlackout(poem, rp)
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
	fmt.Printf("Excerpt of \"%s\" by %s\n\n", poem.Title, poem.Author)
	return nil
}
