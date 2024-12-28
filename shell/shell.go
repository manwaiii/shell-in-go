package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"shell-in-go/shell/execution"
)

// ExecuteCommandWrapper is a wrapper function for execution.ExecuteCommand
var ExecuteCommandWrapper = execution.ExecuteCommand

// StartShell starts the main shell loop
func StartShell() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Display current working directory in prompt
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			continue
		}
		fmt.Printf("%s >> ", cwd)

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nExiting. Goodbye!")
				os.Exit(0)
			}
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Handle exit
		if input == "exit" {
			fmt.Println("Exiting. Goodbye!")
			break
		}

		// Parse and execute the command
		if err := ExecuteCommandWrapper(input); err != nil {
			fmt.Println("[Error in Shell] Error:", err)
		}
	}
}
