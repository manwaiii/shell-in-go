package builtins

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestChangeDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir := os.TempDir()
	defer os.Chdir(tempDir) // Change back to the original directory after the test

	// Test changing to the temporary directory
	err := ChangeDirectory([]string{tempDir})
	if err != nil {
		t.Errorf("ChangeDirectory failed: %v", err)
	}

	// Test changing to a non-existent directory
	err = ChangeDirectory([]string{"/nonexistent"})
	if err == nil {
		t.Error("ChangeDirectory should have failed for non-existent directory")
	}
}

func TestRemoveFileOrDir(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close()

	// Test removing the file
	err = RemoveFileOrDir([]string{tempFile.Name()})
	if err != nil {
		t.Errorf("RemoveFileOrDir failed: %v", err)
	}
}

func TestConcatenateFiles(t *testing.T) {
	// Create a temporary file with content
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	content := "Hello, World!"
	tempFile.WriteString(content)
	tempFile.Close()

	// Capture the output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = ConcatenateFiles([]string{tempFile.Name()})
	if err != nil {
		t.Errorf("ConcatenateFiles failed: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if output != content {
		t.Errorf("ConcatenateFiles output mismatch: got %v, want %v", output, content)
	}
}

func TestCopyFile(t *testing.T) {
	// Create a temporary source file with content
	srcFile, err := os.CreateTemp("", "srcfile")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	content := "Hello, World!"
	srcFile.WriteString(content)
	srcFile.Close()

	// Create a temporary destination file
	destFile, err := os.CreateTemp("", "destfile")
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	destFile.Close()

	// Test copying the file
	err = CopyFile([]string{srcFile.Name(), destFile.Name()})
	if err != nil {
		t.Errorf("CopyFile failed: %v", err)
	}

	// Verify the content of the destination file
	data, err := os.ReadFile(destFile.Name())
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(data) != content {
		t.Errorf("CopyFile content mismatch: got %v, want %v", string(data), content)
	}
}

// Helper function to capture and discard output
func runClearScreenWithOutputCapture() error {
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

	err := ClearScreen()
	w.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	<-done

	return err
}

func TestClearScreen(t *testing.T) {
	err := runClearScreenWithOutputCapture()
	if err != nil {
		t.Errorf("ClearScreen failed: %v", err)
	}
}
