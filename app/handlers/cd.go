package handlers

import (
	"fmt"
	"os"
)

func HandleCd(input []string) {
	var target string

	switch len(input) {
	case 1:
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, MsgHomeNotSet)
			return
		}
		target = home
	case 2:
		target = input[1]
		if target == "~" {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, MsgHomeNotSet)
				return
			}
			target = home
		}
	default:
		fmt.Fprintln(os.Stderr, "cd"+MsgTooManyArgs)
		return
	}

	if err := os.Chdir(target); err != nil {
		fmt.Fprintln(os.Stderr, "cd: "+target+MsgNoSuchFile)
	}
}
