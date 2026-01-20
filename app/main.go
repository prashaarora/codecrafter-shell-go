package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
			case strings.HasPrefix(strings.TrimRight(command, "\n"), "exit"):
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
				fields := strings.Fields(command)
				if len(fields) == 0{
					continue
				}
				cmdName := fields[0]
				argName := fields[1:]
				_, err := exec.LookPath(cmdName)
				if err == nil {
					exeCommand := exec.Command(cmdName, argName...)
					exeCommand.Stdout = os.Stdout
					exeCommand.Stderr = os.Stderr
					err := exeCommand.Run()
					if err != nil{
						fmt.Println(err)
					}
				} else {
					fmt.Println(cmdName + ": command not found")
				}
		}
	}
}
