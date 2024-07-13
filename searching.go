package main

import (
	"log"
	"path/filepath"
	"regexp"
	"runtime"
)

var numCPU = runtime.NumCPU()

const searchFailure = -1

var idChannel = make(chan int, numCPU)
var quitChannel = make(chan bool, numCPU)

func searchPoemFolder(poemFolder string, rp *regexp.Regexp, maxLength int) (int, error) {
	lengths, lengthsErr := getLengths(poemFolder)
	if lengthsErr != nil {
		return searchFailure, lengthsErr
	}
	for idx := 0; idx < numCPU; idx++ {
		log.Printf("Dispatching search goroutine #%d\n", idx)
		go searchingRoutine(idx, poemFolder, rp, lengths, maxLength)
	}
	return <-idChannel, nil
}

func searchingRoutine(routineID int, poemFolder string, rp *regexp.Regexp, lengths []int, maxLength int) {
	for poemID := routineID; poemID < len(lengths); poemID += numCPU {
		select {
		case <-quitChannel:
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
				idChannel <- poemID
				close(idChannel)
				for idx := 0; idx < numCPU; idx++ {
					quitChannel <- true
				}
				return
			}
		}
	}
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
	return rp.MatchString(poem.Text), nil
}
