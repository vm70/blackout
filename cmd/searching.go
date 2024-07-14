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

// default return value for when a searching function fails.
const searchFailure = -1

var (
	// Channel for sending poem IDs that match the given blackout regex.
	found = make(chan int, numCPU)
	// Channel for signifying when a searching goroutine has failed to find a blackout poem.
	failed = make(chan bool, numCPU)
	// Channel for telling the goroutines to stop searching for poems to black out.
	stop = make(chan bool, numCPU)
)

// searchPoemsFolder searches the poems folder for poems smaller than the
// maximum length that match the given blackout regex.
func searchPoemsFolder(poemsFolder string, rp *regexp.Regexp, maxLength int) (int, error) {
	// Get the lengths from the poems folder
	lengths, lengthsErr := getLengths(poemsFolder)
	if lengthsErr != nil {
		return searchFailure, lengthsErr
	}
	// Dispatch the goroutines to search the poems in parallel
	for i := 0; i < numCPU; i++ {
		log.Printf("Starting search goroutine #%d\n", i)
		go searchEveryNPoems(i, poemsFolder, rp, lengths, maxLength)
	}
	nFailed := 0 // How many goroutines have failed.
	for {
		select {
		case poemID := <-found:
			log.Printf("Received poem ID %d to black out\n", poemID)
			for idx := 0; idx < numCPU; idx++ {
				stop <- true
			}
			return poemID, nil
		case <-failed:
			nFailed++
			if nFailed >= numCPU {
				searchErr := errors.New("Failed to find a blackout poem")
				return searchFailure, searchErr
			}
		}
	}
}

// searchEveryNPoems is a goroutine that searches every `numCPU` poems for one that is
// shorter than the maximum length and matches the blackout regex, starting from
// the poem at index `startID` and stopping when there are no more poems.
func searchEveryNPoems(startID int, poemFolder string, rp *regexp.Regexp, lengths []int, maxLength int) {
	for poemID := startID; poemID < len(lengths); poemID += numCPU {
		select {
		case <-stop:
			log.Printf("Goroutine %d stopped\n", startID)
			return
		default:
			log.Printf("Goroutine %d Checking Poem #%d\n", startID, poemID)
			doable, err := canBlackout(rp, poemFolder, poemID, lengths[poemID], maxLength)
			if err != nil {
				log.Printf("Goroutine %d had an error\n", startID)
				log.Fatal(err)
			}
			if doable {
				log.Printf("Goroutine %d found poem %d to black out\n", startID, poemID)
				found <- poemID
				log.Printf("Goroutine %d stopping", startID)
				return
			}
		}
	}
	failed <- true
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
