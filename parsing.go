package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const poemsJson = "https://huggingface.co/datasets/DanFosing/public-domain-poetry/resolve/main/poems.json"

var dataFolder = filepath.Join(xdg.DataHome, "blackout")

var poems_sha256 = [32]byte{0x17, 0x2c, 0xd2, 0xc5, 0xd9, 0x53, 0xc7, 0x02, 0x33, 0x90, 0xa8, 0xd1, 0xf3, 0x37, 0xd0, 0x23, 0xd7, 0xfb, 0xb2, 0xb9, 0x25, 0xdf, 0x0a, 0x66, 0xd0, 0x22, 0x1f, 0x30, 0xc6, 0xad, 0xc3, 0x08}

func poemsFileHashMatches(filename string) error {
	content, readErr := os.ReadFile(filename)
	if readErr != nil {
		return readErr
	}
	return poemsBytesHashMatches(content)
}

func poemsBytesHashMatches(fileBytes []byte) error {
	resp_sum := sha256.Sum256(fileBytes)
	if resp_sum != poems_sha256 {
		errorStr := fmt.Sprintf("Hash %x doesn't match reference %x", resp_sum, resp_sum)
		return errors.New(errorStr)
	}
	return nil
}

func downloadPoems(filename string) error {
	// Check if file exists
	_, fileErr := os.Stat(filename)
	if fileErr == nil {
		log.Printf("File already exists at %s\n", filename)
    return poemsFileHashMatches(filename)
	}
	if errors.Is(fileErr, os.ErrNotExist) {
		log.Println("Downloading poem dataset")
		resp, getErr := http.Get(poemsJson)
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
		file, createErr := os.Create(filename)
		if createErr != nil {
			return createErr
		}
		_, writeErr := io.WriteString(file, string(body))
		if writeErr != nil {
			return writeErr
		}
		defer file.Close()
		return nil
	}
	return fileErr
}

func splitPoems(poemsJson string, poemsFolder string) {

}
