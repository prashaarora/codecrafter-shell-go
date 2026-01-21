package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtins = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
	"pwd": true,
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		input := strings.Fields(command)
		if len(input) == 0 {
			continue
		}
		cmd := input[0]

		switch cmd {
		case "exit":
			handleExit(input)
		case "echo":
			handleEcho(input)
		case "type":
			handleType(input)
		case "pwd":
			handlePwd(input)
		default:
			handleExternal(input)

		}
	}
}

func handleExit(input []string) {
	os.Exit(0)
}

func handleEcho(input []string) {
	output := input[1:]
	fmt.Println(strings.Join(output, " "))
}

func handleType(input []string) {
	if len(input) < 2 {
		fmt.Fprintln(os.Stderr, "type: missing argument")
		return
	}
	cmdName := input[1]

	if builtins[cmdName] {
		fmt.Println(cmdName + " is a shell builtin")
	} else {
		path, err := exec.LookPath(cmdName)
		if err == nil {
			fmt.Println(cmdName + " is " + path)
		} else {
			fmt.Fprintln(os.Stderr, cmdName + ": command not found")
		}
	}

}

func handlePwd(input []string) {
	if len(input) > 1{
		fmt.Fprintln(os.Stderr, "pwd: too many arguments")
		return
	}
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(path)
}

func handleExternal(input []string) {
	cmdName := input[0]
	args := input[1:]
	cmdPath, err := exec.LookPath(cmdName)
	if err == nil {
		exeCommand := exec.Command(cmdPath, args...)
		exeCommand.Stdout = os.Stdout
		exeCommand.Stderr = os.Stderr
		execErr := exeCommand.Run()
		if execErr != nil {
			fmt.Fprintln(os.Stderr, execErr)
		}
	} else {
		fmt.Fprintln(os.Stderr, cmdName + ": command not found")
	}
}
