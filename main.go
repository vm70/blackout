package main

import (
	"log"
	"regexp"
)

var blackoutRP, _ = regexp.Compile("(.*?)(b)(.*?)(l)(.*?)(a)(.*?)(c)(.*?)(k)(.*?)(o)(.*?)(u)(.*?)(t)(.*?)(p)(.*?)(o)(.*?)(e)(.*?)(m)(.*?)")

func main() {
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
	poemID, err := searchPoemFolder("poem_folder", blackoutRP, 400)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found poem ID %d to black out\n", poemID)
// 	for i := 0; i < 100; i++ {
// 		log.Printf("Default %d\n", i)
// 		time.Sleep(10 * time.Millisecond)
// 	}
}
