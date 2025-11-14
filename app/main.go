package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
					pathVar := os.Getenv("PATH")
					pathDirs := strings.Split(pathVar, string(os.PathListSeparator))
					found := false
					for _, i := range pathDirs {
						filename := filepath.Join(i, typeOutput)
						isFileexist := checkFileExists(filename)
						if isFileexist{
							fileInfo, err := os.Stat(filename)
							if err == nil {
								if fileInfo.Mode().Perm() & 0111 != 0{
									fmt.Println(typeOutput + " is " + filename)
									found = true
									break
								}
							}
						} 
					}
					if !found{
						fmt.Println(typeOutput + ": not found")
					}
				}
			default:
				fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
