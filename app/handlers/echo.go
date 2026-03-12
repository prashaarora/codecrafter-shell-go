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
		stdoutFile, fileErr := OpenOutputFile(redirectInfo.StdoutFile, redirectInfo.StdOutAppend)
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
		stderrFile, err := OpenOutputFile(redirectInfo.StderrFile, redirectInfo.StdErrAppend)
		if err != nil {
			fmt.Fprintln(os.Stderr, MsgErrorFileCreation+redirectInfo.StderrFile, err)
			return
		}
		defer stderrFile.Close()
	}
}
