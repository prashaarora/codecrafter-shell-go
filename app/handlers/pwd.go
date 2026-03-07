package handlers

import (
	"fmt"
	"os"
)

func HandlePwd(input []string) {
	if len(input) > 1 {
		fmt.Fprintln(os.Stderr, "pwd"+MsgTooManyArgs)
		return
	}
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(path)
}
