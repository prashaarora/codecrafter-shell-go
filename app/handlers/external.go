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
			stdoutFile, fileErr := OpenOutputFile(redirectInfo.StdoutFile, redirectInfo.StdOutAppend)
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
			stderrFile, fileErr := OpenOutputFile(redirectInfo.StderrFile, redirectInfo.StdErrAppend)
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
