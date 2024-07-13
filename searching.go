package main

import (
	"log"
	"regexp"
)

func checkLengths() {
}

// func searchPoemFolder(poemFolder string) error {
// 	var numCPUs = runtime.NumCPU()
// 	lengths, lengthsErr := getLengths(poemFolder)
// 	if lengthsErr != nil {
// 		return lengthsErr
// 	}
//   entries, err := os.ReadDir()
//   if
//   for idx := 0; idx < numCPUs; idx++ {
//
//   }
// 	return nil
// }

func searchingRoutine(routineID int, numCPU int, rp *regexp.Regexp, poems []Poem, poemLengths []int, maxLength int) bool {
	for idx := routineID; idx < len(poems); idx += numCPU {
		log.Printf("Checking Poem #%d\n", idx)
		if canBlackout(rp, poems[idx], poemLengths[idx], maxLength) {
			return true
		}
	}
	return false
}

func canBlackout(rp *regexp.Regexp, poem Poem, poemLength int, maxLength int) bool {
	if poemLength > maxLength {
		return false
	}
	return rp.MatchString(poem.Text)
}
