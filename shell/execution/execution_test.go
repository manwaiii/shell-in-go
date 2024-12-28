package execution

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// Helper function to capture and discard output
func runWithOutputCapture(cmd string) error {
	r, w, _ := os.Pipe()
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	os.Stdout = w
	os.Stderr = w

	outputBuffer := &bytes.Buffer{}
	done := make(chan struct{})
	go func() {
		io.Copy(outputBuffer, r)
		close(done)
	}()

	err := ExecuteCommand(cmd)
	w.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	<-done

	return err
}

func TestExecuteCommand(t *testing.T) {
	// Test setting an environment variable
	err := runWithOutputCapture("TEST_VAR=test_value")
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}
	if os.Getenv("TEST_VAR") != "test_value" {
		t.Error("ExecuteCommand did not set the environment variable correctly")
	}

	// Test built-in command 'ls'
	err = runWithOutputCapture("ls")
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}

	// Test built-in command 'cd'
	err = runWithOutputCapture("cd /")
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}

	// Test built-in command 'env'
	err = runWithOutputCapture("env")
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}

	// Test built-in command 'rm'
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close()
	err = runWithOutputCapture("rm " + tempFile.Name())
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}

	// Test built-in command 'cat'
	tempFile, err = os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	content := "Hello, World!"
	tempFile.WriteString(content)
	tempFile.Close()
	err = runWithOutputCapture("cat " + tempFile.Name())
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}

	// Test built-in command 'cp'
	srcFile, err := os.CreateTemp("", "srcfile")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	srcFile.WriteString(content)
	srcFile.Close()
	destFile, err := os.CreateTemp("", "destfile")
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	destFile.Close()
	err = runWithOutputCapture("cp " + srcFile.Name() + " " + destFile.Name())
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}
}
