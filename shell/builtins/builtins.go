package builtins

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

// ChangeDirectory changes the current directory
func ChangeDirectory(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: cd [path]")
	}
	return os.Chdir(args[0])
}

// PrintEnvironment prints all environment variables
func PrintEnvironment() error {
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	return nil
}

// ListDirectory mimics the Unix 'ls' command
func ListDirectory(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot access '%s': %w", dir, err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil
}

// RemoveFileOrDir removes files or directories
func RemoveFileOrDir(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: rm [file|directory]")
	}

	for _, target := range args {
		err := os.RemoveAll(target) // Remove file or directory recursively
		if err != nil {
			return fmt.Errorf("cannot remove '%s': %w", target, err)
		}
	}
	return nil
}

// ConcatenateFiles displays the contents of files
func ConcatenateFiles(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: cat [file1] [file2] ...")
	}

	for _, file := range args {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("cannot read file '%s': %w", file, err)
		}
		fmt.Print(string(data))
	}
	return nil
}

// CopyFile copies the source file to the destination
func CopyFile(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: cp [source] [destination]")
	}

	source := args[0]
	destination := args[1]

	// Open the source file
	srcFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("cannot open source file '%s': %w", source, err)
	}
	defer srcFile.Close()

	// Create the destination file
	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("cannot create destination file '%s': %w", destination, err)
	}
	defer destFile.Close()

	// Copy data from source to destination
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying from '%s' to '%s': %w", source, destination, err)
	}

	return nil
}

// ClearScreen clears the terminal screen
func ClearScreen() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
