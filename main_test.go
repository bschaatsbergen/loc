// loc tests
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
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainFunction(t *testing.T) {
	origArgs := os.Args
	origStdout := os.Stdout

	// Restore the original command-line arguments and output after the test
	defer func() {
		os.Args = origArgs
		os.Stdout = origStdout
	}()

	os.Args = []string{"cmd", "-dir=test_dir"}

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = writer

	// Reset the flag package to allow re-parsing of flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Call the main function in a separate goroutine
	done := make(chan struct{})
	go func() {
		main()
		_ = writer.Close()
		close(done)
	}()

	// Read the output from the pipe's reader
	var out bytes.Buffer
	if _, err := io.Copy(&out, reader); err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}

	// Wait for the main function to finish
	<-done

	output := out.String()
	if !strings.Contains(output, "Total lines of code:") {
		t.Errorf("expected output to contain 'Total lines of code:', got %q", output)
	}
}

func setupTestDirectory(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "loc-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	dirs := []string{
		"src",
		"src/tests",
		"build",
		"node_modules/package",
		"gen",
		".git",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
	// Create some test files with various contents
	files := map[string]string{
		"main.go":                       "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
		"src/utils.go":                  "package src\n\nfunc Add(a, b int) int {\n\treturn a + b\n}\n",
		"src/utils_test.go":             "package src\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) {\n\t// test code\n}\n",
		"src/tests/integration.go":      "package tests\n\n// integration test\nfunc TestIntegration() {\n}\n",
		"app.spec.ts":                   "describe('app', () => {\n\tit('works', () => {\n\t\t// test\n\t});\n});\n",
		"component.test.js":             "test('component', () => {\n\t// test code\n});\n",
		"build/output.go":               "// generated file\npackage main\n\nvar Generated = true\n",
		"gen/models.go":                 "// Auto-generated file\npackage gen\n\ntype Model struct{}\n",
		"node_modules/package/index.js": "module.exports = {};\n",
		".git/config":                   "[core]\n\trepositoryformatversion = 0\n",
		"README.md":                     "# Test Project\n\nThis is a test.\n",
	}

	for file, content := range files {
		filePath := filepath.Join(tempDir, file)
		fileDir := filepath.Dir(filePath)
		err := os.MkdirAll(fileDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for file %s: %v", file, err)
		}

		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	config := Config{
		Languages: map[string]LanguageConfig{
			"go": {
				Extensions:   []string{".go"},
				SkipPatterns: []string{`^\s*//`, `^\s*$`}, // Skip comments and empty lines
			},
			"typescript": {
				Extensions:   []string{".ts"},
				SkipPatterns: []string{`^\s*//`, `^\s*$`},
			},
			"javascript": {
				Extensions:   []string{".js"},
				SkipPatterns: []string{`^\s*//`, `^\s*$`},
			},
			"markdown": {
				Extensions:   []string{".md"},
				SkipPatterns: []string{`^\s*$`},
			},
		},
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	configPath := filepath.Join(tempDir, "config.json")
	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	return tempDir
}

func TestExcludePatterns(t *testing.T) {
	testDir := setupTestDirectory(t)
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(testDir)

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func(dir string) {
		_ = os.Chdir(dir)
	}(originalWd)

	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	tests := []struct {
		name            string
		excludePatterns []string
		expectedFiles   []string
		excludedFiles   []string
	}{
		{
			name:            "No exclusions",
			excludePatterns: []string{},
			expectedFiles: []string{
				"main.go",
				"src/utils.go",
				"src/utils_test.go",
				"src/tests/integration.go",
				"app.spec.ts",
				"component.test.js",
				"build/output.go",
				"gen/models.go",
				"node_modules/package/index.js",
				"README.md",
			},
		},
		{
			name:            "Exclude test files",
			excludePatterns: []string{`.*_test\.go$`, `.*\.spec\.ts$`, `.*\.test\.js$`},
			expectedFiles: []string{
				"main.go",
				"src/utils.go",
				"src/tests/integration.go",
				"build/output.go",
				"gen/models.go",
				"node_modules/package/index.js",
				"README.md",
			},
			excludedFiles: []string{
				"src/utils_test.go",
				"app.spec.ts",
				"component.test.js",
			},
		},
		{
			name:            "Exclude directories",
			excludePatterns: []string{`node_modules/.*`, `\.git/.*`, `build/.*`},
			expectedFiles: []string{
				"main.go",
				"src/utils.go",
				"src/utils_test.go",
				"src/tests/integration.go",
				"app.spec.ts",
				"component.test.js",
				"gen/models.go",
				"README.md",
			},
			excludedFiles: []string{
				"build/output.go",
				"node_modules/package/index.js",
			},
		},
		{
			name:            "Exclude generated and test files",
			excludePatterns: []string{`gen/.*`, `.*_test\.go$`, `tests/.*`},
			expectedFiles: []string{
				"main.go",
				"src/utils.go",
				"app.spec.ts",
				"component.test.js",
				"build/output.go",
				"node_modules/package/index.js",
				"README.md",
			},
			excludedFiles: []string{
				"gen/models.go",
				"src/utils_test.go",
				"src/tests/integration.go",
			},
		},
		{
			name: "Complex exclusion patterns",
			excludePatterns: []string{
				`.*_test\.go$`,  // Go test files
				`.*\.spec\.ts$`, // TypeScript spec files
				`.*\.test\.js$`, // JavaScript test files
				`/tests/`,       // Test directories
				`node_modules/`, // Node modules
				`build/`,        // Build directory
			},
			expectedFiles: []string{
				"main.go",
				"src/utils.go",
				"gen/models.go",
				"README.md",
			},
			excludedFiles: []string{
				"src/utils_test.go",
				"src/tests/integration.go",
				"app.spec.ts",
				"component.test.js",
				"build/output.go",
				"node_modules/package/index.js",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := &Loc{
				Directory: testDir,
			}

			config, err := readConfig()
			if err != nil {
				t.Fatalf("Failed to read config: %v", err)
			}
			loc.Config = config

			if len(tt.excludePatterns) > 0 {
				compiledPatterns, err := compileExcludePatterns(tt.excludePatterns)
				if err != nil {
					t.Fatalf("Failed to compile exclude patterns: %v", err)
				}
				loc.ExcludePatterns = compiledPatterns
			}

			processedFiles := make(map[string]bool)

			err = filepath.Walk(loc.Directory, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip if excluded
				if loc.shouldExcludeFile(path) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				// Only track files, not directories
				if !info.IsDir() {
					relPath, _ := filepath.Rel(testDir, path)
					// Skip config.json from tracking
					if relPath != "config.json" {
						processedFiles[relPath] = true
					}
				}

				return nil
			})

			if err != nil {
				t.Fatalf("Failed to walk directory: %v", err)
			}

			// Check that expected files are processed
			for _, expectedFile := range tt.expectedFiles {
				if !processedFiles[expectedFile] {
					t.Errorf("Expected file %s was not processed", expectedFile)
				}
			}

			// Check that excluded files are not processed
			for _, excludedFile := range tt.excludedFiles {
				if processedFiles[excludedFile] {
					t.Errorf("Excluded file %s was processed", excludedFile)
				}
			}

			// Also test the actual line counting
			loc.TotalLines = 0
			err = loc.scan()
			if err != nil {
				t.Fatalf("Failed to scan directory: %v", err)
			}

			if loc.TotalLines == 0 && len(tt.expectedFiles) > 0 {
				t.Error("Expected some lines to be counted but got 0")
			}

			t.Logf("Test %s: Processed %d files, counted %d lines",
				tt.name, len(processedFiles), loc.TotalLines)
		})
	}
}

func TestShouldExcludeFile(t *testing.T) {
	loc := &Loc{
		Directory: "/project",
	}

	patterns := []string{
		`.*_test\.go$`,
		`node_modules/.*`,
		`\.git/.*`,
		`gen/.*`,
	}

	compiledPatterns, err := compileExcludePatterns(patterns)
	if err != nil {
		t.Fatalf("Failed to compile patterns: %v", err)
	}
	loc.ExcludePatterns = compiledPatterns

	tests := []struct {
		path     string
		expected bool
	}{
		{"/project/main.go", false},
		{"/project/src/utils.go", false},
		{"/project/src/utils_test.go", true},
		{"/project/node_modules/package/index.js", true},
		{"/project/.git/config", true},
		{"/project/gen/models.go", true},
		{"/project/build/output.go", false},
		{"/project/app.spec.ts", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := loc.shouldExcludeFile(tt.path)
			if result != tt.expected {
				t.Errorf("shouldExcludeFile(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCompileExcludePatterns(t *testing.T) {
	tests := []struct {
		name        string
		patterns    []string
		expectError bool
	}{
		{
			name:        "Valid patterns",
			patterns:    []string{`.*_test\.go$`, `node_modules/.*`},
			expectError: false,
		},
		{
			name:        "Invalid regex",
			patterns:    []string{`[`}, // Invalid regex
			expectError: true,
		},
		{
			name:        "Empty patterns",
			patterns:    []string{},
			expectError: false,
		},
		{
			name:        "Complex patterns",
			patterns:    []string{`(test|spec)`, `\.(git|build)/`, `^gen/.*\.go$`},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := compileExcludePatterns(tt.patterns)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(compiled) != len(tt.patterns) {
					t.Errorf("Expected %d compiled patterns, got %d", len(tt.patterns), len(compiled))
				}
			}
		})
	}
}

func TestExcludeFlagsType(t *testing.T) {
	var flags excludeFlags

	err := flags.Set("pattern1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = flags.Set("pattern2")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "pattern1, pattern2"
	if flags.String() != expected {
		t.Errorf("Expected %q, got %q", expected, flags.String())
	}

	if len(flags) != 2 {
		t.Errorf("Expected 2 flags, got %d", len(flags))
	}
}

func TestMainWithExcludeFlags(t *testing.T) {
	testDir := setupTestDirectory(t)
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(testDir)

	origArgs := os.Args
	origWd, _ := os.Getwd()
	defer func() {
		os.Args = origArgs
		_ = os.Chdir(origWd)
	}()

	err := os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
		shouldCount bool
	}{
		{
			name:        "Basic directory scan",
			args:        []string{"cmd", "-dir", "."},
			expectError: false,
			shouldCount: true,
		},
		{
			name:        "Exclude test files",
			args:        []string{"cmd", "-dir", ".", "--exclude", `.*_test\.go$`, "--exclude", `.*\.spec\.ts$`},
			expectError: false,
			shouldCount: true,
		},
		{
			name:        "Exclude directories",
			args:        []string{"cmd", "-dir", ".", "--exclude", `node_modules/.*`, "--exclude", `\.git/.*`},
			expectError: false,
			shouldCount: true,
		},
		{
			name:        "Invalid regex pattern",
			args:        []string{"cmd", "-dir", ".", "--exclude", `[`},
			expectError: true,
			shouldCount: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ContinueOnError)

			var buf bytes.Buffer
			flag.CommandLine.SetOutput(&buf)

			os.Args = tt.args

			var loc Loc
			var excludePatterns excludeFlags

			dir := flag.String("dir", ".", "directory to count lines of code")
			flag.Var(&excludePatterns, "exclude", "regex pattern to exclude files/directories")

			err := flag.CommandLine.Parse(tt.args[1:])
			if err != nil && !tt.expectError {
				t.Fatalf("Unexpected flag parsing error: %v", err)
			}
			if err != nil && tt.expectError {
				return // Expected error, test passed
			}

			loc.Directory = *dir

			if len(excludePatterns) > 0 {
				compiledPatterns, err := compileExcludePatterns(excludePatterns)
				if err != nil {
					if tt.expectError {
						return // Expected error
					}
					t.Fatalf("Unexpected error compiling patterns: %v", err)
				}
				loc.ExcludePatterns = compiledPatterns
			}

			config, err := readConfig()
			if err != nil {
				t.Fatalf("Failed to read config: %v", err)
			}
			loc.Config = config

			err = loc.scan()
			if err != nil {
				t.Fatalf("Failed to scan: %v", err)
			}

			if tt.shouldCount && loc.TotalLines == 0 {
				t.Error("Expected some lines to be counted")
			}

			t.Logf("Test %s: Counted %d lines", tt.name, loc.TotalLines)
		})
	}
}
