package main

import (
	"fmt"
	"shell-in-go/shell"
)

func main() {
	fmt.Println("Welcome to Go Unix-Like Shell!")
	fmt.Println("Type 'exit' or ctrl+c to quit.")
	fmt.Println("------------------------------")

	// Start the shell loop
	shell.StartShell()
}
