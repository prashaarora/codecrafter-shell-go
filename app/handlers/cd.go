package handlers

import (
	"fmt"
	"os"
)

func HandleCd(input []string) {
	if len(input) != 2 {
		fmt.Fprintln(os.Stderr, "cd"+MsgExpect1Arg)
		return
	}
	args := input[1]
	if args == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cd"+MsgHomeNotSet)
			return
		}
		args = home
	}
	err := os.Chdir(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cd: "+args+MsgNoSuchFile)
	}
}
