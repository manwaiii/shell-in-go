package execution

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"shell-in-go/shell/builtins"
	"strings"
)

// ExecuteCommand parses and executes the given command
func ExecuteCommand(input string) error {
	// Handle environment variable setting (e.g., VAR=value)
	if strings.Contains(input, "=") && !strings.Contains(input, " ") {
		parts := strings.SplitN(input, "=", 2)
		os.Setenv(parts[0], parts[1])
		return nil
	}

	// Split input for piping
	commands := strings.Split(input, "|")

	// Set up piping between commands
	var prevPipe io.ReadCloser
	for i, cmdLine := range commands {
		// Parse command and arguments
		cmdParts := strings.Fields(strings.TrimSpace(cmdLine))
		if len(cmdParts) == 0 {
			continue
		}
		cmdName := cmdParts[0]
		cmdArgs := expandEnvironmentVariables(cmdParts[1:])

		// Handle built-in commands
		if i == len(commands)-1 { // Execute last command normally
			switch cmdName {
			case "ls":
				return builtins.ListDirectory(cmdArgs)
			case "cd":
				return builtins.ChangeDirectory(cmdArgs)
			case "env":
				return builtins.PrintEnvironment()
			case "rm":
				return builtins.RemoveFileOrDir(cmdArgs)
			case "cat":
				return builtins.ConcatenateFiles(cmdArgs)
			case "cp":
				return builtins.CopyFile(cmdArgs)
			case "clear":
				return builtins.ClearScreen()
			}
		}

		// Prepare the command
		cmd := exec.Command(cmdName, cmdArgs...)

		// Set up input for the current command
		if prevPipe != nil {
			cmd.Stdin = prevPipe
		}

		// Set up output for the current command
		if i < len(commands)-1 {
			// Pipe the output to the next command
			outPipe, err := cmd.StdoutPipe()
			if err != nil {
				return fmt.Errorf("error creating pipe for command '%s': %w", cmdName, err)
			}
			prevPipe = outPipe
		} else {
			// Final command outputs to Stdout
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		// Run the command
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting command '%s': %w", cmdName, err)
		}

		// Wait for the command to complete
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("error executing command '%s': %w", cmdName, err)
		}
	}

	return nil
}

// ExpandEnvironmentVariables expands environment variables in arguments
func expandEnvironmentVariables(args []string) []string {
	for i, arg := range args {
		if strings.HasPrefix(arg, "$") {
			args[i] = os.Getenv(strings.TrimPrefix(arg, "$"))
		}
	}
	return args
}
