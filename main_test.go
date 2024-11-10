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
