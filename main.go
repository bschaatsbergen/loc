// loc - count lines of code in a directory
// BSD 3-Clause License
//
// Copyright (c) 2024, Alex Gaetano Padula
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
//  3. Neither the name of the copyright holder nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Loc is the main struct for the Loc program
type Loc struct {
	TotalLines int     // Total number of lines of code
	Config     *Config // The Loc configuration
	Directory  string  // The directory to scan
}

// LanguageConfig is the configuration for a language
type LanguageConfig struct {
	SkipPatterns []string `json:"skip_patterns"` // Patterns to skip; lines matching these patterns will not be counted
	Extensions   []string `json:"extensions"`    // File extensions to count
}

// Config is the configuration for Loc
type Config struct {
	Languages map[string]LanguageConfig `json:"languages"`
}

// CONFIG_FILE is the name of the configuration file
const CONFIG_FILE = "config.json"

// readConfig reads the Loc configuration file
func readConfig() (*Config, error) {
	// get working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// read the config files contents into memory
	configFile, err := os.ReadFile(wd + string(os.PathSeparator) + CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	// create config variable
	var config Config

	// unmarshal the config file into the config variable
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	// return the config variable
	return &config, nil
}

// scan scans the directory and counts the lines of code
func (loc *Loc) scan() error {
	// Walk the directory
	return filepath.Walk(loc.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil { // if there is an error, return the error
			return err
		}
		if !info.IsDir() { // if the file is not a directory we can count the lines of code

			for _, langConfig := range loc.Config.Languages { // iterate over configured languages
				for _, ext := range langConfig.Extensions { // iterate over the extensions for the language
					if strings.HasSuffix(path, ext) { // if the file has the correct extension
						lines, err := loc.countLines(path, langConfig.SkipPatterns) // count the lines of code
						if err != nil {
							return err
						}
						loc.TotalLines += lines // add the lines of code to the total
					}
				}
			}
		}
		return nil
	})
}

// countLines counts the lines of code in a file
func (loc *Loc) countLines(filePath string, skipPatterns []string) (int, error) {
	// we need to open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close() // defer the closure of the file

	scanner := bufio.NewScanner(file) // create a scanner for the file

	totalLines := 0 // total lines of code in the file

	// create a slice of regular expressions for the skip patterns
	skipRegexps := make([]*regexp.Regexp, len(skipPatterns))
	for i, pattern := range skipPatterns {
		skipRegexps[i] = regexp.MustCompile(pattern) // compile the regular expression
	}

	for scanner.Scan() { // iterate over the lines of the file
		line := scanner.Text() // get the line of the file
		skip := false          // set skip to false
		for _, re := range skipRegexps {
			if re.MatchString(line) { // if the line matches the regular expression we skip it
				skip = true
				break
			}
		}
		if !skip { // if we are not skipping the line we increment the total lines
			totalLines++
		}
	}

	// check for scanner errors
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return totalLines, nil
}

// cloneRepo clones a GitHub repository to a temporary directory
func cloneRepo(repoURL string) (string, error) {
	tempDir, err := os.MkdirTemp("", "loc-repo-") // create a temporary directory
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "clone", repoURL, tempDir)
	err = cmd.Run()
	if err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}

	return tempDir, nil
}

func main() {
	var err error // global error variable
	loc := Loc{}  // create a new Loc struct

	dir := flag.String("dir", ".", "directory to count lines of code")           // create a flag for the directory
	repo := flag.String("repo", ".", "github repository to count lines of code") // create a flag for a repository

	flag.Parse() // parse the flags

	loc.Directory = *dir // set the directory

	// directory supercedes repo
	if loc.Directory == "" { // if the directory is empty
		fmt.Println("Directory is empty") // print an error
		os.Exit(1)
	} else {
		if *repo != "" {
			loc.Directory, err = cloneRepo(*repo)
			if err != nil {
				fmt.Println("Error cloning repository:", err)
				return
			}
			defer os.RemoveAll(loc.Directory)
		} else {
			loc.Directory = *dir
		}
	}

	// Read the config
	loc.Config, err = readConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	// Scan the directory and count lines of code
	err = loc.scan()
	if err != nil {
		fmt.Println("Error scanning directory:", err)
		return
	}

	fmt.Printf("Total lines of code: %d\n", loc.TotalLines) // print the total lines of code
}
