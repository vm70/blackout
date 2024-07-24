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
	"errors"
	"log"
	"path/filepath"
	"regexp"
)

// default return value for when a searching function fails; the maximum integer value.
const searchFailure = int(^uint(0) >> 1)

// A SearchParams struct contains common information about the folder searching operation.
type SearchParams struct {
	NRoutines     int            // The number of goroutines to dispatch when searching.
	PoemsFolder   string         // The file path to the poems folder.
	RP            *regexp.Regexp // The blackout regex pointer.
	MaxPoemLength int            // The maximum poem length [characters].
}

// searchPoemsFolder searches the poems folder for poems smaller than the maximum length that match the given blackout regex.
func searchPoemsFolder(sp SearchParams) (int, error) {
	// Initialize important search parameters

	// Channels that tell each goroutine the first poem ID found.
	foundFirst := make([]chan int, sp.NRoutines)
	for idx := 0; idx < sp.NRoutines; idx++ {
		log.Printf("Initializing foundFirst channel %d\n", idx)
		foundFirst[idx] = make(chan int, 1)
	}
	// Channel where the goroutines send found poem IDs.
	found := make(chan int, sp.NRoutines)
	// The smallest poem ID that can be blacked out.
	smallestPoemID := searchFailure

	// Get the lengths from the poems folder
	lengths, lengthsErr := getLengths(sp.PoemsFolder)
	if lengthsErr != nil {
		return searchFailure, lengthsErr
	}
	// Dispatch the searching goroutines
	for idx := 0; idx < sp.NRoutines; idx++ {
		log.Printf("Starting search goroutine #%d\n", idx)
		go searchEveryNPoems(idx, lengths, sp, found, foundFirst)
	}
	// Find the smallest poem ID from the searching goroutines
	for i := 0; i < sp.NRoutines; i++ {
		foundID := <-found
		smallestPoemID = min(smallestPoemID, foundID)
		log.Printf("Main thread\t: received %d; earliest poem in index to black out has ID %d\n", foundID, smallestPoemID)
	}
	// If all goroutines fail, then the smallest poem ID is invalid
	if smallestPoemID == searchFailure {
		searchErr := errors.New("Failed to find a blackout poem")
		return searchFailure, searchErr
	}
	return smallestPoemID, nil
}

// searchEveryNPoems is a goroutine that searches every `sp.NRoutines` poems for one that is shorter than the maximum length and matches the blackout regex, starting from the poem at index `startID`. If it finds a poem to black out, then it sends the poem ID through the `found` channel. If it is the first (by time) to find a poem, then it sends that ID through the `foundFirst` channels.
func searchEveryNPoems(startID int, lengths []int, sp SearchParams, found chan int, foundFirst []chan int) {
	// Initialize the ID of the first found poem
	firstFoundPoemID := searchFailure
	for poemID := startID; poemID < len(lengths); poemID += sp.NRoutines {
		log.Printf("Goroutine %d\t: Checking Poem ID %d\n", startID, poemID)
		doable, err := canBlackout(sp.RP, sp.PoemsFolder, poemID, lengths[poemID], sp.MaxPoemLength)
		if err != nil {
			log.Printf("Goroutine %d\t: got an error\n", startID)
			log.Fatal(err)
		}
		select {
		case first := <-foundFirst[startID]:
			log.Printf("Goroutine %d\t: got %d is the first poem ID\n", startID, first)
			firstFoundPoemID = first
		default:
			if poemID >= firstFoundPoemID-sp.NRoutines {
				log.Printf("Goroutine %d\t: did not find an earlier poem than ID %d; stopping\n", startID, firstFoundPoemID)
				found <- searchFailure
				return
			}
			if doable {
				log.Printf("Goroutine %d\t: found poem %d to black out\n", startID, poemID)
				if firstFoundPoemID != searchFailure {
					for _, ff := range foundFirst {
						ff <- poemID
					}
				}
				found <- poemID
				log.Printf("Goroutine %d\t: stopping", startID)
				return
			}
		}
	}
	log.Printf("Goroutine %d\t: failed to find a poem\n", startID)
	found <- searchFailure
}

// canBlackout reports whether the given poem in the poem folder is shorter than the maximum length and can be blacked out by the blackout regex.
func canBlackout(rp *regexp.Regexp, poemFolder string, poemID int, poemLength int, maxLength int) (bool, error) {
	if poemLength > maxLength {
		return false, nil
	}
	poemPath := filepath.Join(poemFolder, poemFilename(poemID))
	poem, err := json2poem(poemPath)
	if err != nil {
		return false, err
	}
	delineatedPoem := delineate(poem)
	return rp.MatchString(delineatedPoem), nil
}
