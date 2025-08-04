package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin": 
		cmd = exec.Command("open", url)
	default: 
		cmd = exec.Command("xdg-open", url) 
	}

	return cmd.Start() 
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter something (type 'exit' to quit):")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break 
		}
		input := strings.TrimSpace(scanner.Text())

		commands := strings.Split(input, " ")
		fmt.Printf("Commands entered: %v\n", commands)
		fmt.Printf("Type of var: %T\n", commands)

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		fmt.Println("You entered:", input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}
}
