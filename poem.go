package main

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

const regexEscapes = `.+*?()|[]{}^$`

var blackoutRP = regexp.MustCompile(`[^\t\f\r\n\ ]`)

type Poem struct {
	Title  string
	Author string
	Text   string
}

func poem2json(poem Poem, filename string) error {
	// Marshal to JSON bytes
	poemBytes, err := json.Marshal(poem)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Write file
	err = os.WriteFile(filename, poemBytes, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func json2poem(filename string) (Poem, error) {
	// Read the file name
	var poem Poem
	fileBytes, err := os.ReadFile(filename)
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

func delineate(poem Poem) string {
	return strings.Replace(poem.Text, "\\n", "\n", -1)
}

func blackout(poem Poem, rp *regexp.Regexp) (string, error) {
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

func printPoem(poem Poem) {
  fmt.Printf("\"%s\" by %s\n\n", poem.Title, poem.Author)
  lines := strings.Split(poem.Text, "\\n")
  for _, line := range lines {
    fmt.Println(line)
  }
}

func printBlackoutPoem(poem Poem, rp *regexp.Regexp, message string) error {
  bp, err := blackout(poem, rp)
  if err != nil {
    return err
  }
  lines := strings.Split(bp, "\n")
  for _, line := range lines {
    fmt.Println(line)
  }
  fmt.Println("\n" + message)
  fmt.Printf("Excerpt of \"%s\" by %s\n\n", poem.Title, poem.Author)
  return nil
}
