package main

import "log"

func main() {
	log.Printf("Data Folder is %s\n", dataFolder)
	downloadPoems("poems.json")
	poems, err := readPoemDB("poems.json")
	if err != nil {
		log.Fatal(err)
	}
	err = splitPoems(poems, "poem_folder")
	if err != nil {
		log.Fatal(err)
	}
	err = readPoemFolder("poem_folder")
	if err != nil {
		log.Fatal(err)
	}
}
