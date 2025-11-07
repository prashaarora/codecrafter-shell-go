package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout
func main() {
	for {
	fmt.Fprint(os.Stdout, "$ ")
	// Wait for user input
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}
	if strings.TrimRight(command, "\n") == "exit 0"{
		os.Exit(0)
	} else if strings.HasPrefix(command, "echo"){
       output := strings.TrimPrefix(command, "echo")
	   fmt.Println(strings.TrimSpace(output))
	} else {
		fmt.Println(command[:len(command)-1] + ": command not found")
	}		
	}
}
