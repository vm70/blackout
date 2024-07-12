package main

import (
	"encoding/json"
	"log"
	"os"
)

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
