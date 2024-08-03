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
	PoemsFolder string         // The file path to the poems folder.
	NPoems      int            // The number of poems in the poems folder.
	NThreads    int            // The number of goroutines to dispatch when searching.
	RP          *regexp.Regexp // The blackout regex pointer.
	MaxLength   int            // The maximum poem length [characters]
	Profanities bool           // Whether to allow profanities in searching.
}

// searchPoemsFolder searches the poems folder for poems smaller than the maximum length that match the given blackout regex.
func searchPoemsFolder(sp SearchParams) (int, error) {
	// Initialize important search parameters

	// Channels that tell each goroutine the first poem ID found.
	foundFirst := make([]chan int, sp.NThreads)
	for idx := 0; idx < sp.NThreads; idx++ {
		log.Printf("Initializing foundFirst channel %d\n", idx)
		foundFirst[idx] = make(chan int, 1)
	}
	// Channel where the goroutines send found poem IDs.
	found := make(chan int, sp.NThreads)
	// The smallest poem ID that can be blacked out.
	smallestPoemID := searchFailure

	// Dispatch the searching goroutines
	for idx := 0; idx < sp.NThreads; idx++ {
		log.Printf("Starting search goroutine #%d\n", idx)
		go searchEveryNPoems(idx, sp, found, foundFirst)
	}
	// Find the smallest poem ID from the searching goroutines
	for i := 0; i < sp.NThreads; i++ {
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
func searchEveryNPoems(startID int, sp SearchParams, found chan int, foundFirst []chan int) {
	// Initialize the ID of the first found poem
	firstFoundPoemID := searchFailure
	// Search the poem IDs that this routine is responsible for
	for poemID := startID; poemID < sp.NPoems; poemID += sp.NThreads {
		// Read the current poem
		poemPath := filepath.Join(sp.PoemsFolder, poemFilename(poemID))
		parsedPoem, readErr := json2parsedPoem(poemPath)
		if readErr != nil {
			log.Printf("Goroutine %d\t: got an error trying to read Poem %d\n", startID, poemID)
			log.Fatal(readErr)
		}
		// Check the poem's length
		if parsedPoem.Length > sp.MaxLength {
			log.Printf("Goroutine %d\t: Poem %d is too long (%d > %d)", startID, poemID, parsedPoem.Length, sp.MaxLength)
			continue
		}
		// Check if the poem's profanity level fits within search params
		if parsedPoem.IsProfane && !sp.Profanities {
			log.Printf("Goroutine %d\t: Poem %d contains profane characters", startID, poemID)
			continue
		}
		// Check if it can be blacked out
		doable, err := canBlackout(sp.RP, parsedPoem)
		if err != nil {
			log.Printf("Goroutine %d\t: got an error trying to black out Poem %d\n", startID, poemID)
			log.Fatal(err)
		}
		select {
		case first := <-foundFirst[startID]:
			log.Printf("Goroutine %d\t: got %d is the first poem ID\n", startID, first)
			firstFoundPoemID = first
		default:
			if poemID >= firstFoundPoemID-sp.NThreads {
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

// canBlackout signals whether the given poem in the poems folder can be blacked out with the regex.
func canBlackout(rp *regexp.Regexp, parsedPoem ParsedPoem) (bool, error) {
	delineatedPoem := delineate(parsedPoem)
	return rp.MatchString(delineatedPoem), nil
}
