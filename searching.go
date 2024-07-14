package main

import (
	"errors"
	"log"
	"path/filepath"
	"regexp"
	"runtime"
)

var numCPU = runtime.NumCPU()

const searchFailure = -1

var found = make(chan int, numCPU)
var failed = make(chan bool, numCPU)
var stop = make(chan bool, numCPU)

func searchPoemFolder(poemFolder string, rp *regexp.Regexp, maxLength int) (int, error) {
	lengths, lengthsErr := getLengths(poemFolder)
	if lengthsErr != nil {
		return searchFailure, lengthsErr
	}
	for idx := 0; idx < numCPU; idx++ {
		log.Printf("Starting search goroutine #%d\n", idx)
		go searchingRoutine(idx, poemFolder, rp, lengths, maxLength)
	}
	var nFailed = 0
	for {
		select {
		case poemID := <-found:
			log.Printf("Received poem ID %d to black out\n", poemID)
			for idx := 0; idx < numCPU; idx++ {
				stop <- true
			}
			return poemID, nil
		case <-failed:
			nFailed += 1
			if nFailed == numCPU {
				searchErr := errors.New("Failed to find a blackout poem")
				return searchFailure, searchErr
			}
		}
	}
}

func searchingRoutine(routineID int, poemFolder string, rp *regexp.Regexp, lengths []int, maxLength int) {
	for poemID := routineID; poemID < len(lengths); poemID += numCPU {
		select {
		case <-stop:
			log.Printf("Goroutine %d stopped\n", routineID)
			return
		default:
			log.Printf("Goroutine %d Checking Poem #%d\n", routineID, poemID)
			doable, err := canBlackout(rp, poemFolder, poemID, lengths[poemID], maxLength)
			if err != nil {
				log.Printf("Goroutine %d had an error\n", routineID)
				log.Fatal(err)
			}
			if doable {
				log.Printf("Goroutine %d found poem %d to black out\n", routineID, poemID)
				found <- poemID
				log.Printf("Goroutine %d stopping", routineID)
        return
			}
		}
	}
	failed <- true
}

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
