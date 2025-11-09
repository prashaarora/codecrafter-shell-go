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
		//wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		switch {
			case strings.TrimRight(command, "\n") == "exit 0":
				os.Exit(0)
			case strings.HasPrefix(command, "echo"):
				output := strings.TrimPrefix(command, "echo")
				fmt.Println(strings.TrimSpace(output))
			case strings.HasPrefix(command, "type"):
				typeOutput := strings.TrimSpace(strings.TrimPrefix(command, "type"))
				if typeOutput == "echo" || typeOutput == "exit" || typeOutput == "type"{
					fmt.Println(typeOutput + " is a shell builtin")
				} else {
					fmt.Println(typeOutput + ": command not found")
				}
			default:
				fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}

	// for {
	// fmt.Fprint(os.Stdout, "$ ")
	// // Wait for user input
	// command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error reading input:", err)
	// 	os.Exit(1)
	// }
	// if strings.TrimRight(command, "\n") == "exit 0"{
	// 	os.Exit(0)
	// } else if strings.HasPrefix(command, "echo"){
    //    output := strings.TrimPrefix(command, "echo")
	//    fmt.Println(strings.TrimSpace(output))
	// } else {
	// 	fmt.Println(command[:len(command)-1] + ": command not found")
	// }		
	// }
