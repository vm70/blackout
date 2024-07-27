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
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

// lengthsFilename is the name of the file in the poems folder where each poem's length (in characters) is stored.
const lengthsFilename = "lengths.csv"

// writeLengths writes the poems folder's lengths (stored as a string array) to the CSV stored in the poem folder.
func writeLengths(lengths []string, poemsFolder string) error {
	lengthsPath := filepath.Join(poemsFolder, lengthsFilename)
	file, createErr := os.Create(lengthsPath)
	if createErr != nil {
		return createErr
	}
	lengthsRecord := [][]string{lengths}
	csvWriter := csv.NewWriter(file)
	writeErr := csvWriter.WriteAll(lengthsRecord)
	return writeErr
}

// getLengths reads the lengths file in the given poems folder and returns the array of corresponding poem lengths.
func getLengths(poemsFolder string) ([]int, error) {
	lengths := []int{}
	lengthsPath := filepath.Join(poemsFolder, lengthsFilename)
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
