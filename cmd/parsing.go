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

// poemsURL is the online URL where the public domain poetry database JSON file is stored.
const poemsURL = "https://huggingface.co/datasets/DanFosing/public-domain-poetry/resolve/main/poems.json"

// poemsSha256 is the SHA256 hash of the poem dataase JSON: `172cd2c5d953c7023390a8d1f337d023d7fbb2b925df0a66d0221f30c6adc308`.
var poemsSha256 = [32]byte{0x17, 0x2c, 0xd2, 0xc5, 0xd9, 0x53, 0xc7, 0x02, 0x33, 0x90, 0xa8, 0xd1, 0xf3, 0x37, 0xd0, 0x23, 0xd7, 0xfb, 0xb2, 0xb9, 0x25, 0xdf, 0x0a, 0x66, 0xd0, 0x22, 0x1f, 0x30, 0xc6, 0xad, 0xc3, 0x08}

// dataFolder is this program's data folder. On Linux systems, it would be `~/.local/share/blackout`.
var dataFolder = filepath.Join(xdg.DataHome, "blackout")

// poemsFileHashMatches returns an error if the given file's SHA256 hash doesn't match the hard-coded one above.
func poemsFileHashMatches(filename string) error {
	content, readErr := os.ReadFile(filename)
	if readErr != nil {
		return readErr
	}
	return poemsBytesHashMatches(content)
}

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

func downloadPoems(filename string) error {
	// Check if file exists
	_, fileErr := os.Stat(filename)
	if fileErr == nil {
		log.Printf("File already exists at %s\n", filename)
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
		hashErr := poemsBytesHashMatches(body)
		if hashErr != nil {
			return hashErr
		}
		writeErr := os.WriteFile(filename, body, 0o666)
		if writeErr != nil {
			return writeErr
		}
		return nil
	}
	return fileErr
}

func poemFilename(poemID int) string {
	return "poem" + strconv.Itoa(poemID) + ".json"
}

func splitPoems(poems []Poem, poemFolder string) error {
	_, folderErr := os.Stat(poemFolder)
	if os.IsNotExist(folderErr) {
		log.Printf("Creating poem folder %s\n", poemFolder)
		os.Mkdir(poemFolder, 0o750)
	} else {
		log.Printf("Poem folder %s already exists\n", poemFolder)
		return nil
	}
	lengths := []string{}
	for idx, poem := range poems {
		lengths = append(lengths, fmt.Sprintf("%d", len(poem.Text)))
		poemJSON := filepath.Join(poemFolder, poemFilename(idx))
		poemErr := poem2json(poem, poemJSON)
		if poemErr != nil {
			return poemErr
		}
	}
	return writeLengths(lengths, poemFolder)
}
