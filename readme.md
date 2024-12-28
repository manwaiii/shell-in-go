# Go Unix-Like Shell

This project is a simple Unix-like shell implemented in Go. It supports some common and basic shell commands such as `cd`, `ls`, `env`, `rm`, `cat`, `cp`, and `clear`.

## Built-in Commands

- `cd [path]`: Change the current directory.
- `ls [path]`: List the contents of a directory.
- `env`: Print all environment variables.
- `rm [file|directory]`: Remove files or directories.
- `cat [file1] [file2] ...`: Display the contents of files.
- `cp [source] [destination]`: Copy a file from source to destination.
- `clear`: Clear current terminal screen.

## Running the Shell

To run the shell, use the following command:

```sh
# run the project
go run main.go

# test the project
go test -v ./...

# format the project
go fmt ./...
```

## Usage

Once the shell is running, you can use the built-in commands as well as any other commands available on your system. To exit the shell, type `exit` or press Ctrl+C.

## Example
```she
ls
cd shell
```

## Improvement Inspirations

- Autocomplete: Implement command and file path autocompletion to enhance user experience.
- Configuration File: Allow users to configure the shell using a configuration file (e.g., `.shellrc`).
- Error Handling and Logging: Improve error handling and add logging to capture errors and debug information.
- Makefile: Create Makefile for more custom command