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
	// cacheFolder is this program's cache folder. On Linux systems, it would be `~/.cache/blackout`.
	cacheFolder = filepath.Join(xdg.CacheHome, "blackout")
	// Local path to public domain poetry dataset JSON file.
	cacheFolderJSON = filepath.Join(cacheFolder, "poems.json")
	// Directory where the parsed poem JSONs are stored.
	cacheFolderPoems = filepath.Join(cacheFolder, "poems")
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
func downloadPoemsJSON(poemsPath string) (err error) {
	// Make parent directory if it doesn't exist
	dir, _ := filepath.Split(poemsPath)
	err = os.MkdirAll(dir, 0o750)
	if err != nil {
		return err
	}
	// Check if file exists
	_, err = os.Stat(poemsPath)
	if err == nil {
		log.Printf("File already exists at %s\n", poemsPath)
		return nil
	}
	if errors.Is(err, os.ErrNotExist) {
		log.Println("Downloading poem dataset")
		resp, err := http.Get(poemsURL)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// Safely close HTTP response
		defer func() {
			err = resp.Body.Close()
		}()
		if err != nil {
			return err
		}
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
	return err
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

// setupDataFolder sets up this CLI application's cache folder.
func setupDataFolder() error {
	// Make the cache folder if it doesn't already exist
	_, folderErr := os.Stat(cacheFolder)
	if os.IsNotExist(folderErr) {
		log.Printf("Creating cache folder %s\n", cacheFolder)
		dirErr := os.Mkdir(cacheFolder, 0o750)
		if dirErr != nil {
			return dirErr
		}
	} else {
		log.Printf("Data folder %s already exists\n", cacheFolder)
	}
	// Download the poem database, and put it in the cache folder
	dlErr := downloadPoemsJSON(cacheFolderJSON)
	if dlErr != nil {
		return dlErr
	}
	// Populate the "poems" folder in the cache folder if not already done
	_, poemsFolderErr := os.Stat(filepath.Join(cacheFolder, "poems"))
	if os.IsNotExist(poemsFolderErr) {
		poems, readErr := readPoemsJSON(filepath.Join(cacheFolder, "poems.json"))
		if readErr != nil {
			return readErr
		}
		splitErr := parsePoems(poems, filepath.Join(cacheFolder, "poems"))
		if splitErr != nil {
			return splitErr
		}
	}
	return nil
}
