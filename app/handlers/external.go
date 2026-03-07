package handlers

import (
	"fmt"
	"os"
	"os/exec"
)

func HandleExternal(input []string) {
	cmdName := input[0]
	args := input[1:]
	redirectInfo := HandleRedirection(args)
	_, err := exec.LookPath(cmdName)
	if err == nil {
		exeCommand := exec.Command(cmdName, redirectInfo.CleanArgs...)
		if redirectInfo.HasStdout {
			var stdoutFile *os.File
			var fileErr error
			if redirectInfo.StdOutAppend {
				stdoutFile, fileErr = OpenFileinAppendMode(redirectInfo.StdoutFile)
			} else {
				stdoutFile, fileErr = CreateFile(redirectInfo.StdoutFile)
			}
			if fileErr != nil {
				fmt.Fprintln(os.Stderr, MsgErrorFileCreation+redirectInfo.StdoutFile, fileErr)
				return
			}
			defer stdoutFile.Close()
			exeCommand.Stdout = stdoutFile
		} else {
			exeCommand.Stdout = os.Stdout
		}
		if redirectInfo.HasStderr {
			var stderrFile *os.File
			var fileErr error
			if redirectInfo.StdErrAppend {
				stderrFile, fileErr = OpenFileinAppendMode(redirectInfo.StderrFile)
			} else {
				stderrFile, fileErr = CreateFile(redirectInfo.StderrFile)
			}
			if fileErr != nil {
				fmt.Fprintln(os.Stderr, MsgErrorFileCreation+redirectInfo.StderrFile, fileErr)
				return
			}
			defer stderrFile.Close()
			exeCommand.Stderr = stderrFile
		} else {
			exeCommand.Stderr = os.Stderr
		}
		exeCommand.Run()
	} else {
		fmt.Fprintln(os.Stderr, cmdName+MsgNotFound)
	}
}
