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
	Length int
}

func poem2json(poem Poem, name string) error {
	// Marshal to JSON bytes
	poem_bytes, err := json.Marshal(poem)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Write file
	err = os.WriteFile(name, poem_bytes, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func json2poem(name string) (Poem, error) {
	// Read the file name
	var poem Poem
	file_bytes, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
		return poem, err
	}
	// parse JSON and return poem
	err = json.Unmarshal(file_bytes, &poem)
	if err != nil {
		log.Fatal(err)
		return poem, err
	}
	return poem, err
}
