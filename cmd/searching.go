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
	"runtime"
)

// number of CPUs that the host computer has.
var numCPU = runtime.NumCPU()

// default return value for when a searching function fails; the maximum integer value.
const searchFailure = int(^uint(0) >> 1)

var (
	// Slice of channels for telling each goroutine the first (earliest in time) poem found.
	foundFirst = make([]chan int, numCPU)
	// Channel for sending poem IDs that match the given blackout regex.
	found = make(chan int, numCPU)
)

// searchPoemsFolder searches the poems folder for poems smaller than the
// maximum length that match the given blackout regex.
func searchPoemsFolder(poemsFolder string, rp *regexp.Regexp, maxLength int) (int, error) {
	// Get the lengths from the poems folder
	lengths, lengthsErr := getLengths(poemsFolder)
	if lengthsErr != nil {
		return searchFailure, lengthsErr
	}
	// Initialize foundFirst channel slice
	for i := 0; i < numCPU; i++ {
		log.Printf("Initializing foundFirst channel %d\n", i)
		foundFirst[i] = make(chan int, 1)
	}
	// Dispatch the goroutines to search the poems in parallel
	for i := 0; i < numCPU; i++ {
		log.Printf("Starting search goroutine #%d\n", i)
		go searchEveryNPoems(i, poemsFolder, rp, lengths, maxLength)
	}
	earliestPoemID := searchFailure
	for i := 0; i < numCPU; i++ {
		foundID := <-found
		earliestPoemID = min(earliestPoemID, foundID)
		log.Printf("Main thread:\t received %d; earliest poem in index to black out has ID %d\n", foundID, earliestPoemID)
	}
	if earliestPoemID == searchFailure {
		searchErr := errors.New("Failed to find a blackout poem")
		return searchFailure, searchErr
	}
	return earliestPoemID, nil
}

// searchEveryNPoems is a goroutine that searches every `numCPU` poems for one that is
// shorter than the maximum length and matches the blackout regex, starting from
// the poem at index `startID` and stopping when there are no more poems.
func searchEveryNPoems(startID int, poemFolder string, rp *regexp.Regexp, lengths []int, maxLength int) {
	earliestPoemID := searchFailure
	for poemID := startID; poemID < len(lengths); poemID += numCPU {
		log.Printf("Goroutine %d:\t Checking Poem ID %d\n", startID, poemID)
		doable, err := canBlackout(rp, poemFolder, poemID, lengths[poemID], maxLength)
		if err != nil {
			log.Printf("Goroutine %d:\t got an error\n", startID)
			log.Fatal(err)
		}
		select {
		case first := <-foundFirst[startID]:
			log.Printf("Goroutine %d:\t got %d is the first poem ID", startID, first)
			earliestPoemID = first
		default:
			if poemID >= earliestPoemID-numCPU {
				log.Printf("Goroutine %d did not find an earlier poem than ID %d; stopping", startID, earliestPoemID)
				found <- searchFailure
				return
			}
			if doable {
				log.Printf("Goroutine %d found poem %d to black out\n", startID, poemID)
				if earliestPoemID != searchFailure {
					for idx := 0; idx < numCPU; idx++ {
						foundFirst[idx] <- poemID
					}
				}
				found <- poemID
				log.Printf("Goroutine %d stopping", startID)
				return
			}
		}
	}
	found <- searchFailure
}

// canBlackout reports whether the given poem in the poem folder is shorter
// than the maximum length and can be blacked out by the blackout regex.
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
