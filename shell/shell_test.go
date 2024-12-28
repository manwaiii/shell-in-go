package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// Mock the ExecuteCommand function
func mockExecuteCommand(input string) error {
	if input == "mock_error" {
		return fmt.Errorf("mock error")
	}
	return nil
}

// Helper function to replace os.Stdin and os.Stdout
func withMockedIO(input string, testFunc func()) (string, error) {
	// Create a pipe to simulate user input
	inputReader, inputWriter, _ := os.Pipe()
	// Create a pipe to capture the output
	outputReader, outputWriter, _ := os.Pipe()

	// Replace os.Stdin and os.Stdout with the pipes
	originalStdin := os.Stdin
	originalStdout := os.Stdout
	os.Stdin = inputReader
	os.Stdout = outputWriter
	defer func() {
		os.Stdin = originalStdin
		os.Stdout = originalStdout
	}()

	// Write test commands to the input pipe
	go func() {
		defer inputWriter.Close()
		inputWriter.WriteString(input)
	}()

	// Capture the output
	var outputBuffer bytes.Buffer
	done := make(chan struct{})
	go func() {
		defer outputReader.Close()
		io.Copy(&outputBuffer, outputReader)
		close(done)
	}()

	// Run the test function
	testFunc()

	// Close the writer to finish capturing
	outputWriter.Close()
	<-done

	// Read the captured output
	return outputBuffer.String(), nil
}

func TestStartShell(t *testing.T) {
	// Mock the ExecuteCommandWrapper function by temporarily replacing it
	originalExecuteCommandWrapper := ExecuteCommandWrapper
	ExecuteCommandWrapper = mockExecuteCommand
	defer func() { ExecuteCommandWrapper = originalExecuteCommandWrapper }()

	// Define the input and expected output
	input := "exit\n"
	expectedOutput := "Exiting. Goodbye!"

	// Run the StartShell function with mocked IO
	output, err := withMockedIO(input, StartShell)
	if err != nil {
		t.Fatalf("Failed to capture output: %v", err)
	}

	// Check the output
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain %q, but got %q", expectedOutput, output)
	}
}
