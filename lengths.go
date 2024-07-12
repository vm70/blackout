package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

const lengthsFile = "lengths.csv"

func writeLengths(lengths []string, poemFolder string) error {
  lengthsPath := filepath.Join(poemFolder, lengthsFile)
  file, createErr := os.Create(lengthsPath)
  if createErr != nil {
    return createErr
  }
  var lengthsRecord = [][]string{lengths}
  csvWriter := csv.NewWriter(file)
  writeErr := csvWriter.WriteAll(lengthsRecord)
  return writeErr
}

func getLengths(poemFolder string) ([]int, error) {
  var lengths = []int{}
  lengthsPath := filepath.Join(poemFolder, lengthsFile)
  file, openErr := os.Open(lengthsPath)
  if openErr != nil {
    return lengths, openErr
  }
  csvReader := csv.NewReader(file)
  data, readErr := csvReader.ReadAll()
  if readErr != nil {
    return lengths, readErr
  }
  for _, lengthString := range data[0] {
    length, convErr := strconv.Atoi(lengthString)
    if convErr != nil {
      return lengths, convErr
    }
    lengths = append(lengths, length)
  }
  return lengths, nil
}


