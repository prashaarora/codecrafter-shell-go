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

type RedirectInfo struct {
	CleanArgs  []string
	StdoutFile string
	StderrFile string
	HasStdout  bool
	HasStderr  bool
	StdOutAppend bool
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
	args := input[1:]
	redirectInfo := handleRedirection(args)

	output := strings.Join(redirectInfo.CleanArgs, " ")

	if redirectInfo.HasStdout {
		var stdoutFile *os.File
		var fileErr error
		if redirectInfo.StdOutAppend {
			stdoutFile, fileErr = openFileinAppendMode(redirectInfo.StdoutFile)
		} else {
			stdoutFile, fileErr = createFile(redirectInfo.StdoutFile)
		}
		if fileErr != nil {
			fmt.Fprintln(os.Stderr, msgErrorFileCreation+redirectInfo.StdoutFile, fileErr)
			return
		}
		defer stdoutFile.Close()

		fmt.Fprintln(stdoutFile, output)
	} else {
		fmt.Println(output)
	}
	if redirectInfo.HasStderr {
		stderrFile, err := os.Create(redirectInfo.StderrFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, msgErrorFileCreation, err)
			return
		}
		defer stderrFile.Close()
	}
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
		fmt.Fprintln(os.Stderr, "cd"+msgExpect1Arg)
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

func openFileinAppendMode(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func createFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

func handleExternal(input []string) {
	cmdName := input[0]
	args := input[1:]
	redirectInfo := handleRedirection(args)
	_, err := exec.LookPath(cmdName)
	if err == nil {
		exeCommand := exec.Command(cmdName, redirectInfo.CleanArgs...)
		if redirectInfo.HasStdout {
			var stdoutFile *os.File
			var fileErr error
			if redirectInfo.StdOutAppend {
				stdoutFile, fileErr = openFileinAppendMode(redirectInfo.StdoutFile)
			} else {
				stdoutFile, fileErr = createFile(redirectInfo.StdoutFile)
			}
			if fileErr != nil {
				fmt.Fprintln(os.Stderr, msgErrorFileCreation+redirectInfo.StdoutFile, fileErr)
				return
			}
			defer stdoutFile.Close()
			exeCommand.Stdout = stdoutFile
		} else {
			exeCommand.Stdout = os.Stdout
		}
		if redirectInfo.HasStderr {
			stderrFile, fileErr := createFile(redirectInfo.StderrFile)
			if fileErr != nil {
				fmt.Fprintln(os.Stderr, msgErrorFileCreation+redirectInfo.StderrFile, fileErr)
				return
			}
			defer stderrFile.Close()
			exeCommand.Stderr = stderrFile
		} else {
			exeCommand.Stderr = os.Stderr
		}
		exeCommand.Run()
	} else {
		fmt.Fprintln(os.Stderr, cmdName+msgNotFound)
	}
}

func handleRedirection(args []string) RedirectInfo {
	stdoutRedirectIndex := -1
	stderrRedirectIndex := -1
	stdOutAppend := false
	
	for i, arg := range args {
		switch arg {
		case ">>", "1>>":
			stdoutRedirectIndex = i
			stdOutAppend = true
		case ">", "1>":
			stdoutRedirectIndex = i
			stdOutAppend = false
		case "2>":
			stderrRedirectIndex = i
		}
	}
	
	stdoutFile := ""
	stderrFile := ""
	hasStdout := false
	hasStderr := false

	if stdoutRedirectIndex != -1 && stdoutRedirectIndex+1 < len(args) {
		stdoutFile = args[stdoutRedirectIndex+1]
		hasStdout = true
	}
	
	if stderrRedirectIndex != -1 && stderrRedirectIndex+1 < len(args) {
		stderrFile = args[stderrRedirectIndex+1]
		hasStderr = true
	}

	cleanArgs := make([]string, 0)
	for i, arg := range args {
	   if stdoutRedirectIndex != -1 && (i == stdoutRedirectIndex || i == stdoutRedirectIndex+1) {
	       continue
	   }
	   
	   if stderrRedirectIndex != -1 && (i == stderrRedirectIndex || i == stderrRedirectIndex+1) {
	       continue
	   }
	   
	   cleanArgs = append(cleanArgs, arg)
	}
	
	return RedirectInfo{
		CleanArgs:  cleanArgs,
		StdoutFile: stdoutFile,
		StderrFile: stderrFile,
		HasStdout:  hasStdout,
		HasStderr:  hasStderr,
		StdOutAppend: stdOutAppend,
	}
}
