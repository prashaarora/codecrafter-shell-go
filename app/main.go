package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var _ = fmt.Fprint
var _ = os.Stdout
func checkFileExists(filepath string)bool{
	_, error := os.Stat(filepath)
	return !errors.Is(error, os.ErrNotExist)
}
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
				builtins := map[string]bool{
					"echo" : true,
					"exit" : true,
					"type" : true,
				}

				if builtins[typeOutput]{
					fmt.Println(typeOutput + " is a shell builtin")
				} else {
					path, err := exec.LookPath(typeOutput)
					if err == nil {
						fmt.Println(typeOutput + " is " + path)
					} else{
						fmt.Println(typeOutput + ": not found")
					}
				}
			default:
				fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
