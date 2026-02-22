package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

var builtins = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

const (
	msgNotFound          = ": not found"
	msgNoSuchFile        = ": No such file or directory"
	msgMissingArg        = ": missing argument"
	msgTooManyArgs       = ": too many arguments"
	msgIsBuiltin         = " is a shell builtin"
	msgIs                = " is "
	msgHomeNotSet        = "cd: HOME not set"
	msgErrorReading      = "Error reading input:"
	msgExpect1Arg        = ": expect 1 argument atleast"
	msgErrorFileCreation = "error creating file: "
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := reader.ReadString('\n') 
		if err != nil {
			fmt.Fprintln(os.Stderr, msgErrorReading, err)
			os.Exit(1)
		}
		input := parseCommand(command)
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
		case "cd":
			handleCd(input)
		default:
			handleExternal(input)

		}
	}
}

func parseCommand(cmd string) []string {
	p := parser.NewCommandParser()
	return p.Parse(cmd)
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
		fmt.Fprintln(os.Stderr, "type"+msgMissingArg)
		return
	}
	cmdName := input[1]

	if builtins[cmdName] {
		fmt.Println(cmdName + msgIsBuiltin)
	} else {
		path, err := exec.LookPath(cmdName)
		if err == nil {
			fmt.Println(cmdName + msgIs + path)
		} else {
			fmt.Fprintln(os.Stderr, cmdName+msgNotFound)
		}
	}

}

func handlePwd(input []string) {
	if len(input) > 1 {
		fmt.Fprintln(os.Stderr, "pwd"+msgTooManyArgs)
		return
	}
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(path)
}

func handleCd(input []string) {
	if len(input) != 2 {
		fmt.Fprintln(os.Stderr,"cd"+msgExpect1Arg)
		return
	}
	args := input[1]
	if args == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cd"+msgHomeNotSet)
			return
		}
		args = home
	}
	err := os.Chdir(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cd: "+args+msgNoSuchFile)
	}
}

func handleExternal(input []string) {
	cmdName := input[0]
	args := input[1:]
	cleanArgs, filename, hasRedirect := handleRedirection(args)
	_, err := exec.LookPath(cmdName)
	if err == nil {
		exeCommand := exec.Command(cmdName, cleanArgs...)
		if hasRedirect {
			outputFile, fileErr := os.Create(filename)
			if fileErr != nil {
				fmt.Fprintln(os.Stderr, msgErrorFileCreation+filename, fileErr)
				return
			}
			defer outputFile.Close()
			exeCommand.Stdout = outputFile
		} else {
			exeCommand.Stdout = os.Stdout
		}
		exeCommand.Stderr = os.Stderr
		execErr := exeCommand.Run()
		if execErr != nil {
			fmt.Fprintln(os.Stderr, execErr)
		}
	} else {
		fmt.Fprintln(os.Stderr, cmdName+msgNotFound)
	}
}

func handleRedirection(args []string) ([]string, string, bool) { 
	redirectIndex := -1
	for i, arg := range args{
		if arg == ">" || arg == "1>" {
			redirectIndex = i
			break
		}
	}
	if redirectIndex == -1 {
		return args, "", false
	}

	if redirectIndex+1 >= len(args) {
		return args, "", false
	}

	cmdArgs := args[:redirectIndex]
	filename := args[redirectIndex+1]
	return cmdArgs, filename, true
}
