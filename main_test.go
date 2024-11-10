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
	"flag"
	"io"
	"os"
	"testing"
)

func TestMainFunction(t *testing.T) {
	// Save the original command-line arguments and output
	origArgs := os.Args
	origStdout := os.Stdout

	// Restore the original command-line arguments and output after the test
	defer func() {
		os.Args = origArgs
		os.Stdout = origStdout
	}()

	// Set the command-line arguments
	os.Args = []string{"cmd", "-dir=test_dir"}

	// Create a pipe to capture the output
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	// Redirect stdout to the pipe's writer
	os.Stdout = writer

	// Reset the flag package to allow re-parsing of flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Call the main function in a separate goroutine
	done := make(chan struct{})
	go func() {
		main()
		writer.Close()
		close(done)
	}()

	// Read the output from the pipe's reader
	var out bytes.Buffer
	if _, err := io.Copy(&out, reader); err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}

	// Wait for the main function to finish
	<-done

	// Verify the output
	expectedOutput := "Total lines of code: 126\n"
	if !bytes.Equal(out.Bytes(), []byte(expectedOutput)) {
		t.Errorf("expected output to contain %q, got %q", expectedOutput, out.String())
	}
}
