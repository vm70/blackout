package main

import (
	"log"
	"path/filepath"
	"regexp"
	"time"
)

var bpMessage = "blackout poem"
var poemRP = regexp.MustCompile(msg2regex(bpMessage))

func main_old() {
	log.Printf("Data Folder is %s\n", dataFolder)
	downloadPoems("poems.json")
	dbPoems, err := readPoemDB("poems.json")
	if err != nil {
		log.Fatal(err)
	}
	err = splitPoems(dbPoems, "poem_folder")
	if err != nil {
		log.Fatal(err)
	}
	poemID, err := searchPoemFolder("poem_folder", poemRP, 400)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found poem ID %d to black out\n", poemID)
	poem, err := json2poem(filepath.Join("poem_folder", poemFilename(poemID)))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Poem ID %d is \"%s\"\n", poemID, poem.Title)
	time.Sleep(1 * time.Second)
	print("\n\n\n")
	printPoem(poem)
	print("\n\n\n")
	printBlackoutPoem(poem, poemRP, bpMessage)
	print("\n\n\n")
}
