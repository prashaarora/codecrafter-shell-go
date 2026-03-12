package handlers

import (
	"os"
)

type RedirectInfo struct {
	CleanArgs    []string
	StdoutFile   string
	StderrFile   string
	HasStdout    bool
	HasStderr    bool
	StdOutAppend bool
	StdErrAppend bool
}

func HandleRedirection(args []string) RedirectInfo {
	stdoutRedirectIndex := -1
	stderrRedirectIndex := -1
	stdOutAppend := false
	stdErrAppend := false

	for i, arg := range args {
		switch arg {
		case ">>", "1>>":
			stdoutRedirectIndex = i
			stdOutAppend = true
		case ">", "1>":
			stdoutRedirectIndex = i
			stdOutAppend = false
		case "2>>":
			stderrRedirectIndex = i
			stdErrAppend = true
		case "2>":
			stderrRedirectIndex = i
			stdErrAppend = false
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
		CleanArgs:    cleanArgs,
		StdoutFile:   stdoutFile,
		StderrFile:   stderrFile,
		HasStdout:    hasStdout,
		HasStderr:    hasStderr,
		StdOutAppend: stdOutAppend,
		StdErrAppend: stdErrAppend,
	}
}

func OpenFileInAppendMode(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func CreateFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

func OpenOutputFile(filename string, appendMode bool) (*os.File, error) {
	if appendMode {
		return OpenFileInAppendMode(filename)
	}
	return CreateFile(filename)
}
