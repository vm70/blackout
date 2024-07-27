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
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

const longDescription = `Blackout is a command-line application that automates the process of making
simple blackout poems, where characters and words of an an original text source
are removed to create an entirely new piece. It combs through a database of
public-domain poetry to find one with characters that match a given message,
then prints the resulting blacked-out poem to standard output.`

const examples = `blackout --help
blackout 'lorem ipsum' --max-length 800`

var (
	Verbose       bool // Whether to print verbose results.
	MaxLength     int  // Maximum poem length to black out.
	PrintOriginal bool // Whether to print the original poem before blacking it out.
	Profanities   bool // Whether to filter out poems with offensive words while searching.
)

// rootCmd represents the base command when called without any sub-commands.
var rootCmd = &cobra.Command{
	Use:     "blackout <message>",
	Short:   "Make a blackout poem with the given hidden message",
	Long:    longDescription,
	Version: Version,
	Args:    cobra.ExactArgs(1),
	Run:     runApp,
	Example: examples,
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
	rootCmd.PersistentFlags().BoolVarP(&PrintOriginal, "print-original", "o", false, "print original poem before blacking out")
	rootCmd.PersistentFlags().BoolVarP(&Profanities, "allow-profanities", "p", false, "allow blacking out poems with profanities")
}

// runApp runs the CLI application.
func runApp(cmd *cobra.Command, args []string) {
	if !Verbose {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stdout)
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
	sp := SearchParams{runtime.NumCPU(), dataFolderPoems, blackoutRegex, MaxLength, Profanities}
	poemID, err := searchPoemsFolder(sp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found poem ID %d to black out\n", poemID)
	poem, err := json2poem(filepath.Join(dataFolderPoems, poemFilename(poemID)))
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
	printErr := PrintBlackoutPoem(poem, args[0])
	if printErr != nil {
		log.Fatal(err)
	}
}
