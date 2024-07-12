package main

import (
	"os"
	"regexp"
	"runtime"
)

func checkLengths() {
}

func searchPoemFolder(poemFolder string) error {
	var numCPUs = runtime.NumCPU()
	lengths, lengthsErr := getLengths(poemFolder)
	if lengthsErr != nil {
		return lengthsErr
	}
  entries, err := os.ReadDir()
  if 
  for idx := 0; idx < numCPUs; idx++ {
  }
	return nil
}

func canBlackout(rgx regexp.Regexp, poem Poem, poemLength int, maxLength int) bool {
	if poemLength > maxLength {
		return false
	}
	return rgx.MatchString(poem.Text)
}
