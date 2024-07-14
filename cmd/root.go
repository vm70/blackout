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
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/spf13/cobra"
)

// BlackoutVersion is the version number of `blackout`.
const BlackoutVersion = "0.1.0"

// Verbose determines whether to print verbose results.
var (
	Verbose bool
	// MaxLength determines the maximum poem length to black out.
	MaxLength int
	// PrintOriginal determines whether to print the original poem before blacking it out.
	PrintOriginal bool
)

// rootCmd represents the base command when called without any sub-commands.
var rootCmd = &cobra.Command{
	Use:     "blackout <message>",
	Short:   "Make a blackout poem with the given hidden message",
	Version: BlackoutVersion,
	Args:    cobra.ExactArgs(1),
	Run:     runApp,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init sets up the flags of the CLI application.
func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "verbose output")
	rootCmd.PersistentFlags().IntVarP(&MaxLength, "max-length", "l", 400, "maximum poem length")
	rootCmd.PersistentFlags().BoolVarP(&PrintOriginal, "print-original", "p", false, "print original poem before blacking out")
}

// runApp runs the CLI application.
func runApp(cmd *cobra.Command, args []string) {
	if !Verbose {
		log.SetOutput(io.Discard)
	}
	log.Printf("Running %s\n", cmd.Name())
	log.Printf("Data Folder is %s\n", dataFolder)
	regexpString := msg2regex(args[0])
	blackoutRegex := regexp.MustCompile(regexpString)
	log.Printf("Message = %s\n", args[0])
	log.Printf("Regex = %s\n", regexpString)
	setupErr := setupDataFolder()
	if setupErr != nil {
		log.Fatalf(setupErr.Error())
	}
	poemID, err := searchPoemsFolder("poem_folder", blackoutRegex, MaxLength)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found poem ID %d to black out\n", poemID)
	poem, err := json2poem(filepath.Join("poem_folder", poemFilename(poemID)))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Poem ID %d is \"%s\"\n", poemID, poem.Title)
	if Verbose {
		time.Sleep(1 * time.Second)
	}
	if PrintOriginal {
		PrintPoem(poem)
		print("\n")
	}
	PrintBlackoutPoem(poem, args[0])
}
