package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const poems_json = "https://huggingface.co/datasets/DanFosing/public-domain-poetry/resolve/main/poems.json"

var poems_location = filepath.Join(xdg.DataHome, "blackout")

var poems_sha256 = [32]byte{0x17, 0x2c, 0xd2, 0xc5, 0xd9, 0x53, 0xc7, 0x02, 0x33, 0x90, 0xa8, 0xd1, 0xf3, 0x37, 0xd0, 0x23, 0xd7, 0xfb, 0xb2, 0xb9, 0x25, 0xdf, 0x0a, 0x66, 0xd0, 0x22, 0x1f, 0x30, 0xc6, 0xad, 0xc3, 0x08}

func downloadPoems(filename string) error {
	// Check if file exists
	_, err := os.Stat(filename)
	if err == nil {
		// Do nothing if the file already has been downloaded
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		// Download the poem dataset
		resp, err := http.Get(poems_json)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)

		// Check SHA256 sum
		resp_sum := sha256.Sum256(body)
		if resp_sum != poems_sha256 {
			error_str := fmt.Sprintf("Hash %x doesn't match reference %x", resp_sum, resp_sum)
			return errors.New(error_str)
		}

		// Create local file
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		// Write to the local file
		_, err = io.WriteString(file, string(body))
		return err
	} else {
		return err
	}
}
