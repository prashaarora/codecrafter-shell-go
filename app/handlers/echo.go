package handlers

import (
	"fmt"
	"os"
	"strings"
)

func HandleEcho(input []string) {
	args := input[1:]
	redirectInfo := HandleRedirection(args)

	output := strings.Join(redirectInfo.CleanArgs, " ")

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

		fmt.Fprintln(stdoutFile, output)
	} else {
		fmt.Println(output)
	}

	if redirectInfo.HasStderr {
		var stderrFile *os.File
		var err error
		if redirectInfo.StdErrAppend {
			stderrFile, err = OpenFileinAppendMode(redirectInfo.StderrFile)
		} else {
			stderrFile, err = CreateFile(redirectInfo.StderrFile)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, MsgErrorFileCreation+redirectInfo.StderrFile, err)
			return
		}
		defer stderrFile.Close()
	}
}
