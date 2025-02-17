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
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/adrg/xdg"
)

// The online URL where the public domain poetry database JSON file is stored.
const poemsURL = "https://huggingface.co/datasets/DanFosing/public-domain-poetry/resolve/main/poems.json"

var (
	// SHA256 hash of the poem database JSON.
	poemsSha256 = [32]byte{0x17, 0x2c, 0xd2, 0xc5, 0xd9, 0x53, 0xc7, 0x02, 0x33, 0x90, 0xa8, 0xd1, 0xf3, 0x37, 0xd0, 0x23, 0xd7, 0xfb, 0xb2, 0xb9, 0x25, 0xdf, 0x0a, 0x66, 0xd0, 0x22, 0x1f, 0x30, 0xc6, 0xad, 0xc3, 0x08}
	// dataFolder is this program's data folder. On Linux systems, it would be `~/.local/share/blackout`.
	dataFolder = filepath.Join(xdg.DataHome, "blackout")
	// Local path to public domain poetry dataset JSON file.
	dataFolderJSON = filepath.Join(dataFolder, "poems.json")
	// Directory where the parsed poem JSONs are stored.
	dataFolderPoems = filepath.Join(dataFolder, "poems")
)

// poemsBytesHashMatches returns an error if the given byte array's SHA256 hash doesn't match the hard-coded one above.
func poemsBytesHashMatches(fileBytes []byte) error {
	respSum := sha256.Sum256(fileBytes)
	if respSum != poemsSha256 {
		errorStr := fmt.Sprintf("Hash %x doesn't match reference %x", respSum, respSum)
		return errors.New(errorStr)
	}
	return nil
}

// readPoemsJSON reads the poem database JSON file and converts it into an array of Poems.
func readPoemsJSON(poemsJSON string) ([]Poem, error) {
	// Read the file name
	var poemArr []Poem
	fileBytes, err := os.ReadFile(poemsJSON)
	if err != nil {
		return poemArr, err
	}
	// parse JSON and return poem
	err = json.Unmarshal(fileBytes, &poemArr)
	if err != nil {
		log.Fatal(err)
		return poemArr, err
	}
	return poemArr, nil
}

// downloadPoemsJSON downloads the poem JSON file and places it in the given path. If the poems JSON file already exists (from a previous run), then it returns nil.
func downloadPoemsJSON(poemsPath string) error {
	// Make parent directory if it doesn't exist
	dir, _ := filepath.Split(poemsPath)
	mkdirErr := os.MkdirAll(dir, 0o750)
	if mkdirErr != nil {
		return mkdirErr
	}
	// Check if file exists
	_, fileErr := os.Stat(poemsPath)
	if fileErr == nil {
		log.Printf("File already exists at %s\n", poemsPath)
		return nil
	}
	if errors.Is(fileErr, os.ErrNotExist) {
		log.Println("Downloading poem dataset")
		resp, getErr := http.Get(poemsURL)
		if getErr != nil {
			return getErr
		}
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return readErr
		}
		defer resp.Body.Close()
		hashErr := poemsBytesHashMatches(body)
		if hashErr != nil {
			return hashErr
		}
		writeErr := os.WriteFile(poemsPath, body, 0o666)
		if writeErr != nil {
			return writeErr
		}
		return nil
	}
	return fileErr
}

// poemFilename returns the poem's file name by its ID.
func poemFilename(poemID int) string {
	return "poem" + strconv.Itoa(poemID) + ".json"
}

// Parse an array of poems, and split them into JSON files in the poems folder.
func parsePoems(poems []Poem, poemsFolder string) error {
	_, folderErr := os.Stat(poemsFolder)
	if os.IsNotExist(folderErr) {
		log.Printf("Creating poems folder %s\n", poemsFolder)
		dirErr := os.Mkdir(poemsFolder, 0o750)
		if dirErr != nil {
			return dirErr
		}
	} else {
		log.Printf("Poems folder %s already exists\n", poemsFolder)
		return nil
	}
	for idx, poem := range poems {
		parsedPoem := NewParsedPoem(poem)
		poemJSON := filepath.Join(poemsFolder, poemFilename(idx))
		poemErr := parsedPoem2json(parsedPoem, poemJSON)
		if poemErr != nil {
			return poemErr
		}
	}
	return nil
}

// setupDataFolder sets up this CLI application's data folder.
func setupDataFolder() error {
	// Make the data folder if it doesn't already exist
	_, folderErr := os.Stat(dataFolder)
	if os.IsNotExist(folderErr) {
		log.Printf("Creating data folder %s\n", dataFolder)
		dirErr := os.Mkdir(dataFolder, 0o750)
		if dirErr != nil {
			return dirErr
		}
	} else {
		log.Printf("Data folder %s already exists\n", dataFolder)
	}
	// Download the poem database, and put it in the data folder
	dlErr := downloadPoemsJSON(dataFolderJSON)
	if dlErr != nil {
		return dlErr
	}
	// Populate the "poems" folder in the data folder if not already done
	_, poemsFolderErr := os.Stat(filepath.Join(dataFolder, "poems"))
	if os.IsNotExist(poemsFolderErr) {
		poems, readErr := readPoemsJSON(filepath.Join(dataFolder, "poems.json"))
		if readErr != nil {
			return readErr
		}
		splitErr := parsePoems(poems, filepath.Join(dataFolder, "poems"))
		if splitErr != nil {
			return splitErr
		}
	}
	return nil
}
